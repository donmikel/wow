package service

import (
	"fmt"
	"net"

	"github.com/DarthPestilane/easytcp"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	stdprom "github.com/prometheus/client_golang/prometheus"

	"github.com/donmikel/wow/applications/client/adapters/tcp/guide"
	"github.com/donmikel/wow/applications/client/interfaces"
	"github.com/donmikel/wow/pkg/guidetour"
)

// Service contains shared business logic for the application.
type Service struct {
	server interfaces.Server
	guides map[string]interfaces.Guide
	packer easytcp.Packer
	codec  easytcp.Codec
	logger log.Logger
}

// NewService creates a new instance of Service.
func NewService(
	server interfaces.Server,
	packer easytcp.Packer,
	codec easytcp.Codec,
	logger log.Logger,
) *Service {
	return &Service{
		server: server,
		guides: map[string]interfaces.Guide{},
		packer: packer,
		codec:  codec,
		logger: logger,
	}
}

// GetQuote returns random quote from the World of Wisdom book.
func (s *Service) GetQuote() (string, error) {
	// Before calling the GetQuote endpoint, we must complete a guide tour.
	startHash, finishHash, err := s.runChallenge()
	if err != nil {
		return "", fmt.Errorf("runChallenge error: %w", err)
	}

	return s.server.GetQuote(startHash, finishHash)
}

func (s *Service) runChallenge() (string, string, error) {
	// Retrieves init puzzle hash.
	challenge, err := s.server.MakeChallenge()
	if err != nil {
		return "", "", fmt.Errorf("start challenge error: %w", err)
	}

	// Updates connections to guide instances.
	err = s.updateGuideConnections(challenge.Guides)
	if err != nil {
		return "", "", fmt.Errorf("updating guide connection error: %w", err)
	}

	// Run a tour guide.
	finishHash, err := guidetour.RunTour(challenge.Hash, len(challenge.Guides), challenge.TourLength, func(h0 string, stepNumber int, tourIndex int, tourLength int) (string, error) {
		guideAddress := challenge.Guides[tourIndex]
		gd, found := s.guides[guideAddress]
		if !found {
			return "", fmt.Errorf("guide not found, address: %s", guideAddress)
		}
		// Gets next step hash from the guide.
		return gd.GetTourHash(h0, stepNumber, tourLength)
	})
	if err != nil {
		return "", "", fmt.Errorf("RunTour error: %w", err)
	}
	return challenge.Hash, finishHash, nil
}

func (s *Service) updateGuideConnections(guides []string) error {
	for _, connString := range guides {
		if _, ok := s.guides[connString]; !ok {
			provider, err := s.connectToGuide(connString)
			if err != nil {
				return fmt.Errorf("guide connection error: %w", err)
			}

			s.guides[connString] = provider
		}
	}

	return nil
}

func (s *Service) connectToGuide(addr string) (interfaces.Guide, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcp server connection error: %w", err)
	}

	var gp interfaces.Guide
	gp, err = guide.NewProvider(conn, s.packer, s.codec, s.logger)
	if err != nil {
		return nil, fmt.Errorf("guide provider error: %w", err)
	}

	guideHistMetric := kitprom.NewHistogramFrom(stdprom.HistogramOpts{
		Namespace:   "wow",
		Subsystem:   "client",
		Name:        "guide_duration_seconds",
		Help:        "Guide requests duration in seconds by method.",
		ConstLabels: stdprom.Labels{"guide": addr},
	}, []string{"method", "error"})

	gp = guide.NewInstrumentingMiddleware(guideHistMetric, gp)

	if err != nil {
		return nil, fmt.Errorf("guide provider creation error: %w", err)
	}

	return gp, nil
}
