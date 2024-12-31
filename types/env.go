package types

import (
	dockerTypes "github.com/docker/docker/api/types"
)

type EnvService interface {
	CreateEnv(CreateEnvPayload) (string, string, string, error)
	KillEnv(string) (string, error)
	GetAllEnvs() ([]GetImageResponse, error)
	GetTerminal(string) (dockerTypes.HijackedResponse, error)
	AddNginxUpstream(int, string, string, string) error
	RemoveNginxUpstream(string) error
}

type EnvStore interface {
	SaveContainer(string, int) error
	DeleteContainer(string) error
	GetAllEnvs() ([]GetImageResponse, error)
	GetContainerByID(string) (Env, error)
	UpdateContainerAccessLink(string, string) error
}

type Env struct {
	ID          string `json:"id"`
	ImageName   string `json:"imageName"`
	ContainerID string `json:"containerID"`
	AccessLink  string `json:"accessLink"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"createdAt"`
}

type CreateEnvPayload struct {
	ImageID int `json:"imageID"`
}

type GetImageResponse struct {
	ID         string `json:"id"`
	ImageName  string `json:"imageName"`
	AccessLink string `json:"accessLink"`
	CreatedAt  string `json:"createdAt"`
}
