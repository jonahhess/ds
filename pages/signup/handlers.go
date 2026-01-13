package signup

import (
	"myapp/auth2"
	"myapp/db"
	"myapp/layouts"
	"net/http"

	"myapp/utils"
)

func Page(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Signup", Signup()).
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

	sess, ok := utils.SessionFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	if name == "" || email == "" || password == "" {
		sess.AddFlash("Invalid email or password")
		sess.Save(r, w)
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	hash, err := auth2.HashPassword(password)
	if err != nil {
		sess.AddFlash("Invalid email or password")
		sess.Save(r, w)
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		name, email, hash,
	)
	if err != nil {
		sess.AddFlash("email already exists")
		sess.Save(r, w)
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
