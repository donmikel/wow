package server

import (
	"fmt"

	"github.com/DarthPestilane/easytcp"

	"github.com/donmikel/wow/applications/client"
	"github.com/donmikel/wow/pkg/protocol"
)

// GetQuote returns random quote from the World of Wisdom book.
func (p *Provider) GetQuote(startHash string, finishHash string) (string, error) {
	req := &protocol.GetQuoteRequestMessage{
		GuideTourPuzzleHeader: protocol.GuideTourPuzzleHeader{
			StartHash:  startHash,
			FinishHash: finishHash,
		},
	}
	data, err := p.codec.Encode(req)
	if err != nil {
		return "", fmt.Errorf("request encoding error: %w", err)
	}

	msg := easytcp.NewMessage(protocol.GetQuoteRequest, data)

	b, err := p.packer.Pack(msg)
	if err != nil {
		return "", fmt.Errorf("request packing error: %w", err)
	}

	if _, err := p.conn.Write(b); err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}

	return p.readGetQuoteResponse()
}

// MakeChallenge initialize puzzle that informs the client to complete a guided tour.
func (p *Provider) MakeChallenge() (client.ChallengeHeader, error) {
	msg := easytcp.NewMessage(protocol.ChallengeRequest, nil)

	b, err := p.packer.Pack(msg)
	if err != nil {
		return client.ChallengeHeader{}, fmt.Errorf("request packing error: %w", err)
	}

	if _, err := p.conn.Write(b); err != nil {
		return client.ChallengeHeader{}, fmt.Errorf("connection error: %w", err)
	}

	return p.readChallengeResponse()
}

func (p *Provider) readChallengeResponse() (client.ChallengeHeader, error) {
	msg, err := p.packer.Unpack(p.conn)
	if err != nil {
		return client.ChallengeHeader{}, fmt.Errorf("unpack response error: %w", err)
	}

	if msg.ID() == protocol.ChallengeRequiredResponse {
		var getChallengeResp protocol.POWChallengeRequiredResponse
		if err = p.codec.Decode(msg.Data(), &getChallengeResp); err != nil {
			return client.ChallengeHeader{}, fmt.Errorf("decode VerifyChallengeResponseMessage error: %w", err)
		}

		return client.ChallengeHeader{
			Hash:       getChallengeResp.Hash,
			TourLength: getChallengeResp.TourLength,
			Guides:     getChallengeResp.Guides,
		}, nil
	}

	return client.ChallengeHeader{}, fmt.Errorf("wrong meesage received, msg_id: %v", msg.ID())
}

func (p *Provider) readGetQuoteResponse() (string, error) {
	msg, err := p.packer.Unpack(p.conn)
	if err != nil {
		return "", fmt.Errorf("unpack response error: %w", err)
	}

	switch msg.ID() {
	case protocol.GetQuoteResponse:
		var getQuoteResp protocol.GetQuoteResponseMessage
		if err = p.codec.Decode(msg.Data(), &getQuoteResp); err != nil {
			return "", fmt.Errorf("decode GetQuoteResponseMessage error: %w", err)
		}

		return getQuoteResp.Quote, nil
	case protocol.Error:
		var errorResp protocol.ErrorResponse
		err = p.codec.Decode(msg.Data(), &errorResp)
		if err != nil {
			return "", fmt.Errorf("decode POWChallengeRequiredResponse error: %w", err)
		}

		return "", fmt.Errorf("readGetQuoteResponse received error: %s", errorResp.Message)
	}

	return "", fmt.Errorf("wrong meesage received, msg_id: %v", msg.ID())
}
