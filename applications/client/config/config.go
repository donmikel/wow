package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Client contains all configuration settings related to client binary.
type Client struct {
	Ops    Ops `yaml:"ops"`
	Server TCP `yaml:"server"`
}

// Ops section describes settings for private operations-related API.
type Ops struct {
	// HTTPAddr is TCP address service's private operations-related API (e.g. exposing Prometheus metrics) listens on.
	HTTPAddr string `yaml:"http_addr"`
}

// TCP section describes settings for tcp server.
type TCP struct {
	// TCPAddr is TCP address service's API.
	TCPAddr string `yaml:"tcp_addr"`
}

// Validate validates some configuration settings to catch configuration errors early.
func (cfg *Client) Validate() error {
	return nil
}

// Parse YAML configuration file.
func Parse(filePath string) (Client, error) {
	c := Client{}

	f, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return c, fmt.Errorf("cannot read from file: %s, err: %w", filePath, err)
	}

	d := yaml.NewDecoder(f)
	// Because there are multiple consumers of the configuration being parsed here, it is much simpler to ignore extra fields
	// compared to adjusting configuration source (Helm, inventory folder) to contain only those fields we've defined in this code.
	d.SetStrict(false)

	err = d.Decode(&c)
	if err != nil {
		return c, fmt.Errorf("error decoding file: %s, err: %w", filePath, err)
	}

	err = f.Close()
	if err != nil {
		return c, fmt.Errorf("cannot close file: %s, err: %w", filePath, err)
	}

	return c, nil
}
