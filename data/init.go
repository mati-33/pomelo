package data

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", path))
	if err != nil {
		return nil, fmt.Errorf("SetubDB: failed to open sqlite3 db: %v", err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, fmt.Errorf("SetupDB: failed to ping sqlite3 db: %v", pingErr)
	}

	return db, nil
}

func SetupTables(conn *sql.DB) error {
	_, err := conn.Exec(`
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS lists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_done INTEGER NOT NULL DEFAULT 0,
    list_id INTEGER NOT NULL,
    FOREIGN KEY (list_id) REFERENCES lists(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

		`)
	if err != nil {
		return fmt.Errorf("SetupTables: failed to create sql tables: %v", err)
	}

	return nil
}
