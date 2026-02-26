package types

import "database/sql"

type CourseForm struct {
	ID          int
	Title       string
	Version     int
	Description sql.NullString
}