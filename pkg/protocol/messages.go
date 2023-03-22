package protocol

// POWChallengeRequiredResponse represents a message format for the ChallengeRequiredResponse id.
type POWChallengeRequiredResponse struct {
	Hash       string   `json:"hash"`
	TourLength int      `json:"tour_length"`
	Guides     []string `json:"guides"`
}

// GuideTourPuzzleHeader represents a challenge header.
type GuideTourPuzzleHeader struct {
	StartHash  string `json:"start_hash"`
	FinishHash string `json:"finish_hash"`
}

// GetQuoteResponseMessage represents a message format for the GetQuoteResponse id.
type GetQuoteResponseMessage struct {
	Quote string
}

// GetQuoteRequestMessage represents a message format for the GetQuoteRequest id.
type GetQuoteRequestMessage struct {
	GuideTourPuzzleHeader
}

// GuideRequestMessage represents a message format for the GuideRequest id.
type GuideRequestMessage struct {
	Hash       string `yaml:"hash"`
	StepNumber int    `yaml:"step_number"`
	TourLength int    `yaml:"tour_length"`
}

// GuideResponseMessage represents a message format for the GuideResponse id.
type GuideResponseMessage struct {
	Hash string `yaml:"hash"`
}

// ErrorResponse represents a message format for the Error id.
type ErrorResponse struct {
	// Code is a machine-readable code.
	Code string `json:"code,omitempty"`
	// Message is a human-readable message.
	Message string `json:"message"`
	// Inner is a wrapped error that is never shown to API consumers.
	Inner error `json:"-"`
}
