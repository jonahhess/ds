package course

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/errors"
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
			errors.HandleBadRequest(w, r, "course id not found")
			return
		}

	myData, err := GetCourseData(DB, userID, courseID)
	if err != nil {
		 errors.HandleNotFound(w, r, "Course")
		 return
	}

	csrfToken := auth.CSRFToken(r)
	 if err := layouts.
	 Base("Course", Course(userID, courseID, *myData, csrfToken)).
	 Render(r.Context(), w);  err != nil {
		 errors.HandleInternalError(w, r, err)
		 return
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
		WHERE c.id = ? AND c.version > 0
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
				errors.HandleUnauthorized(w, r)
				return
			}
			
		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			errors.HandleBadRequest(w, r, "course id not found")
			return
		}

		if _, err := DB.Exec("INSERT INTO user_courses (user_id, course_id, current_lesson) VALUES (?, ?, 0)", userID, courseID); err != nil {
			errors.HandleError(w, r, err, http.StatusConflict, "Already enrolled in this course")
			return
		}

		myData, err := GetCourseData(DB, userID, courseID)
		if err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}

		csrfToken := auth.CSRFToken(r)
		if err := layouts.
		Base("Course", Course(userID, courseID, *myData, csrfToken)).
		Render(r.Context(), w);  err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}
	}
}