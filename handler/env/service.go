package env

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/justsushant/envbox/types"
	"github.com/justsushant/envbox/utils"
)

var containerAddr = fmt.Sprintf("http://127.0.0.1:%s/tree?token=", utils.DEFAULT_CONTAINER_PORT)

type Service struct {
	store types.EnvStore
}

func NewService(store types.EnvStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAllEnvs() ([]types.Env, error) {
	return s.store.GetAllEnvs()
}

func (s *Service) KillEnv(client *client.Client, id string) (string, error) {
	ctx := context.Background()

	// get the container details from db
	env, err := s.store.GetContainerByID(id)
	fmt.Println(env)
	if err != nil {
		return "", fmt.Errorf("failed to get the container details: %v", err)
	}
	if !env.Active {
		return "", fmt.Errorf("container is already stopped")
	}

	// currently, stops the container immediately
	// could use ContainerStop to gracefully stop the container 
	if err := client.ContainerKill(ctx, env.ContainerID, "SIGKILL"); err != nil {
		return "", fmt.Errorf("failed to kill the container: %v", err)
	}

	// removes the container from host
	if err := client.ContainerRemove(ctx, env.ContainerID, container.RemoveOptions{Force:true}); err != nil {
		return "", fmt.Errorf("failed to remove the container: %v", err)
	}

	// sets the active to false in the db
	if err := s.store.DeleteContainer(id); err != nil {
		return "", fmt.Errorf("failed to inactive the container in database: %v", err)
	}

	return "container stopped and removed successfully", nil
}

func (s *Service) CreateEnv(client *client.Client, p types.CreateEnvPayload) (string, error) {
	ctx := context.Background()

	// gets a random free port on the host
	hostPort, err := utils.GetRandomFreePort()
	if err != nil {
		return "", fmt.Errorf("failed to get a free port: %v", err)
	}
	
	// create port binding for the container
	portBinding, err := createPortBinding("tcp", hostPort, utils.DEFAULT_CONTAINER_PORT)
	if err != nil {
		return "", err
	}

	// random password of 8 characters for notebook
	notebookPwd := utils.GenerateRandomPassword(8) 

	// creates the container
	resp, err := client.ContainerCreate(ctx, &container.Config{
		Image: p.ImageName,
		Cmd:   generateJupyterNoteBookStartCommand(notebookPwd, utils.DEFAULT_CONTAINER_PORT),
		ExposedPorts: nat.PortSet{
			utils.DEFAULT_CONTAINER_PORT + "/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: portBinding,
	}, nil, nil, "")
	if err != nil {
		fmt.Println(err)
	}

	// starts the container
	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Println(err)
	}
	
	// save the container details in the db
	if err := s.store.SaveContainer(resp.ID, p.ImageName); err != nil {
		fmt.Println("error while saving container details in the database: ", err)
	}

	fmt.Println("container id: ", resp.ID)

	// gets the logs from the container
	// we parse it to know when the notebook service has started
	logsChan, errChan := getContainerLogs(client, resp.ID)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				fmt.Println(err)
				return "", err
			}
		case log := <-logsChan:
			if strings.Contains(log, containerAddr) {
				return fmt.Sprintf("http://127.0.0.1:%s/tree?token=%s",hostPort,notebookPwd), nil
			}
		}
	}
}


func createPortBinding(proto string, hostPort string, containerPort string) (nat.PortMap, error) {
	port, err := nat.NewPort(proto, containerPort)
	if err != nil {
		return nil, fmt.Errorf("unable to get the container port: %v", err)
	}

	return nat.PortMap{port: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}}}, nil
}


func getContainerLogs(cli *client.Client, containerID string) (<-chan string, <-chan error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}

	// channels for logs and errors
	errChan := make(chan error)
	logsChan := make(chan string)

	// goroutine to get the logs from the container
	go func() {
		responseBody, err := cli.ContainerLogs(context.Background(), containerID, options)
		errChan <- err
		defer responseBody.Close()

		// reads the logs to find the notebook url
		// if notebook url is found, means the notebook has started
		buf := make([]byte, 1024)
		for {
			n, err := responseBody.Read(buf)
			if err != nil {
				errChan <- err
			}
			logsChan <- string(buf[:n])
		}
	}()

	return logsChan, errChan
}

// generates the notebook start command for container creation
func generateJupyterNoteBookStartCommand(token, port string) []string {
	cmd := fmt.Sprintf("jupyter notebook --allow-root --ip 0.0.0.0 --NotebookApp.allow_origin=* --NotebookApp.token=%s --no-browser --port=%s", token, port)
	return strings.Split(cmd, " ")
}