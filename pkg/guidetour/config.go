package guidetour

// GuideTour provides logic for the Guided tour puzzle protocol.
type GuideTour struct {
	Address string
	// GetTourLength returns length of the guide tour.
	// We can implement more complex logic here to adjust the length of the tour depending on
	// the load or suspiciousness of the client
	GetTourLength func() int
	GetTimeStamp  func() int64
}

// NewGuideTour creates a new instance of GuideTour.
func NewGuideTour(address string, getTourLength func() int, getTimeStamp func() int64) *GuideTour {
	return &GuideTour{
		Address:       address,
		GetTourLength: getTourLength,
		GetTimeStamp:  getTimeStamp,
	}
}
