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

func (s *Store) SaveContainer(containerID, imageName string) error {
	_, err := s.db.Exec("INSERT INTO containers_running (containerID, imageName) VALUES (?, ?)", containerID, imageName)
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
	rows, err := s.db.Query("SELECT * FROM containers_running")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []types.Env
	for rows.Next() {
		var env types.Env
		err := rows.Scan(&env.ID, &env.ContainerID, &env.ImageName, &env.Active, &env.CreatedAt)
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
	err := row.Scan(&env.ID, &env.ContainerID, &env.ImageName, &env.Active, &env.CreatedAt)
	if err != nil {
		return types.Env{}, err
	}
	return env, nil
}