package service

import (
	"math/rand"

	"github.com/donmikel/wow/applications/server/interfaces"
)

// Service contains shared business logic for the application.
type Service struct {
	wowRepository interfaces.WOWRepository
}

// NewService creates a new instance of Service.
func NewService(wowRepository interfaces.WOWRepository) *Service {
	return &Service{
		wowRepository: wowRepository,
	}
}

// GetQuote returns random quote from the World of Wisdom book.
// nolint:gosec
func (s *Service) GetQuote() string {
	return s.wowRepository.GetQuote(func(len int) int {
		return rand.Intn(len)
	})
}
