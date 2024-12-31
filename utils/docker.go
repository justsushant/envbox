package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func NewDockerClient() (*client.Client, error) {
	client, err := client.NewClientWithOpts(client.WithHost("tcp://10.128.0.2:2375"), client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error while creating docker client: %v", err)
	}

	return client, nil
}

func ReloadNginxConf() error {
	ctx := context.Background()

	cli, err := NewDockerClient()
	if err != nil {
		return fmt.Errorf("error while creating docker client: %v", err)
	}

	resp, err := cli.ContainerExecCreate(ctx, "nginx-reverse-proxy", types.ExecConfig{
		Cmd: []string{"nginx", "-s", "reload"},
	})
	if err != nil {
		return fmt.Errorf("error while creating exec: %v", err)
	}

	err = cli.ContainerExecStart(ctx, resp.ID, container.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("error while starting exec: %v", err)
	}

	return nil
}
