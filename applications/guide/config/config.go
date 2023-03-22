package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Guide contains all configuration settings related to guide binary.
type Guide struct {
	Ops Ops `yaml:"ops"`
	// Key is a secret key shared with the server.
	Key string `yaml:"key"`
	TCP TCP    `yaml:"api"`
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
func (cfg *Guide) Validate() error {
	if cfg.Key == "" {
		return fmt.Errorf("empty cfg.Key")
	}

	return nil
}

// Parse YAML configuration file.
func Parse(filePath string) (Guide, error) {
	c := Guide{}

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
