package guidetour

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// MakePuzzle initialize puzzle that informs the client to complete a guided tour.
func (g *GuideTour) MakePuzzle(Ks string) string {
	data := fmt.Sprintf("%s:%d:%d:%s", g.getClientID(), g.GetTourLength(), g.GetTimeStamp(), Ks)

	return hash(data)
}

// RunTour is a process of solving guided tour puzzle.
func RunTour(
	h0 string,
	guidesCount int,
	tourLength int,
	sendToTG func(h0 string, stepNumber int, tourIndex int, tourLength int) (string, error),
) (string, error) {
	hi := h0
	L := tourLength
	var err error

	for i := 0; i < L; i++ {
		tourIndex := getTourIndex(hi, guidesCount)
		hi, err = sendToTG(hi, i+1, tourIndex, L)
		if err != nil {
			return "", fmt.Errorf("send to guide error: %w", err)
		}
	}

	return hi, nil
}

func getTourIndex(hash string, n int) int {
	i := new(big.Int)
	result := new(big.Int)
	num := big.NewInt(int64(n))
	i.SetString(hash, 16)
	return int(result.Mod(i, num).Int64())
}

func (g *GuideTour) getClientID() string {
	return g.Address
}

// VerifyPuzzle verifies a result of guide tour puzzle solving.
func (g *GuideTour) VerifyPuzzle(firstHash string, lastHash string, keys []string) (bool, error) {
	hi := firstHash
	N := len(keys)
	L := g.GetTourLength()

	for i := 0; i < L; i++ {
		tourIndex := getTourIndex(hi, N)
		if tourIndex >= len(keys) {
			return false, fmt.Errorf("shared key for the guide with index equals %d, not found", tourIndex)
		}
		key := keys[tourIndex]

		hi = g.TourGuideHash(hi, i+1, L, key)
	}
	if hi == lastHash {
		return true, nil
	}

	return false, nil
}

// TourGuideHash solves puzzle of the tour guide at the l-th stop.
func (g *GuideTour) TourGuideHash(previousHash string, stepNumber int, length int, key string) string {
	data := fmt.Sprintf("%s:%d:%d:%s:%d:%s", previousHash, stepNumber, length, g.getClientID(), g.GetTimeStamp(), key)

	return hash(data)
}

func hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))

	return hex.EncodeToString(hasher.Sum(nil))
}
