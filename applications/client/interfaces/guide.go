package interfaces

// Guide provides an interface to interact with the Guide API.
//
//go:generate moqc -stub -out mock/guide.go -pkg mock . Guide
type Guide interface {
	// GetTourHash solves puzzle of the tour guide at the l-th stop.
	GetTourHash(hash string, stepNumber int, tourLength int) (hi string, err error)
}
