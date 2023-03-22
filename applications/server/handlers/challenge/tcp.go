package challenge

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"

	"github.com/donmikel/wow/pkg/protocol"
)

// ServerInterface is a server interface.
//
//go:generate moqc -stub -out mock/server_interface.go -pkg mock . ServerInterface
type ServerInterface interface {
	// MakeChallenge initialize puzzle that informs the client to complete a guided tour.
	MakeChallenge(c easytcp.Context)
}

// TCPServer implements tcp handlers for Server interface.
type TCPServer struct {
	// array of guide addresses
	guides []string
	// array of guide secret keys
	guideKeys []string
	// GetTourLength returns length of the guide tour.
	// We can implement more complex logic here to adjust the length of the tour depending on
	// the load or suspiciousness of the client
	GetTourLength func() int
	GetTimeStamp  func() int64
	logger        log.Logger
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(
	guides []string,
	guideKeys []string,
	getTourLength func() int,
	getTimeStamp func() int64,
	logger log.Logger,
) *TCPServer {
	return &TCPServer{
		guides:        guides,
		guideKeys:     guideKeys,
		GetTourLength: getTourLength,
		GetTimeStamp:  getTimeStamp,
		logger:        logger,
	}
}

// RegisterHandlers registers tcp handlers.
func RegisterHandlers(t ServerInterface, tcpServer *easytcp.Server, middlewares ...easytcp.MiddlewareFunc) {
	tcpServer.AddRoute(protocol.ChallengeRequest, t.MakeChallenge)
}
