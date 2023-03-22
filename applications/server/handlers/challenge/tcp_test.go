package challenge

import (
	"net"
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/pkg/protocol"
	server_mock "github.com/donmikel/wow/pkg/server/mock"
)

func TestTCPServer_MakeChallenge(t *testing.T) {
	type fields struct {
		guides        []string
		guideKeys     []string
		GetTourLength func() int
		GetTimeStamp  func() int64
	}

	tests := []struct {
		name          string
		fields        fields
		want          interface{}
		wantMessageID int
	}{
		{
			name: "make new puzzle",
			fields: fields{
				guides: []string{
					"test1",
					"test2",
				},
				guideKeys: []string{
					"key1",
					"key2",
				},
				GetTimeStamp: func() int64 {
					return 10
				},
				GetTourLength: func() int {
					return 5
				},
			},
			wantMessageID: protocol.ChallengeRequiredResponse,
			want: &protocol.POWChallengeRequiredResponse{
				Hash:       "f14d9f03a63d21ce32f2acaa6b82cd8eaafee11402c24190851259bc89714393",
				TourLength: 5,
				Guides: []string{
					"test1",
					"test2",
				},
			},
		},
		{
			name: "empty guide list",
			fields: fields{
				guides: []string{},
				GetTimeStamp: func() int64 {
					return 10
				},
				GetTourLength: func() int {
					return 5
				},
			},
			wantMessageID: protocol.Error,
			want: &protocol.ErrorResponse{
				Code:    server.ErrBadRequest,
				Message: "guide keys list is empty",
			},
		},
	}

	for _, tt := range tests {
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
		t.Run(tt.name, func(t *testing.T) {
			tcp := &TCPServer{
				guides:        tt.fields.guides,
				guideKeys:     tt.fields.guideKeys,
				GetTourLength: tt.fields.GetTourLength,
				GetTimeStamp:  tt.fields.GetTimeStamp,
				logger:        log.NewNopLogger(),
			}

			tcp.MakeChallenge(c)
			want, err := c.Session().Codec().Encode(tt.want)
			assert.NoError(t, err)

			wantMsg := easytcp.NewMessage(tt.wantMessageID, want)
			assert.Equal(t, wantMsg, c.Response())
		})
	}
}
