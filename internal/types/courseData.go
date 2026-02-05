package types

import "time"

type  CourseData struct {
	CourseID int
    Title string
	Description string
	CreatedBy string
	CreatedAt time.Time
	Lessons string
	UserCurrentlyEnrolled bool
}