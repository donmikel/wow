package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	want := Guide{
		Ops: Ops{HTTPAddr: "0.0.0.0:8003"},
		Key: "key1",
		TCP: TCP{
			TCPAddr: ":8100",
		},
	}

	got, err := Parse("config.yml")

	assert.NoError(t, got.Validate())
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}
