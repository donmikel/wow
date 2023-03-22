package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	want := Client{
		Ops: Ops{HTTPAddr: "0.0.0.0:8005"},
		Server: TCP{
			TCPAddr: "server:8099",
		},
	}

	got, err := Parse("config.yml")

	assert.NoError(t, got.Validate())
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}
