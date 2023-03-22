package server

// Service contains shared business logic for the application
//
//go:generate moqc -stub -out mock/service.go -pkg mock . Service
type Service interface {
	// GetQuote returns random quote from the World of Wisdom book.
	GetQuote() string
}
