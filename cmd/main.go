package main

import (
	"fmt"
	"log"

	"github.com/justsushant/envbox/config"
	"github.com/justsushant/envbox/db"
	"github.com/justsushant/envbox/server"

	"github.com/justsushant/envbox/utils/docker"
)

func main() {
	sqliteDB, err := db.NewSqlLiteStorage(config.Envs.SqliteDB)
	if err != nil {
		log.Fatal("Error while creating sqlite storage")
	}

	// err = db.MigrateSqliteDBUp(sqliteDB)
	// if err != nil {
	// 	log.Fatal("Error while migrating sqlite db: ", err)
	// }


	s := server.NewServer(fmt.Sprintf(":%s", config.Envs.Port), sqliteDB, utils.NewDockerClient())
	s.Run()
}