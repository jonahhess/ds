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

// Example: create a users table
func CreateTables() error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        username TEXT
    );
    `
    _, err := DB.Exec(query)
    return err
}

// Example: insert or update a user
func UpsertUser(id, name, email, username string) error {
    query := `
    INSERT INTO users (id, name, email, username)
    VALUES (?, ?, ?, ?)
    ON CONFLICT(id) DO UPDATE SET
        name=excluded.name,
        email=excluded.email,
        username=excluded.username;
    `
    _, err := DB.Exec(query, id, name, email, username)
    return err
}

// Example: fetch a user by ID
type User struct {
    ID       string
    Name     string
    Email    string
    Username string
}

func GetUserByID(id string) (*User, error) {
    u := &User{}
    query := `SELECT id, name, email, username FROM users WHERE id = ?`
    row := DB.QueryRow(query, id)
    err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Username)
    if err != nil {
        return nil, err
    }
    return u, nil
}
