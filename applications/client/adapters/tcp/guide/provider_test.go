package guide

import (
	"bytes"
	"testing"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"

	"github.com/donmikel/wow/pkg/protocol"
	"github.com/donmikel/wow/pkg/server/mock"
)

func TestProvider_GetTourHash(t *testing.T) {
	type fields struct {
		hash       string
		stepNumber int
		tourLength int
	}
	tests := []struct {
		name      string
		fields    fields
		wantMsgID int
		wantMsg   string
		wantErr   bool
	}{
		{
			name:      "get guide hash - success",
			wantMsgID: protocol.GuideResponse,
			wantMsg:   "test",
		},
		{
			name:      "bad response",
			wantMsgID: protocol.GetQuoteResponse,
			wantMsg:   "wrong meesage received, msg_id: 7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			en, _ := (&easytcp.JsonCodec{}).Encode(protocol.GuideResponseMessage{Hash: tt.wantMsg})
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

			got, err := p.GetTourHash(tt.fields.hash, tt.fields.stepNumber, tt.fields.tourLength)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, got, tt.wantMsg)
		})
	}
}
