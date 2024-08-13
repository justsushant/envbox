package types

import (
	"github.com/docker/docker/client"
)

type EnvService interface {
	CreateEnv(*client.Client, CreateEnvPayload) (string, error)
	KillEnv(*client.Client, string) (string, error)
	GetAllEnvs() ([]Env, error)
}

type EnvStore interface {
	SaveContainer(string, string) error
	DeleteContainer(string) error
	GetAllEnvs() ([]Env, error)
	GetContainerByID(string) (Env, error)
}

type Env struct {
	ID          string `json:"id"`
	ContainerID string `json:"containerID"`
	ImageName   string `json:"imageName"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"createdAt"`
}

type CreateEnvPayload struct {
	// ImageID int `json:"imageId"`
	ImageName string `json:"imageName"`
}
