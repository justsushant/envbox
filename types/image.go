package types

type ImageService interface {
	GetImages() ([]Image, error)
}

type Image struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type ImageStore interface {
	GetImages() ([]Image, error)
	GetImageByID(int) (Image, error)
}
