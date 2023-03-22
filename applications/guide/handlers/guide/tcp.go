package guide

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"

	"github.com/donmikel/wow/pkg/protocol"
)

// ServerInterface is a server interface.
//
//go:generate moqc -stub -out mock/server_interface.go -pkg mock . ServerInterface
type ServerInterface interface {
	// GetTourHash solves puzzle of the tour guide at the l-th stop.
	GetTourHash(c easytcp.Context)
}

// TCPServer implements tcp handlers for Server interface.
type TCPServer struct {
	// array of guide secret keys
	guideKey string
	// GetTourLength returns length of the guide tour.
	GetTourLength func() int
	GetTimeStamp  func() int64
	logger        log.Logger
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(
	guideKey string,
	getTourLength func() int,
	getTimeStamp func() int64,
	logger log.Logger,
) *TCPServer {
	return &TCPServer{
		guideKey:      guideKey,
		GetTimeStamp:  getTimeStamp,
		GetTourLength: getTourLength,
		logger:        logger,
	}
}

// RegisterHandlers registers tcp handlers.
func RegisterHandlers(si ServerInterface, tcpServer *easytcp.Server) {
	tcpServer.AddRoute(protocol.GuideRequest, si.GetTourHash)
}
