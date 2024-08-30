package migration

import (
	"database/sql"
)

func MigrateSqliteDBUp(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func insertImageRecordInDB(db *sql.DB, name, path string) error {
	rows, err := db.Query("SELECT * FROM mst_images WHERE name = ? AND path = ?", name, path)
	if err != nil {
		return err
	}
	if rows.Next() {
		return nil
	}

	_, err = db.Exec("INSERT INTO mst_images (name, path) VALUES (?, ?)", name, path)
	if err != nil {
		return err
	}
	return nil
}