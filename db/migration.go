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

`

func MigrateSqliteDBUp(db *sql.DB) error {
	_, err := db.Exec(tableCreate)
	if err != nil {
		return err
	}
	return nil
}