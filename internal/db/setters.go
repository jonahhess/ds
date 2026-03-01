package db

import (
	"database/sql"
)

func CreateCourse(db *sql.DB, userID int, title string, description string) (int, error) {
	result, err := db.Exec(
			"INSERT INTO courses (title, description, created_by) VALUES (?, ?, ?)",
			title, description, userID,
		)
		if err != nil {
			return 0, err
		}

		courseID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(courseID), nil
}