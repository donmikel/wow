package interfaces

// WOWRepository repository of stored guide.
//
//go:generate moqc -stub -out mock/wow_repository.go -pkg mock . WOWRepository
type WOWRepository interface {
	// GetQuote returns a random quote from the list.
	GetQuote(sel func(len int) int) string
}
