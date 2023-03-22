package protocol

// Header of TCP-message in protocol.
const (
	Error                     = iota
	GuideRequest              // from client to guide request a new step to solve the puzzle.
	GuideResponse             // from guide to client message with challenge for client
	ChallengeRequest          // from client to server request to star tour guide.
	ChallengeRequiredResponse // from server to client response with puzzle challenge.
	GetQuoteRequest           // from client to server - message with solved challenge
	GetQuoteResponse          // from server to client - message with useful info is solution is correct, or with error if not
)
