package quotes

import (
	"net"
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/applications/server/mock"
	"github.com/donmikel/wow/pkg/protocol"
	server_mock "github.com/donmikel/wow/pkg/server/mock"
)

func TestTCPServer_GetQuote(t *testing.T) {
	type fields struct {
		service server.Service
	}
	tests := []struct {
		name          string
		fields        fields
		want          interface{}
		wantMessageID int
	}{
		{
			name: "get quote - success",
			fields: fields{
				service: &mock.Service{
					GetQuoteFn: func() string {
						return "test"
					},
				},
			},
			wantMessageID: protocol.GetQuoteResponse,
			want: &protocol.GetQuoteResponseMessage{
				Quote: "test",
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

			tcp := &TCPServer{
				service: tt.fields.service,
				logger:  log.NewNopLogger(),
			}
			tcp.GetQuote(c)

			want, err := c.Session().Codec().Encode(tt.want)
			assert.NoError(t, err)

			wantMsg := easytcp.NewMessage(tt.wantMessageID, want)
			assert.Equal(t, wantMsg, c.Response())
		})
	}
}
