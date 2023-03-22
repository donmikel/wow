package guide

import (
	"net"
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/guide"
	"github.com/donmikel/wow/pkg/protocol"
	server_mock "github.com/donmikel/wow/pkg/server/mock"
)

func TestTCPServer_GetTourHash(t *testing.T) {
	type fields struct {
		GetTourLength func() int
		GetTimeStamp  func() int64
		guideKey      string
	}
	tests := []struct {
		name          string
		fields        fields
		requestID     int
		request       interface{}
		wantMessage   interface{}
		wantMessageID int
	}{
		{
			name: "get guide hash - success",
			fields: fields{
				GetTimeStamp: func() int64 {
					return 10
				},
				GetTourLength: func() int {
					return 5
				},
				guideKey: "key",
			},
			requestID: protocol.GuideRequest,
			request: &protocol.GuideRequestMessage{
				Hash:       "hash",
				StepNumber: 5,
				TourLength: 10,
			},
			wantMessageID: protocol.GuideResponse,
			wantMessage: &protocol.GuideResponseMessage{
				Hash: "1a0a1186e316bb62471fa4a39ba687e338d8b25ab51d19c919631ac007401771",
			},
		},
		{
			name: "wrong guide step - fail",
			fields: fields{
				GetTimeStamp: func() int64 {
					return 10
				},
				GetTourLength: func() int {
					return 5
				},
				guideKey: "key",
			},
			requestID: protocol.GuideRequest,
			request: &protocol.GuideRequestMessage{
				Hash:       "hash",
				StepNumber: 11,
				TourLength: 10,
			},
			wantMessageID: protocol.Error,
			wantMessage: &protocol.ErrorResponse{
				Code:    guide.ErrBadRequest,
				Message: "step number must be less or equal to the length of the tour",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := easytcp.NewContext()
			c.SetSession(&server_mock.Session{
				ConnFn: func() net.Conn {
					return &server_mock.Conn{
						RemoteAddrFn: func() net.Addr {
							return &server_mock.Addr{
								StringFn: func() string {
									return "127.0.0.1"
								},
							}
						},
					}
				},
				IDFn: func() interface{} {
					return "20"
				},
				CodecFn: func() easytcp.Codec {
					return &easytcp.JsonCodec{}
				},
			})

			tcp := NewTCPServer(tt.fields.guideKey, tt.fields.GetTourLength, tt.fields.GetTimeStamp, log.NewNopLogger())
			err := c.SetRequest(tt.requestID, tt.request)
			assert.NoError(t, err)

			tcp.GetTourHash(c)

			want, err := c.Session().Codec().Encode(tt.wantMessage)
			assert.NoError(t, err)

			wantMsg := easytcp.NewMessage(tt.wantMessageID, want)
			assert.Equal(t, wantMsg, c.Response())
		})
	}
}
