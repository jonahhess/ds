package myCourse

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

	courseID, ok := params.IntFrom(ctx, "courseID")
	if !ok {
		 return;
	}

	myCourseData, err := GetMyCourseData(DB, userID, courseID)
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		 return;
	}

	csrfToken := auth.CSRFToken(r)
	 if err := layouts.
	 Base("MyCourse", MyCourse(*myCourseData, csrfToken)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

 func GetMyCourseData(DB *sql.DB, userID int, courseID int) (*types.MyCourseData, error) {
    var data types.MyCourseData

    err := DB.QueryRow(`
		SELECT
			uc.user_id,
			uc.course_id,
			c.title,
			c.description,
			u.name,
			u.created_at,
			uc.current_lesson,
			COALESCE(l.title, '') AS current_lesson_name,
			(
				SELECT COUNT(*)
				FROM lessons
				WHERE course_id = uc.course_id
			) AS total_lessons,
			 c.version
		FROM user_courses AS uc
		JOIN courses AS c ON c.id = uc.course_id
		JOIN users AS u ON u.id = c.created_by
		LEFT JOIN lessons l
			ON l.course_id = uc.course_id
			AND l.lesson_index = uc.current_lesson
		WHERE uc.user_id = ?
		AND uc.course_id = ?
    `, userID, courseID).Scan(
        &data.UserID,
        &data.CourseID,
        &data.Title,
        &data.Description,
        &data.CreatedBy,
        &data.CreatedAt,
		&data.CurrentLesson,
		&data.CurrentLessonName,
		&data.TotalLessons,
		&data.Version,
    )

    if err != nil {
        return nil, err
    }

    return &data, nil
}

func Remove(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		courseID, ok := params.IntFrom(ctx, "courseID")
		if !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		_, err := DB.Exec("DELETE FROM user_courses WHERE user_id = ? AND course_id = ?", userID, courseID)
		if err != nil {
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/courses", http.StatusSeeOther)
	}
}