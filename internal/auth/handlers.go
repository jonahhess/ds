package auth

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/validation"
)

const sessionUserIDKey = "user_id"


func LoginHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate email format
		if !validation.ValidateEmail(email) {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}

		var (
			userID int
			hash   string
		)

		err := DB.QueryRow(
			"SELECT id, password_hash FROM users WHERE email = ?",
			email,
		).Scan(&userID, &hash)

		// Always perform password check to prevent timing attacks
		// If user doesn't exist, use a dummy hash to maintain constant time
		if err != nil {
			// Use a valid bcrypt hash to perform comparison (takes similar time)
			// This is a hash of "dummy" to ensure constant-time comparison
			dummyHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
			_ = CheckPassword(password, dummyHash)
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}

		if err := CheckPassword(password, hash); err != nil {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}

		sess, ok := SessionFromContext(r.Context())
		if !ok {
			http.Error(w, "Failed to get session", http.StatusInternalServerError)
			return
		}
		sess.Values[sessionUserIDKey] = userID

		if err := sess.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, ok := SessionFromContext(r.Context())
	if !ok {
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}