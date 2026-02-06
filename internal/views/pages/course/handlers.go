package course

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
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

	courseIDStr := chi.URLParam(r, "courseID")
    if !ok {
        http.Error(w, "course id missing from context", http.StatusBadRequest)
        return
    }

	courseID, err := strconv.Atoi(courseIDStr)
        if err != nil {
            http.Error(w, "invalid course id", http.StatusBadRequest)
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
        if err == sql.ErrNoRows {
            return nil, err // or a custom "not found" error
        }
        return nil, err
    }

    return &data, nil
}

