package guide

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log/level"

	"github.com/donmikel/wow/applications/guide"
	"github.com/donmikel/wow/pkg/guidetour"
	"github.com/donmikel/wow/pkg/protocol"
	"github.com/donmikel/wow/pkg/server"
)

// GetTourHash solves puzzle of the tour guide at the l-th stop.
func (t *TCPServer) GetTourHash(c easytcp.Context) {
	var reqData protocol.GuideRequestMessage

	if err := c.Bind(&reqData); err != nil {
		level.Error(t.logger).Log("msg", "error parsing GuideRequestMessage",
			"err", err,
		)

		return
	}

	if reqData.Hash == "" {
		err := c.SetResponse(protocol.Error, &protocol.ErrorResponse{
			Code:    guide.ErrBadRequest,
			Message: "hash is required",
		})
		if err != nil {
			level.Error(t.logger).Log("msg", "Error failed to set response",
				"err", err,
			)
		}

		return
	}
	if reqData.StepNumber > reqData.TourLength {
		err := c.SetResponse(protocol.Error, &protocol.ErrorResponse{
			Code:    guide.ErrBadRequest,
			Message: "step number must be less or equal to the length of the tour",
		})
		if err != nil {
			level.Error(t.logger).Log("msg", "Error failed to set response",
				"err", err,
			)
		}

		return
	}

	clientAddress := server.GetClientAddress(c.Session().Conn().RemoteAddr())
	guideTour := guidetour.NewGuideTour(clientAddress, t.GetTourLength, t.GetTimeStamp)

	hi := guideTour.TourGuideHash(reqData.Hash, reqData.StepNumber, reqData.TourLength, t.guideKey)

	err := c.SetResponse(protocol.GuideResponse, &protocol.GuideResponseMessage{
		Hash: hi,
	})
	if err != nil {
		level.Error(t.logger).Log("msg", "GuideResponse failed to set response",
			"err", err,
		)
	}
}
