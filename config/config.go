package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Worker   Worker   `yaml:"worker"`
	Streamer Streamer `yaml:"streamer"`
}

type Worker struct {
	Address string        `yaml:"address"`
	MaxJobs int           `yaml:"max-jobs"`
	Timeout time.Duration `yaml:"timeout"`
}

type Streamer struct {
	Address         string        `yaml:"address"`
	Interval        time.Duration `yaml:"interval"`
	ProcessTimeSeed int           `yaml:"process-time-seed"`
}

// New reads the config file, parses its content and returns a new
// instance of [*Config].
func New(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg *Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return cfg, nil
}
