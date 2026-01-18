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
	users := `
	CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
    `

	DB.Exec(users)

	courses := `
	CREATE TABLE courses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  description TEXT,
  created_by TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
    `
 DB.Exec(courses)

 	lessons := `
	CREATE TABLE lessons (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  course_id INTEGER FOREIGN KEY,
  title TEXT NOT NULL,
  text TEXT NOT NULL,
  quiz_id INTEGER FOREIGN KEY, 
  created_by TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
    `
 DB.Exec(lessons);

  questions := `
	CREATE TABLE questions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  lesson_id INTEGER FOREIGN KEY,
  title TEXT NOT NULL,
  text TEXT NOT NULL,
  created_by TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
    `
 DB.Exec(questions);

  answers := `
	CREATE TABLE answers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id INTEGER FOREIGN KEY,
  text TEXT NOT NULL
);
    `
 DB.Exec(answers);

  correctAnswer := `
	CREATE TABLE correctAnswer (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id INTEGER FOREIGN KEY,
  answer_id INTEGER FOREIGN KEY
);
    `
 DB.Exec(correctAnswer);

	return nil
}
