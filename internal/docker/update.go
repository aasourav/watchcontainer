package docker

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullLatestImage(ctx context.Context, cli *client.Client, imageStr string) error {
	out, err := cli.ImagePull(ctx, imageStr, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.ReadAll(out)
	return nil
}

func DeleteOldDigest(ctx context.Context, cli *client.Client, imageID string, imageName string) error {
	args := filters.NewArgs()
	args.Add("dangling", "true")

	images, err := cli.ImageList(ctx, image.ListOptions{
		Filters: args,
	})

	if err != nil {
		log.Fatalf("failed to list images: %v", err)
	}

	if len(images) == 0 {
		fmt.Println("no untagged images found.")
		return nil
	}

	for _, img := range images {
		if imageID == img.ID {
			_, err := cli.ImageRemove(ctx, img.ID, image.RemoveOptions{
				Force:         true,
				PruneChildren: true,
			})

			if err != nil {
				fmt.Printf("failed to remove image %s: %v\n", imageName, err)
			} else {
				fmt.Printf("removing untagged image: %s\n", imageName)
			}
		}
	}

	return nil
}

func RestartContainer(ctx context.Context, cli *client.Client, old types.ContainerJSON) (string, error) {

	timeout := 10
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
