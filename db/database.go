package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize opens a SQLite database file
func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	// Optional: set pragmas for performance
	DB.Exec("PRAGMA foreign_keys = ON;")

	return DB.Ping()
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func CreateTables() error {
	query := `
	CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
    `
	_, err := DB.Exec(query)
	return err
}
