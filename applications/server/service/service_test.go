package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/server/interfaces"
	"github.com/donmikel/wow/applications/server/interfaces/mock"
)

func TestService_GetQuote(t *testing.T) {
	type fields struct {
		wowRepository interfaces.WOWRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				wowRepository: &mock.WOWRepository{
					GetQuoteFn: func(sel func(len int) int) string {
						return "test"
					},
				},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				wowRepository: tt.fields.wowRepository,
			}
			got := s.GetQuote()
			assert.Equal(t, tt.want, got)
		})
	}
}
