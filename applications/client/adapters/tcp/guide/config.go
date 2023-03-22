package guide

import (
	"net"

	"github.com/DarthPestilane/easytcp"
	"github.com/go-kit/log"
)

// Provider provides an interface to interact with the Guide API.
type Provider struct {
	conn   net.Conn
	packer easytcp.Packer
	codec  easytcp.Codec
	logger log.Logger
}

// NewProvider creates new Provider instance.
func NewProvider(conn net.Conn, packer easytcp.Packer, codec easytcp.Codec, logger log.Logger) (*Provider, error) {
	p := &Provider{
		logger: logger,
		packer: packer,
		codec:  codec,
		conn:   conn,
	}

	return p, nil
}
