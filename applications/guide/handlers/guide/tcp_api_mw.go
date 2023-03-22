package guide

import (
	"time"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/kit/metrics"
)

var _ ServerInterface = &instrumentingMiddleware{}

// NewInstrumentingMiddleware creates metrics middleware.
// It monitors all requests to TCP service.
func NewInstrumentingMiddleware(histogram metrics.Histogram, t ServerInterface) ServerInterface {
	return &instrumentingMiddleware{
		requestHistogram: histogram,
		next:             t,
	}
}

type instrumentingMiddleware struct {
	requestHistogram metrics.Histogram
	next             ServerInterface
}

// GetTourHash solves puzzle of the tour guide at the l-th stop.
func (mw *instrumentingMiddleware) GetTourHash(c easytcp.Context) {
	begin := time.Now()
	mw.next.GetTourHash(c)
	labels := []string{
		"method", "GetTourHash",
	}
	mw.requestHistogram.With(labels...).Observe(time.Since(begin).Seconds())
}
