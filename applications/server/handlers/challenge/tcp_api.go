package challenge

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log/level"

	"github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/pkg/guidetour"
	"github.com/donmikel/wow/pkg/protocol"
	srv "github.com/donmikel/wow/pkg/server"
)

// MakeChallenge initialize puzzle that informs the client to complete a guided tour.
func (t *TCPServer) MakeChallenge(c easytcp.Context) {
	sessionID := c.Session().ID().(string)
	clientAddress := srv.GetClientAddress(c.Session().Conn().RemoteAddr())
	if len(t.guideKeys) == 0 {
		err := c.SetResponse(protocol.Error, &protocol.ErrorResponse{
			Code:    server.ErrBadRequest,
			Message: "guide keys list is empty",
		})
		level.Error(t.logger).Log("msg", "guide keys list is empty")
		if err != nil {
			level.Error(t.logger).Log("msg", "failed to send error response",
				"err", err)
		}

		return
	}

	g := guidetour.NewGuideTour(clientAddress, t.GetTourLength, t.GetTimeStamp)

	t.makePuzzle(c, g, sessionID)
}

func (t *TCPServer) makePuzzle(c easytcp.Context, g *guidetour.GuideTour, sessionID string) {
	h0 := g.MakePuzzle(sessionID)
	err := c.SetResponse(protocol.ChallengeRequiredResponse, &protocol.POWChallengeRequiredResponse{
		Hash:       h0,
		TourLength: g.GetTourLength(),
		Guides:     t.guides,
	})
	if err != nil {
		level.Error(t.logger).Log("msg", "failed to send POWChallengeRequiredResponse", "err", err)
	}
}
