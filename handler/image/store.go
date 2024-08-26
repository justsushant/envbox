package image

import (
	"database/sql"
	"github.com/justsushant/envbox/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) types.ImageStore {
	return &Store{
		db: db,
	}
}

func (s *Store) GetImages() ([]types.Image, error) {
	rows, err := s.db.Query("SELECT id, name, path FROM mst_images")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := []types.Image{}
	for rows.Next() {
		var image types.Image
		err := rows.Scan(&image.ID, &image.Name, &image.Path)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (s *Store) GetImageByID(id int) (types.Image, error) {
	var image types.Image
	err := s.db.QueryRow("SELECT id, name, path FROM mst_images WHERE id = ?", id).Scan(&image.ID, &image.Name, &image.Path)
	if err != nil {
		return types.Image{}, err
	}
	return image, nil
}