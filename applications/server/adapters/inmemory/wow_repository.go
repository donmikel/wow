package inmemory

// WOWRepository is an inmemory storage for guide.
type WOWRepository struct {
	storage []string
}

// NewWOWRepository creates new instance of WOWRepository.
func NewWOWRepository(quotes []string) *WOWRepository {
	return &WOWRepository{
		storage: quotes,
	}
}

// GetQuote returns a random quote from the list.
func (w *WOWRepository) GetQuote(sel func(len int) int) string {
	qLen := len(w.storage)
	if qLen == 0 {
		return ""
	}
	index := sel(qLen)

	return w.storage[index]
}
