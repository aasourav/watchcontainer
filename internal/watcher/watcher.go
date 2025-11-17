package watcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"

	"github.com/aasourav/watchcontainer/internal/config"
	"github.com/aasourav/watchcontainer/internal/docker"
	"github.com/aasourav/watchcontainer/internal/notify"
)

type Watcher struct {
	Cfg *config.Config
	Cli *client.Client
}

func New(cfg *config.Config, cli *client.Client) *Watcher {
	return &Watcher{Cfg: cfg, Cli: cli}
}

func (w *Watcher) Run(ctx context.Context) {
	for {
		w.scan(ctx)
		time.Sleep(w.Cfg.Interval)
	}
}

func (w *Watcher) scan(ctx context.Context) {
	containers, err := docker.ListWatchedContainers(ctx, w.Cli)
	if err != nil {
		log.Println("List error:", err)
		return
	}

	for _, c := range containers {
		log.Println(c.Names)
		w.process(ctx, c)
	}
}

func (w *Watcher) process(ctx context.Context, c types.Container) {
	oldDetails, err := w.Cli.ContainerInspect(ctx, c.ID)
	if err != nil {
		log.Println("inspect err:", err)
		return
	}

	oldImg, _, err := w.Cli.ImageInspectWithRaw(ctx, c.Image)
	if err != nil {
		log.Println("local image err:", err)
		return
	}
	log.Println(oldImg.ID)

	oldDigest := docker.GetImageDigest(oldImg)

	// pull latest
	if err := docker.PullLatestImage(ctx, w.Cli, strings.Split(oldDigest, "@")[0]); err != nil {
		log.Println("pull err:", err)
		return
	}

	newImg, _, err := w.Cli.ImageInspectWithRaw(ctx, c.Image)
	if err != nil {
		log.Println("new image err:", err)
		return
	}
	return

	newDigest := docker.GetImageDigest(newImg)

	if newDigest == oldDigest {
		log.Println("SAME IMAGE")
		return
	}
	log.Println("Diffrent IMAGE")
	// update now
	_, err = docker.RestartContainer(ctx, w.Cli, oldDetails)
	if err != nil {
		log.Println("restart err:", err)
		return
	}

	docker.DeleteOldDigest(ctx, w.Cli, oldImg)

	labels := oldDetails.Config.Labels
	if labels["io.watcher.slack.enable"] != "true" {
		return
	}

	// choose webhook (container override > global)
	webhook := labels["io.watcher.slack.webhook"]
	if webhook == "" {
		webhook = w.Cfg.SlackWebhook
	}

	channel := labels["io.watcher.slack.channel"]

	msg := fmt.Sprintf(
		"*Container Updated*\nName: `%s`\nOld Hash: `%s`\nNew Hash: `%s`",
		c.Names[0], oldDigest, newDigest,
	)

	if webhook != "" {
		notify.Slack(webhook, msg, channel)
	}
}
