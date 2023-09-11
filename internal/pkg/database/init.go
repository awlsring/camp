package database

import "database/sql"

func doInit(db *sql.DB) error {
	err := createTables(db)
	if err != nil {
		return err
	}
	return nil
}

func createTables(db *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS machines (
        id UUID NOT NULL,
        PRIMARY KEY (id)
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}
