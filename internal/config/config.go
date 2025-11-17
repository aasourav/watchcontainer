package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WatchInterval   time.Duration `yaml:"watch_interval"`
	SlackWebhook    string        `yaml:"slack_global_webhook"`
	SlackChannel    string        `yaml:"slack_global_channel"`
	IsCleanOldImage bool          `yaml:"clean_old_image"`
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

	if v := os.Getenv("GLOBAL_SLACK_WEBHOOK"); v != "" {
		cfg.SlackWebhook = v
	}

	if v := os.Getenv("GLOBAL_SLACK_CHANNEL"); v != "" {
		cfg.SlackChannel = v
	}

	if v := os.Getenv("WATCH_INTERVAL"); v != "" {
		secs, err := strconv.Atoi(v)
		if err != nil {
			log.Println("invalid WATCH_INTERVAL: must be number", err.Error())
		}
		cfg.WatchInterval = time.Duration(secs) * time.Second
	}

	if v := os.Getenv("CLEAN_OLD_IMAGE"); v != "" {
		isClean, err := strconv.ParseBool(v)
		if err != nil {
			log.Println("invalid WATCH_INTERVAL. must be true or false", err.Error())
		}
		cfg.IsCleanOldImage = isClean
	}

	return &cfg, nil
}
