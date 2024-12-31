package main

import (
	"fmt"

	"github.com/justsushant/envbox/config"
	"github.com/justsushant/envbox/db"
	"github.com/justsushant/envbox/migration"
	"github.com/justsushant/envbox/server"
	"github.com/justsushant/envbox/utils"
)

func main() {
	sqliteDB, err := db.NewSqlLiteStorage(config.Envs.SqliteDB)
	if err != nil {
		panic(err)
	}

	dockerClient, err := utils.NewDockerClient()
	if err != nil {
		panic(err)
	}

	// run migrate and exit on the basis of migrate env variable
	if config.Envs.Migrate == "true" {
		err := migration.MigrateDbAndImage(dockerClient, sqliteDB, utils.CREATE_TABLES_STMT, utils.IMAGE_MIG_TABLE)
		if err != nil {
			panic(err)
		}
		fmt.Println("Migration completed successfully")
		return
	}

	s := server.NewServer(fmt.Sprintf(":%s", config.Envs.Port), sqliteDB, dockerClient)
	s.Run()
}
