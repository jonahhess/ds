package signup

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/validation"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(w http.ResponseWriter, r *http.Request) {

	// extract error if exists
	errMsg := ""
	if r.URL.Query().Get("error") == "1" {
		errMsg = "Invalid Credentials"
	}

	csrfToken := auth.CSRFToken(r)
	err := layouts.
		Base("Signup", Signup(errMsg, csrfToken)).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignupHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if !validation.ValidateName(name) {
			http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
			return
		}

		if !validation.ValidateEmail(email) {
			http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
			return
		}

		if !validation.ValidatePassword(password) {
			http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
			return
		}

		hash, err := auth.HashPassword(password)
		if err != nil {
			http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
			return
		}

		_, err = DB.Exec(
			"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
			name, email, hash,
		)
		if err != nil {
			http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
