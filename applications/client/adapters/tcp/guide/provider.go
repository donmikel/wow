package guide

import (
	"fmt"

	"github.com/DarthPestilane/easytcp"

	"github.com/donmikel/wow/pkg/protocol"
)

// GetTourHash solves puzzle of the tour guide at the l-th stop.
func (p *Provider) GetTourHash(hash string, stepNumber int, tourLength int) (string, error) {
	req := &protocol.GuideRequestMessage{
		Hash:       hash,
		StepNumber: stepNumber,
		TourLength: tourLength,
	}
	data, err := p.codec.Encode(req)
	if err != nil {
		return "", fmt.Errorf("request encoding error: %w", err)
	}

	msg := easytcp.NewMessage(protocol.GuideRequest, data)

	b, err := p.packer.Pack(msg)
	if err != nil {
		return "", fmt.Errorf("request packing error: %w", err)
	}

	if _, err := p.conn.Write(b); err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}

	return p.readGuideResponse()
}

func (p *Provider) readGuideResponse() (string, error) {
	msg, err := p.packer.Unpack(p.conn)
	if err != nil {
		return "", fmt.Errorf("unpack response error: %w", err)
	}

	if msg.ID() == protocol.GuideResponse {
		var guideResp protocol.GuideResponseMessage
		if err = p.codec.Decode(msg.Data(), &guideResp); err != nil {
			return "", fmt.Errorf("decode GuideResponseMessage error: %w", err)
		}

		return guideResp.Hash, nil
	}

	return "", fmt.Errorf("wrong meesage received, msg_id: %v, msg: %s", msg.ID(), msg.Data())
}
