package types

import "time"

type MyCourseData struct {
    UserID   int // user_courses
	CourseID int
    Title string // join courses on c.id === course_id
	Description string
	CreatedBy string // join users on u.id === created_by
	CreatedAt time.Time
	CurrentLesson int
	TotalLessons  int       // count of lessons
	CurrentLessonName string
	Version int
}