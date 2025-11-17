package docker

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullLatestImage(ctx context.Context, cli *client.Client, imageStr string) error {
	log.Println("ImMMM: ", imageStr)
	out, err := cli.ImagePull(ctx, imageStr, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	b, _ := io.ReadAll(out)
	log.Println(string(b))
	return nil
}

func DeleteOldDigest(ctx context.Context, cli *client.Client, img types.ImageInspect) error {
	log.Println("ENTERING DIGEST::::::::::::::::::::::::")
	log.Println("DIGESTS: ", img.RepoDigests)
	for i, repoDigest := range img.RepoDigests {
		if i == 0 {
			log.Println("OOOOOOOOOOOOOOOOOOOOOOOOOo")
			continue
		}
		_, err := cli.ImageRemove(ctx, repoDigest, image.RemoveOptions{})
		log.Println("Digest Removed: ", repoDigest)
		if err != nil {
			return err
		}
	}
	return nil
}

func RestartContainer(ctx context.Context, cli *client.Client, old types.ContainerJSON) (string, error) {

	timeout := 10
	log.Println("NOW Trying to restart")

	log.Println("starting restgart::::::::::")

	cli.ContainerStop(ctx, old.ID, container.StopOptions{Timeout: &timeout})

	cli.ContainerRemove(ctx, old.ID, container.RemoveOptions{
		Force: true,
	})

	resp, err := cli.ContainerCreate(
		ctx,
		old.Config,
		old.HostConfig,
		nil,
		nil,
		old.Name[1:], // remove leading slash
	)

	if err != nil {
		return "", err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	return resp.ID, err
}
