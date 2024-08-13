package db

import (
	"database/sql"
)

var tableCreate = `
CREATE TABLE IF NOT EXISTS images (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	path TEXT
);

CREATE TABLE IF NOT EXISTS containers_running (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	containerID TEXT,
	imageName TEXT,
	accessLink TEXT,
	active BOOLEAN NOT NULL DEFAULT 1 CHECK (active IN (0, 1)),
	createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func MigrateSqliteDBUp(db *sql.DB) error {
	_, err := db.Exec(tableCreate)
	if err != nil {
		return err
	}
	return nil
}