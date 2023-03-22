package quotes

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"

	app "github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/pkg/protocol"
)

// ServerInterface is a server interface.
//
//go:generate moqc -stub -out mock/server_interface.go -pkg mock . ServerInterface
type ServerInterface interface {
	// GetQuote returns random quote from the World of Wisdom book.
	GetQuote(c easytcp.Context)
}

// TCPServer implements tcp handlers for Server interface.
type TCPServer struct {
	service app.Service
	logger  log.Logger
}

// NewTCPServer creates a new instance of TCPServer.
func NewTCPServer(service app.Service, logger log.Logger) *TCPServer {
	return &TCPServer{
		service: service,
		logger:  logger,
	}
}

// RegisterHandlers registers tcp handlers.
func RegisterHandlers(t ServerInterface, tcpServer *easytcp.Server, middlewares ...easytcp.MiddlewareFunc) {
	tcpServer.AddRoute(protocol.GetQuoteRequest, t.GetQuote, middlewares...)
}
