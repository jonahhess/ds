package edit

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/types"

	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courseIDStr := chi.URLParam(r, "courseID")
		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		var title, description sql.NullString
		err = db.QueryRow(
			"SELECT title, description FROM courses WHERE id = ?",
			courseID,
		).Scan(&title, &description)

		if err == sql.ErrNoRows {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		course := types.CourseForm{
			ID:          courseID,
			Title:       title.String,
			Description: description,
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		err = layouts.Base("Edit Course", EditCourse(course, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}