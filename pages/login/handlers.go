package login

import (
	"myapp/auth2"
	"myapp/db"
	"myapp/layouts"
	"myapp/utils"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Login", Login()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	var (
		userID int
		hash   string
	)

	err := db.DB.QueryRow(
		"SELECT id, password_hash FROM users WHERE email = ?",
		email,
	).Scan(&userID, &hash)

	if err != nil {
		return
	}

	if err := auth2.CheckPassword(password, hash); err != nil {
		return
	}

	sess, ok := utils.SessionFromContext(r.Context())
	if !ok {
		return
	}

	sess.Values["user_id"] = userID

	// IMPORTANT: save BEFORE writing response
	if err := sess.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
