package types

import (
	"github.com/docker/docker/client"
	dockerTypes "github.com/docker/docker/api/types"
)

type EnvService interface {
	CreateEnv(*client.Client, CreateEnvPayload) (string, error)
	KillEnv(*client.Client, string) (string, error)
	GetAllEnvs() ([]Env, error)
	GetTerminal(*client.Client, string) (dockerTypes.HijackedResponse, error)
}

type EnvStore interface {
	SaveContainer(string, string) error
	DeleteContainer(string) error
	GetAllEnvs() ([]Env, error)
	GetContainerByID(string) (Env, error)
	UpdateContainerAccessLink(string, string) error
}

type Env struct {
	ID          string `json:"id"`
	ImageName   string `json:"imageName"`
	ContainerID string `json:"containerID"`
	AccessLink string  `json:"accessLink"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"createdAt"`
}

type CreateEnvPayload struct {
	ImageID int `json:"imageID"`
}
