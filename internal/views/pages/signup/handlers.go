package signup

import (
	"myapp/internal/auth"
	"myapp/internal/db"
	"myapp/internal/views/layouts"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	// extract error if exists
	errMsg := ""
	if r.URL.Query().Get("error") == "1" {
		errMsg = "Invalid Credentials"
	}

	err := layouts.
		Base("Signup", Signup(errMsg)).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" {
		http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		name, email, hash,
	)
	if err != nil {
		http.Redirect(w, r, "/signup?error=1", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
