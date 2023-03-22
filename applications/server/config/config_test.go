package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	want := Server{
		Ops: Ops{HTTPAddr: "0.0.0.0:8002"},
		Guides: []Guide{
			{
				Host: "guide-1:8100",
				Key:  "key1",
			},
			{
				Host: "guide-2:8101",
				Key:  "key2",
			},
		},
		TCP: TCP{
			TCPAddr: ":8099",
		},
	}

	got, err := Parse("config.yml")

	assert.NoError(t, got.Validate())
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}
