package quotes

import (
	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log/level"

	"github.com/donmikel/wow/pkg/protocol"
)

// GetQuote returns random quote from the World of Wisdom book.
func (t *TCPServer) GetQuote(c easytcp.Context) {
	quote := t.service.GetQuote()

	err := c.SetResponse(protocol.GetQuoteResponse, &protocol.GetQuoteResponseMessage{
		Quote: quote,
	})
	if err != nil {
		level.Error(t.logger).Log("msg", "GetQuote failed to set response",
			"err", err,
		)
	}
}
