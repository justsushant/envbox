package types

type ImageService interface {
	GetImages() ([]Image, error)
}

type Image struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}