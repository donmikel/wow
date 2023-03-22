package challenge

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/pkg/guidetour"
	"github.com/donmikel/wow/pkg/protocol"
	srv "github.com/donmikel/wow/pkg/server"
)

// POW implements a middleware which verifies the guide tour puzzle.
type POW struct {
	guideKeys []string
	// GetTourLength returns length of the guide tour.
	// We can implement more complex logic here to adjust the length of the tour depending on
	// the load or suspiciousness of the client
	GetTourLength func() int
	GetTimeStamp  func() int64
	logger        log.Logger
}

// NewPOWMiddleware creates an instance of POW.
func NewPOWMiddleware(guideKeys []string, getTourLength func() int, getTimeStamp func() int64, logger log.Logger) *POW {
	return &POW{
		logger:        logger,
		guideKeys:     guideKeys,
		GetTourLength: getTourLength,
		GetTimeStamp:  getTimeStamp,
	}
}

// GuideTour is a middleware for checking the guide tour puzzle.
func (p *POW) GuideTour(next easytcp.HandlerFunc) easytcp.HandlerFunc {
	return func(ctx easytcp.Context) {
		clientAddress := srv.GetClientAddress(ctx.Session().Conn().RemoteAddr())
		var puzzleHeader protocol.GuideTourPuzzleHeader
		// Each TCP request message must contain guide tour puzzle header.
		err := ctx.Bind(&puzzleHeader)
		if err != nil {
			err = ctx.SetResponse(protocol.Error, &protocol.ErrorResponse{
				Code:    server.ErrInternal,
				Message: "header parsing error",
				Inner:   err,
			})
			level.Error(p.logger).Log("msg", "header parsing error")
			if err != nil {
				level.Error(p.logger).Log("msg", "failed to send error response",
					"err", err)
			}

			return
		}

		if puzzleHeader.FinishHash == "" || puzzleHeader.StartHash == "" {
			err = ctx.SetResponse(protocol.Error, &protocol.ErrorResponse{
				Code:    server.ErrBadRequest,
				Message: "both hashes are required",
			})
			level.Error(p.logger).Log("msg", "both hashes are required")
			if err != nil {
				level.Error(p.logger).Log("msg", "failed to send error response",
					"err", err)
			}

			return
		}

		g := guidetour.NewGuideTour(clientAddress, p.GetTourLength, p.GetTimeStamp)

		isValid, err := g.VerifyPuzzle(puzzleHeader.StartHash, puzzleHeader.FinishHash, p.guideKeys)
		if err != nil {
			level.Error(p.logger).Log("msg", "verifying puzzle error", "err", err)
		}
		if !isValid {
			err = ctx.SetResponse(protocol.Error, &protocol.ErrorResponse{
				Code:    server.ErrBadRequest,
				Message: "puzzle is invalid",
			})
			level.Error(p.logger).Log("msg", "puzzle is invalid")
			if err != nil {
				level.Error(p.logger).Log("msg", "failed to send error response",
					"err", err)
			}

			return
		}

		next(ctx)
	}
}
