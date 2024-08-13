package utils

import (
	"github.com/docker/docker/client"
)

func NewDockerClient() *client.Client {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic("error while creating docker client: " + err.Error())
	}

	return client
}