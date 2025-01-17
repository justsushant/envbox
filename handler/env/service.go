package env

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/justsushant/envbox/config"
	localTypes "github.com/justsushant/envbox/types"
	"github.com/justsushant/envbox/utils"
)

type Service struct {
	dockerClient   *client.Client
	envStore       localTypes.EnvStore
	imgStore       localTypes.ImageStore
	nginxCfgWriter localTypes.NginxCfgWriterService
}

func NewService(dockerClient *client.Client, envStore localTypes.EnvStore, imgStore localTypes.ImageStore, nginxCfgWriter localTypes.NginxCfgWriterService) *Service {
	return &Service{
		dockerClient:   dockerClient,
		envStore:       envStore,
		imgStore:       imgStore,
		nginxCfgWriter: nginxCfgWriter,
	}
}

func (s *Service) GetAllEnvs() ([]localTypes.GetImageResponse, error) {
	return s.envStore.GetAllEnvs()
}

func (s *Service) CreateEnv(p localTypes.CreateEnvPayload) (string, string, string, error) {
	ctx := context.Background()

	// get image details from db
	image, err := s.imgStore.GetImageByID(p.ImageID)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get the image details: %v", err)
	}

	// gets a random free port on the host
	hostPort, err := utils.GetRandomFreePort()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get a free port: %v", err)
	}

	// create port binding for the container
	portBinding, err := createPortBinding("tcp", hostPort, utils.DEFAULT_CONTAINER_PORT)
	if err != nil {
		return "", "", "", err
	}

	cmd := []string{}
	internalLink := ""
	accessLink := ""

	// set the startup command and link according to the image
	switch image.Name {
	case "Jupyter Notebook":
		cmd = generateJupyterNoteBookStartCommand(hostPort)
		internalLink = generateJupyterNotebookInternalLink(utils.DEFAULT_CONTAINER_PORT, hostPort)
		accessLink = generateJupyterNotebookAccessLink(hostPort)
	case "VS Code (Go)":
		cmd = generateVScodeGoStartCommand()
		internalLink = generateVsCodeGoInternalLink(utils.DEFAULT_CONTAINER_PORT)
		accessLink = generateVsCodeGoAccessLink(hostPort)
	}

	// creates the container
	resp, err := s.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: image.Path,
		Cmd:   cmd,
		// Cmd:   generateJupyterNoteBookStartCommand(),
		ExposedPorts: nat.PortSet{
			utils.DEFAULT_CONTAINER_PORT + "/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: portBinding,
	}, nil, nil, "")
	if err != nil {
		return "", "", "", fmt.Errorf("error while creating the container: %v", err)
	}

	// starts the container
	if err := s.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", "", "", fmt.Errorf("error while starting the container: %v", err)
	}

	// save the container details in the db
	if err := s.envStore.SaveContainer(resp.ID, image.ID); err != nil {
		return "", "", "", fmt.Errorf("error while saving container details in the database: %v", err)
	}

	// gets the logs from the container
	// we parse it to know when the app has started inside container
	logsChan, errChan := getContainerLogs(s.dockerClient, resp.ID)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				return "", "", "", err
			}
		case log := <-logsChan:
			// if notebook url is found, means the notebook has started
			if strings.Contains(log, internalLink) {
				// accessLink := fmt.Sprintf("http://127.0.0.1:%s/tree", hostPort)

				if err := s.envStore.UpdateContainerAccessLink(resp.ID, accessLink); err != nil {
					return "", "", "", fmt.Errorf("error while updating the access link in the database: %v", err)
				}

				return hostPort, accessLink, resp.ID, nil
			}
		}
	}
}

func (s *Service) KillEnv(id string) (string, error) {
	ctx := context.Background()

	// get the container details from db
	env, err := s.envStore.GetContainerByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no container found: %v", err)
		}
		return "", fmt.Errorf("failed to get the container details: %v", err)
	}
	if !env.Active {
		return "", fmt.Errorf("container is already stopped")
	}

	// currently, stops the container immediately
	// could use ContainerStop to gracefully stop the container
	if err := s.dockerClient.ContainerKill(ctx, env.ContainerID, "SIGKILL"); err != nil {
		return "", fmt.Errorf("failed to kill the container: %v", err)
	}

	// removes the container from host
	if err := s.dockerClient.ContainerRemove(ctx, env.ContainerID, container.RemoveOptions{Force: true}); err != nil {
		return "", fmt.Errorf("failed to remove the container: %v", err)
	}

	// sets the active to false in the db
	if err := s.envStore.DeleteContainer(id); err != nil {
		return "", fmt.Errorf("failed to inactive the container in database: %v", err)
	}

	return env.ContainerID, nil
}

func (s *Service) GetTerminal(id string) (types.HijackedResponse, error) {
	// getting container data from store
	env, err := s.envStore.GetContainerByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.HijackedResponse{}, fmt.Errorf("no container found: %v", err)
		}
		return types.HijackedResponse{}, fmt.Errorf("failed to get the container details: %v", err)
	}
	if !env.Active {
		return types.HijackedResponse{}, fmt.Errorf("container is already stopped")
	}

	// creating exec for the container
	execID, err := s.dockerClient.ContainerExecCreate(context.Background(), env.ContainerID, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		return types.HijackedResponse{}, err
	}

	// extracting response object from exec
	hijackedResp, err := s.dockerClient.ContainerExecAttach(context.Background(), execID.ID, container.ExecStartOptions{Tty: true})
	if err != nil {
		return types.HijackedResponse{}, err
	}

	return hijackedResp, nil
}

func (s *Service) AddNginxUpstream(imageID int, containerID, name, addr string) error {
	var isRewrite bool

	if imageID == 1 {
		isRewrite = false
	} else if imageID == 2 {
		isRewrite = true
	}

	// store the details in db
	if err := s.nginxCfgWriter.AddUpstreamAndRender(containerID, name, addr, isRewrite); err != nil {
		return fmt.Errorf("failed to update the nginx config details: %v", err)
	}

	return nil
}

func (s *Service) RemoveNginxUpstream(containerID string) error {
	// remove the details from db
	if err := s.nginxCfgWriter.RemoveUpstreamAndRender(containerID); err != nil {
		return fmt.Errorf("failed to delete the nginx config details: %v", err)
	}

	return nil
}

func createPortBinding(proto string, hostPort string, containerPort string) (nat.PortMap, error) {
	port, err := nat.NewPort(proto, containerPort)
	if err != nil {
		return nil, fmt.Errorf("unable to get the container port: %v", err)
	}

	return nat.PortMap{port: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}}}, nil
}

func getContainerLogs(client *client.Client, containerID string) (<-chan string, <-chan error) {
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
		// closing the channels
		defer close(errChan)
		defer close(logsChan)

		responseBody, err := client.ContainerLogs(context.Background(), containerID, options)
		errChan <- err
		defer responseBody.Close()

		// reads the logs and push to logs channel
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
func generateJupyterNoteBookStartCommand(baseUrl string) []string {
	cmd := fmt.Sprintf("jupyter notebook --allow-root --ip 0.0.0.0 --NotebookApp.allow_origin=* --NotebookApp.allow_remote_access=True --NotebookApp.base_url='/%s' --NotebookApp.token='' --no-browser --port=8080", baseUrl)
	return strings.Split(cmd, " ")
}

// func generateJupyterNoteBookStartCommand() []string {
// 	cmd := "jupyter notebook --allow-root --ip 0.0.0.0 --NotebookApp.allow_origin=* --NotebookApp.token='' --no-browser --port=8080"
// 	return strings.Split(cmd, " ")
// }

func generateVScodeGoStartCommand() []string {
	cmd := "code-server --auth none /envbox"
	return strings.Split(cmd, " ")
}

func generateJupyterNotebookInternalLink(port, baseUrl string) string {
	return fmt.Sprintf("http://127.0.0.1:%s/%s/tree", port, baseUrl)
}

// func generateJupyterNotebookInternalLink(port string) string {
// 	return fmt.Sprintf("http://127.0.0.1:%s/tree", port)
// }

func generateVsCodeGoInternalLink(port string) string {
	return fmt.Sprintf("http://0.0.0.0:%s/", port)
}

func generateJupyterNotebookAccessLink(port string) string {
	return fmt.Sprintf("http://%s/%s/tree?", config.Envs.Public, port)
}

func generateVsCodeGoAccessLink(port string) string {
	return fmt.Sprintf("http://%s/%s/", config.Envs.Public, port)
}

// func generateJupyterNotebookAccessLink(port string) string {
// 	return fmt.Sprintf("http://127.0.0.1:%s/tree", port)
// }

// func generateVsCodeGoAccessLink(port string) string {
// 	return fmt.Sprintf("http://127.0.0.1:%s/", port)
// }
