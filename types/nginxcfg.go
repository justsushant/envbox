package types

type NginxUpstreamConfig struct {
	ContainerID string `json:"containerID"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	IsRewrite bool   `json:"isRewrite"`
}

type NginxCfgWriterService interface {
	AddUpstreamAndRender(containerID, name, address string, isRewrite bool) error
	RemoveUpstreamAndRender(containerID string) error
}
