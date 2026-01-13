package signup

import (
	"myapp/auth2"
	"myapp/db"
	"myapp/layouts"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
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

	sse := datastar.NewSSE(w, r)

	if name == "" || email == "" || password == "" {
		sse.PatchElements(`<p id="error">Invalid Credentials</p>`)
		return
	}

	hash, err := auth2.HashPassword(password)
	if err != nil {
		sse.PatchElements(`<p id="error">Invalid Credentials</p>`)
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

	sse.Redirect("/")
}
