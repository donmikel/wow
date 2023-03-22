package server

import (
	"bytes"
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/applications/client"
	"github.com/donmikel/wow/pkg/protocol"
	"github.com/donmikel/wow/pkg/server/mock"
)

func TestProvider_GetQuote(t *testing.T) {
	type fields struct {
		startHash  string
		finishHash string
	}
	tests := []struct {
		name      string
		fields    fields
		wantMsgID int
		want      string
		wantErr   bool
	}{
		{
			name:      "get quote - success",
			wantMsgID: protocol.GetQuoteResponse,
			want:      "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			en, _ := (&easytcp.JsonCodec{}).Encode(protocol.GetQuoteResponseMessage{Quote: tt.want})
			bt, _ := easytcp.NewDefaultPacker().Pack(easytcp.NewMessage(tt.wantMsgID, en))
			buf := bytes.NewBuffer(bt)

			p := &Provider{
				conn: &mock.Conn{
					ReadFn: func(b []byte) (int, error) {
						return buf.Read(b)
					},
				},
				packer: easytcp.NewDefaultPacker(),
				codec:  &easytcp.JsonCodec{},
				logger: log.NewNopLogger(),
			}

			got, err := p.GetQuote(tt.fields.startHash, tt.fields.finishHash)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, got, tt.want)
		})
	}
}

func TestProvider_MakeChallenge(t *testing.T) {
	type fields struct {
		hash       string
		tourLength int
		guides     []string
	}
	tests := []struct {
		name      string
		fields    fields
		wantMsgID int
		want      client.ChallengeHeader
		wantErr   bool
	}{
		{
			name:      "make challenge - success",
			wantMsgID: protocol.ChallengeRequiredResponse,
			fields: fields{
				hash:       "hash",
				tourLength: 10,
				guides:     []string{"test"},
			},
			want: client.ChallengeHeader{
				Hash:       "hash",
				TourLength: 10,
				Guides:     []string{"test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			en, _ := (&easytcp.JsonCodec{}).Encode(protocol.POWChallengeRequiredResponse{
				Hash:       tt.fields.hash,
				TourLength: tt.fields.tourLength,
				Guides:     tt.fields.guides,
			})
			bt, _ := easytcp.NewDefaultPacker().Pack(easytcp.NewMessage(tt.wantMsgID, en))
			buf := bytes.NewBuffer(bt)

			p := &Provider{
				conn: &mock.Conn{
					ReadFn: func(b []byte) (int, error) {
						return buf.Read(b)
					},
				},
				packer: easytcp.NewDefaultPacker(),
				codec:  &easytcp.JsonCodec{},
				logger: log.NewNopLogger(),
			}

			got, err := p.MakeChallenge()
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, got, tt.want)
		})
	}
}
