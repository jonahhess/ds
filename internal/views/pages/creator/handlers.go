package creator

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch courses created by this user
		rows, err := db.Query("SELECT id, title, description, created_at FROM courses WHERE created_by = ? ORDER BY created_at DESC", userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var courses []CreatorCourse
		for rows.Next() {
			var c CreatorCourse
			err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.CreatedAt)
			if err != nil {
				continue
			}
			courses = append(courses, c)
		}

		err = layouts.Base("Creator Dashboard", CreatorDashboard(courses)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
