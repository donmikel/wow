package client

// ChallengeHeader represents an information to start guided tour puzzle solving.
type ChallengeHeader struct {
	// Hash is a start hash value.
	Hash string
	// TourLength is a guide tour length.
	TourLength int
	// Guides is an array of guide addresses needed to complete the tour.
	Guides []string
}
