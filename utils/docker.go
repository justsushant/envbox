package utils

import (
	"fmt"
	"github.com/docker/docker/client"
)

func NewDockerClient() (*client.Client, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while creating docker client: %v", err)
	}

	return client, nil
}