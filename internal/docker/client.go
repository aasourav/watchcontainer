package docker

import (
	"context"

	"github.com/docker/docker/client"
)

func New() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(context.Background())
	return cli, nil
}
