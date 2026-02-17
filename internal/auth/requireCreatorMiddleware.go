package auth

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// RequireCreatorMiddleware checks that the authenticated user created the specified course.
// The course ID should be available as a URL parameter named "courseID".
func RequireCreatorMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			courseIDStr := chi.URLParam(r, "courseID")
			courseID, err := strconv.Atoi(courseIDStr)
			if err != nil {
				http.Error(w, "Invalid course ID", http.StatusBadRequest)
				return
			}

			// Check if user created this course
			var createdBy int
			err = db.QueryRow(
				"SELECT created_by FROM courses WHERE id = ?",
				courseID,
			).Scan(&createdBy)

			if err == sql.ErrNoRows {
				http.Error(w, "Course not found", http.StatusNotFound)
				return
			}
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			if createdBy != userID {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

