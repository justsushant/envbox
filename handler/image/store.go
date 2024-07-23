package image

import (
	"database/sql"
	"github.com/justsushant/envbox/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db: db,
	}
}

func (s *Store) GetImages() ([]types.Image, error) {
	rows, err := s.db.Query("SELECT id, name, path FROM images")
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