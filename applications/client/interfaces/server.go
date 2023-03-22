package interfaces

import "github.com/donmikel/wow/applications/client"

// Server provides an interface to interact with the Server API.
//
//go:generate moqc -stub -out mock/server.go -pkg mock . Server
type Server interface {
	// GetQuote returns random quote from the World of Wisdom book.
	GetQuote(startHash string, finishHash string) (string, error)
	// MakeChallenge initialize puzzle that informs the client to complete a guided tour.
	MakeChallenge() (client.ChallengeHeader, error)
}
