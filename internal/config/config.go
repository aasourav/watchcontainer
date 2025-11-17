package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	WatchInterval   time.Duration
	SlackWebhook    string
	SlackChannel    string
	IsCleanOldImage bool
}

func Load() (*Config, error) {
	var cfg Config
	cfg.IsCleanOldImage = false
	cfg.SlackChannel = ""
	cfg.SlackWebhook = ""
	cfg.WatchInterval = 10

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
