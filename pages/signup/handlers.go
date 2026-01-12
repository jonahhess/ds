package signup

import (
	"myapp/auth2"
	"myapp/db"
	"myapp/layouts"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Signup", Signup()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SignupUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	hash, err := auth2.HashPassword(password)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, hash,
	)
	if err != nil {
		http.Error(w, "Email already exists", 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
