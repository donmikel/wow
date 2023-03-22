package server

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/donmikel/wow/applications/client"
	"github.com/donmikel/wow/applications/client/interfaces"
)

var _ interfaces.Server = &instrumentingMiddleware{}

// NewInstrumentingMiddleware creates metrics middleware.
// It monitors all requests to Server service.
func NewInstrumentingMiddleware(histogram metrics.Histogram, p interfaces.Server) interfaces.Server {
	return &instrumentingMiddleware{
		requestHistogram: histogram,
		next:             p,
	}
}

type instrumentingMiddleware struct {
	requestHistogram metrics.Histogram
	next             interfaces.Server
}

// GetQuote returns random quote from the World of Wisdom book.
func (mw *instrumentingMiddleware) GetQuote(startHash string, finishHash string) (string, error) {
	begin := time.Now()
	quote, err := mw.next.GetQuote(startHash, finishHash)
	labels := []string{
		"method", "GetQuote",
		"error", fmt.Sprint(err != nil),
	}
	mw.requestHistogram.With(labels...).Observe(time.Since(begin).Seconds())

	return quote, err
}

// MakeChallenge initialize puzzle that informs the client to complete a guided tour.
func (mw *instrumentingMiddleware) MakeChallenge() (client.ChallengeHeader, error) {
	begin := time.Now()
	challenge, err := mw.next.MakeChallenge()
	labels := []string{
		"method", "MakeChallenge",
		"error", fmt.Sprint(err != nil),
	}
	mw.requestHistogram.With(labels...).Observe(time.Since(begin).Seconds())

	return challenge, err
}
