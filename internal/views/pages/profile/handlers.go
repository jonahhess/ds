package profile

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func ViewPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var profile UserProfile
		err := db.QueryRow(
			"SELECT id, email, name, created_at FROM users WHERE id = ?",
			userID,
		).Scan(&profile.ID, &profile.Email, &profile.Name, &profile.CreatedAt)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = layouts.Base("My Profile", ViewProfile(profile)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func EditPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var profile UserProfile
		db.QueryRow(
			"SELECT id, email, name, created_at FROM users WHERE id = ?",
			userID,
		).Scan(&profile.ID, &profile.Email, &profile.Name, &profile.CreatedAt)

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		w.Header().Set("Content-Type", "text/html")
		err := layouts.Base("Edit Profile", EditProfile(profile, csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Update(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")

		if name == "" || email == "" {
			http.Error(w, "Name and email are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(
			"UPDATE users SET name = ?, email = ? WHERE id = ?",
			name, email, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func ChangePasswordPage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		w.Header().Set("Content-Type", "text/html")
		err := layouts.Base("Change Password", ChangePassword(csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ChangePasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		if currentPassword == "" || newPassword == "" || confirmPassword == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if newPassword != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// TODO: Implement password verification with bcrypt
		// For now, placeholder logic

		_, err := db.Exec(
			"UPDATE users SET password = ? WHERE id = ?",
			newPassword, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func DeletePage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		csrfToken := auth.CSRFTokenFromContext(r.Context())
		w.Header().Set("Content-Type", "text/html")
		err := layouts.Base("Delete Account", ConfirmDeleteAccount(csrfToken)).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Delete(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Delete user (cascades remove courses and enrollments due to foreign key constraints)
		_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Clear session
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
