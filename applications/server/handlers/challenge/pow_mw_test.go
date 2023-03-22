package challenge

import (
	"net"
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/server"
	"github.com/donmikel/wow/pkg/protocol"
	"github.com/donmikel/wow/pkg/server/mock"
)

func TestPOW_GuideTour(t *testing.T) {
	type fields struct {
		guides        []string
		guideKeys     []string
		GetTourLength func() int
		GetTimeStamp  func() int64
	}

	tests := []struct {
		name           string
		fields         fields
		requestID      int
		requestMessage interface{}
		want           interface{}
		wantMessageID  int
	}{
		{
			name:           "request without challenge headed - failed",
			requestID:      protocol.GetQuoteRequest,
			requestMessage: nil,
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
			wantMessageID: protocol.Error,
			want: &protocol.ErrorResponse{
				Code:    server.ErrBadRequest,
				Message: "both hashes are required",
			},
		},
		{
			name:      "request with valid challenge headed",
			requestID: protocol.GetQuoteRequest,
			requestMessage: &protocol.GetQuoteRequestMessage{
				GuideTourPuzzleHeader: protocol.GuideTourPuzzleHeader{
					StartHash:  "30965d9592d652e45d2583dce1cd07463136b22fe47f58aa92c43832d615a071",
					FinishHash: "6161d662b589e5b3ad0a74b5bc6d3f460fd2e11bfd096e1f3eb0a725100c4073",
				},
			},
			fields: fields{
				guides: []string{
					"test1",
					"test2",
				},
				guideKeys: []string{
					"1111",
					"2222",
				},
				GetTimeStamp: func() int64 {
					return 10
				},
				GetTourLength: func() int {
					return 5
				},
			},
			wantMessageID: protocol.GetQuoteResponse,
			want: &protocol.GetQuoteResponseMessage{
				Quote: "test",
			},
		},
		{
			name:      "request with outdated challenge headed",
			requestID: protocol.GetQuoteRequest,
			requestMessage: &protocol.GetQuoteRequestMessage{
				GuideTourPuzzleHeader: protocol.GuideTourPuzzleHeader{
					StartHash:  "30965d9592d652e45d2583dce1cd07463136b22fe47f58aa92c43832d615a071",
					FinishHash: "6161d662b589e5b3ad0a74b5bc6d3f460fd2e11bfd096e1f3eb0a725100c4073",
				},
			},
			fields: fields{
				guides: []string{
					"test1",
					"test2",
				},
				guideKeys: []string{
					"1111",
					"2222",
				},
				GetTimeStamp: func() int64 {
					return 11
				},
				GetTourLength: func() int {
					return 5
				},
			},
			wantMessageID: protocol.Error,
			want: &protocol.ErrorResponse{
				Code:    server.ErrBadRequest,
				Message: "puzzle is invalid",
			},
		},
	}

	for _, tt := range tests {
		c := easytcp.NewContext()
		c.SetSession(&mock.Session{
			ConnFn: func() net.Conn {
				return &mock.Conn{
					RemoteAddrFn: func() net.Addr {
						return &mock.Addr{
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
			c.SetRequest(tt.requestID, tt.requestMessage)

			mw := NewPOWMiddleware(tt.fields.guideKeys, tt.fields.GetTourLength, tt.fields.GetTimeStamp, log.NewNopLogger())

			mw.GuideTour(func(ctx easytcp.Context) {
				ctx.SetResponse(protocol.GetQuoteResponse, &protocol.GetQuoteResponseMessage{Quote: "test"})
			})(c)

			want, err := c.Session().Codec().Encode(tt.want)
			assert.NoError(t, err)

			wantMsg := easytcp.NewMessage(tt.wantMessageID, want)
			assert.Equal(t, wantMsg, c.Response())
		})
	}
}
