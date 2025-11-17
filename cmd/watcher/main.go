package main

import (
	"context"
	"log"

	"github.com/aasourav/watchcontainer/internal/config"
	"github.com/aasourav/watchcontainer/internal/docker"
	"github.com/aasourav/watchcontainer/internal/watcher"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config:", err)
	}

	cli, err := docker.New()
	if err != nil {
		log.Fatal("docker:", err)
	}

	w := watcher.New(cfg, cli)
	w.Run(context.Background())
}
