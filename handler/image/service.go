package image

import (
	"github.com/justsushant/envbox/types"
)

type Service struct {
	store types.ImageStore
}

func NewService(store types.ImageStore) types.ImageService {
	return &Service{
		store: store,
	}
}

func (s *Service) GetImages() ([]types.Image, error) {
	data, err := s.store.GetImages()
	if err != nil {
		return nil, err
	}
	return data, nil
}