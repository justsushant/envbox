package migration

import (
	"database/sql"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/justsushant/envbox/types"
)

func MigrateDbAndImage(cli *client.Client, db *sql.DB, query string, imgOptions []types.ImageMigration) error {
	err := MigrateSqliteDBUp(db, query)
	if err != nil {
		return fmt.Errorf("error while migrating sqlite db: %v", err.Error())
	}

	err = MigrateDockerImgUp(cli, db, imgOptions)
	if err != nil {
		return fmt.Errorf("error while migrating docker image: %v", err.Error())
	}

	return nil
}
