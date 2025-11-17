package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ListWatchedContainers returns all containers with label io.watcher.enable=true
func ListWatchedContainers(ctx context.Context, cli *client.Client) ([]types.Container, error) {
	args := filters.NewArgs()
	args.Add("label", "io.watcher.enable=true")

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: args,
		All:     true,
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func GetImageDigest(image types.ImageInspect) string {
	if len(image.RepoDigests) == 0 {
		return ""
	}
	return image.RepoDigests[0]
}
