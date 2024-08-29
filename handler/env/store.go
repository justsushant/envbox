package env

import (
	"database/sql"

	"github.com/justsushant/envbox/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) types.EnvStore {
	return &Store{
		db: db,
	}
}

func (s *Store) SaveContainer(containerID string, imageID int) error {
	_, err := s.db.Exec("INSERT INTO containers_running (containerID, imageID) VALUES (?, ?)", containerID, imageID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteContainer(id string) error {
	_, err := s.db.Exec("UPDATE containers_running SET active = 0 WHERE id =  ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAllEnvs() ([]types.Env, error) {
	rows, err := s.db.Query("SELECT * FROM containers_running WHERE active = 1")
	if err != nil {
		if err == sql.ErrNoRows {
			return []types.Env{}, sql.ErrNoRows
		}
		return nil, err
	}
	defer rows.Close()

	var envs []types.Env
	for rows.Next() {
		var env types.Env
		err := rows.Scan(&env.ID, &env.ImageName, &env.ContainerID, &env.AccessLink, &env.Active, &env.CreatedAt)
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}

	return envs, nil
}

func (s *Store) GetContainerByID(id string) (types.Env, error) {
	var env types.Env
	row := s.db.QueryRow("SELECT * FROM containers_running WHERE id = ?", id)
	err := row.Scan(&env.ID, &env.ImageName, &env.ContainerID, &env.AccessLink, &env.Active, &env.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Env{}, sql.ErrNoRows
		}
		return types.Env{}, err
	}
	return env, nil
}

func (s *Store) UpdateContainerAccessLink(containerID, accessLink string) error {
	_, err := s.db.Exec("UPDATE containers_running SET accessLink = ? WHERE containerID = ?", accessLink, containerID)
	if err != nil {
		return err
	}
	return nil
}