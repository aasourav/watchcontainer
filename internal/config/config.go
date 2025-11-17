package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval     time.Duration `yaml:"interval"`
	SlackWebhook string        `yaml:"slack_webhook"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	if v := os.Getenv("SLACK_WEBHOOK"); v != "" {
		cfg.SlackWebhook = v
	}

	return &cfg, nil
}
