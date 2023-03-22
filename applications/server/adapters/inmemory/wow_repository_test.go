package inmemory

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint:gosec
func TestWOWRepository_GetQuote(t *testing.T) {
	type fields struct {
		storage []string
	}
	tests := []struct {
		name   string
		fields fields
		wants  []string
	}{
		{
			name: "get first random quote",
			fields: fields{
				storage: []string{"test", "test2"},
			},
			wants: []string{
				"test2",
			},
		},
		{
			name: "get random quote two times",
			fields: fields{
				storage: []string{
					"test",
					"test2",
					"test3",
					"test4",
				},
			},
			wants: []string{
				"test2",
				"test4",
			},
		},
	}
	for _, tt := range tests {
		s := rand.NewSource(1)
		r := rand.New(s)

		t.Run(tt.name, func(t *testing.T) {
			w := &WOWRepository{
				storage: tt.fields.storage,
			}

			for i, want := range tt.wants {
				got := w.GetQuote(func(len int) int {
					return r.Intn(len)
				})
				assert.Equal(t, got, want, "iteration %d", i)
			}
		})
	}
}
