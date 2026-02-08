package myCourse

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

	userIDStrFromURL := chi.URLParam(r, "userID")
	if !ok {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	
	userIDFromURL, err := strconv.Atoi(userIDStrFromURL)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if userID != userIDFromURL {
		http.Error(w, "invalid user id", http.StatusBadRequest)
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

	myCourseData, err := GetMyCourseData(DB, userID, courseID)
	if err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		 return;
	}

	 if err := layouts.
	 Base("MyCourse", MyCourse(userID, *myCourseData)).
	 Render(r.Context(), w);  err != nil {
		 http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

 func GetMyCourseData(DB *sql.DB, userID int, courseID int) (*types.MyCourseData, error) {
    var data types.MyCourseData

    err := DB.QueryRow(`
		WITH ordered_lessons AS (
			SELECT
				l.course_id,
				l.title,
				ROW_NUMBER() OVER (PARTITION BY l.course_id ORDER BY l.id ASC) - 1 AS lesson_index
			FROM lessons l
		)
		SELECT
			uc.user_id,
			uc.course_id,
			c.title,
			c.description,
			u.name,
			u.created_at,
			uc.current_lesson,
			COALESCE(ol.title, '') AS current_lesson_name,
			(
				SELECT COUNT(*)
				FROM lessons l
				WHERE l.course_id = uc.course_id
			) AS total_lessons
		FROM user_courses AS uc
		JOIN courses AS c ON c.id = uc.course_id
		JOIN users AS u ON u.id = c.created_by
		LEFT JOIN ordered_lessons ol
			ON ol.course_id = uc.course_id
			AND ol.lesson_index = uc.current_lesson
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
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err // or a custom "not found" error
        }
        return nil, err
    }

    return &data, nil
}

