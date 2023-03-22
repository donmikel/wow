package quotes

import (
	"time"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/kit/metrics"
)

var _ ServerInterface = &instrumentingMiddleware{}

// NewInstrumentingMiddleware creates metrics middleware.
// It monitors all requests to TCP service.
func NewInstrumentingMiddleware(histogram metrics.Histogram, si ServerInterface) ServerInterface {
	return &instrumentingMiddleware{
		requestHistogram: histogram,
		next:             si,
	}
}

type instrumentingMiddleware struct {
	requestHistogram metrics.Histogram
	next             ServerInterface
}

// GetQuote returns random quote from the World of Wisdom book.
func (mw *instrumentingMiddleware) GetQuote(c easytcp.Context) {
	begin := time.Now()
	mw.next.GetQuote(c)
	labels := []string{
		"method", "GetQuote",
	}
	mw.requestHistogram.With(labels...).Observe(time.Since(begin).Seconds())
}
