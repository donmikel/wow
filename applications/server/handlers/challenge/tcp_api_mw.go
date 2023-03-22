package challenge

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

func (mw *instrumentingMiddleware) MakeChallenge(c easytcp.Context) {
	begin := time.Now()
	mw.next.MakeChallenge(c)
	labels := []string{
		"method", "MakeChallenge",
	}
	mw.requestHistogram.With(labels...).Observe(time.Since(begin).Seconds())
}
