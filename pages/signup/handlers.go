package signup

import (
	"myapp/auth2"
	"myapp/db"
	"myapp/layouts"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Signup", Signup()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	hash, err := auth2.HashPassword(password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		name, email, hash,
	)
	if err != nil {
		http.Error(w, "Email already exists", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
