package guidetour

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuideTour_MakePuzzle(t *testing.T) {
	type fields struct {
		Address       string
		getTimeStamp  int64
		getTourLength int
	}
	type args struct {
		Ks string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "make puzzle",
			fields: fields{
				Address:       "127.0.0.1",
				getTourLength: 5,
				getTimeStamp:  10,
			},
			args: args{
				Ks: "short-lived_secret",
			},
			want: "30965d9592d652e45d2583dce1cd07463136b22fe47f58aa92c43832d615a071",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuideTour{
				Address: tt.fields.Address,
				GetTourLength: func() int {
					return tt.fields.getTourLength
				},
				GetTimeStamp: func() int64 {
					return tt.fields.getTimeStamp
				},
			}
			got := g.MakePuzzle(tt.args.Ks)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGuideTour_RunTour(t *testing.T) {
	type fields struct {
		Address       string
		getTimeStamp  int64
		getTourLength int
	}
	type args struct {
		h0         string
		tourLength int
		keys       []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "run and verify",
			fields: fields{
				Address:       "127.0.0.1",
				getTourLength: 5,
				getTimeStamp:  10,
			},
			args: args{
				h0:         "30965d9592d652e45d2583dce1cd07463136b22fe47f58aa92c43832d615a071",
				tourLength: 5,
				keys: []string{
					"1111",
					"2222",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuideTour{
				Address: tt.fields.Address,
				GetTourLength: func() int {
					return tt.fields.getTourLength
				},
				GetTimeStamp: func() int64 {
					return tt.fields.getTimeStamp
				},
			}
			got, err := RunTour(tt.args.h0, len(tt.args.keys), tt.fields.getTourLength, func(h0 string, stepNumber int, tourIndex int, tourLength int) (string, error) {
				hi := g.TourGuideHash(h0, stepNumber, tourLength, tt.args.keys[tourIndex])
				return hi, nil
			})
			assert.NoError(t, err)

			v, err := g.VerifyPuzzle(tt.args.h0, got, tt.args.keys)
			fmt.Println(got)

			assert.NoError(t, err)
			assert.True(t, v)
		})
	}
}
