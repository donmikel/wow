package service

import (
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/client/interfaces"
	"github.com/donmikel/wow/applications/client/interfaces/mock"
)

func TestService_GetQuote(t *testing.T) {
	type fields struct {
		server interfaces.Server
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				server: &mock.Server{
					GetQuoteFn: func(startHash string, finishHash string) (string, error) {
						return "test", nil
					},
				},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(
				tt.fields.server,
				easytcp.NewDefaultPacker(),
				&easytcp.JsonCodec{},
				log.NewNopLogger(),
			)
			got, err := s.GetQuote()
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
