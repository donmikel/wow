package guide

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/donmikel/wow/applications/client/interfaces"
)

var _ interfaces.Guide = &instrumentingMiddleware{}

// NewInstrumentingMiddleware creates metrics middleware.
// It monitors all requests to Guide service.
func NewInstrumentingMiddleware(histogram metrics.Histogram, p interfaces.Guide) interfaces.Guide {
	return &instrumentingMiddleware{
		requestHistogram: histogram,
		next:             p,
	}
}

type instrumentingMiddleware struct {
	requestHistogram metrics.Histogram
	next             interfaces.Guide
}

// GetTourHash solves puzzle of the tour guide at the l-th stop.
func (mw *instrumentingMiddleware) GetTourHash(hash string, stepNumber int, tourLength int) (string, error) {
	begin := time.Now()
	hi, err := mw.next.GetTourHash(hash, stepNumber, tourLength)
	labels := []string{
		"method", "GetTourHash",
		"error", fmt.Sprint(err != nil),
	}
	mw.requestHistogram.With(labels...).Observe(time.Since(begin).Seconds())

	return hi, err
}
