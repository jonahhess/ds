package course

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/types"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return
	}

	courseID, ok := params.IntFrom(r.Context(), "courseID")
		if !ok {
			http.Error(w, "course id not found", http.StatusInternalServerError)
			return
		}

	myData, err := GetCourseData(DB, userID, courseID)
	if err != nil {
		 http.Error(w, "invalid course id", http.StatusInternalServerError)
		 return
	}

	 if err := layouts.
	 Base("Course", Course(userID, courseID, *myData)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

 func GetCourseData(DB *sql.DB, userID int, courseID int) (*types.CourseData, error) {
    var data types.CourseData

    err := DB.QueryRow(`
		SELECT
			c.title,
			c.description,
			u.name,
			u.created_at,
			COALESCE((
            SELECT GROUP_CONCAT(l.title, '; ')
            FROM lessons l
            WHERE l.course_id = c.id
            ORDER BY l.lesson_index ASC
        ), '') AS lesson_titles,
			EXISTS (
				SELECT 1
				FROM user_courses
				WHERE user_id = ?
				AND course_id = ?
			) AS user_currently_enrolled
		FROM courses AS c
		JOIN users AS u ON u.id = c.created_by
		WHERE c.id = ?
		`, userID, courseID, courseID).Scan(
		&data.Title,
		&data.Description,
		&data.CreatedBy,
		&data.CreatedAt,
		&data.Lessons,
		&data.UserCurrentlyEnrolled,
	)

    if err != nil {
        return nil, err
    }

    return &data, nil
}

func Enroll(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
				http.Error(w, "invalid user id", http.StatusInternalServerError)
				return
			}
			
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "course id not found", http.StatusInternalServerError)
			return
		}

		if _, err := DB.Exec("INSERT INTO user_courses (user_id, course_id, current_lesson) VALUES (?, ?, 0)", userID, courseID); err != nil {
			http.Error(w, "Enroll error: ", http.StatusConflict)
			return
		}

		myData, err := GetCourseData(DB, userID, courseID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := layouts.
		Base("Course", Course(userID, courseID, *myData)).
		Render(r.Context(), w);  err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}