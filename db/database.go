package db

import (
	"database/sql"
	"myapp/types"

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
func CreateUsersTable() error {
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

func CreateSessionsTable() error {
    query := `
    CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    data BLOB,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
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

func GetUserByID(id int64) (*types.User, error) {
    u := &types.User{}
    query := `SELECT id, name, email, username FROM users WHERE id = ?`
    row := DB.QueryRow(query, id)
    err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Username)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func GetUserByEmail(email string) (*types.User, error) {
    u := &types.User{}
    query := `SELECT id, name, email, username FROM users WHERE email = ?`
    row := DB.QueryRow(query, email)
    err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Username)
    if err != nil {
        return nil, err
    }
    return u, nil
}

