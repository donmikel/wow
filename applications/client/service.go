package client

// Service contains shared business logic for the application
//
//go:generate moqc -stub -out mock/service.go -pkg mock . Service
type Service interface {
	GetQuote() (string, error)
}
