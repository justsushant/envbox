package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqlLiteStorage(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}