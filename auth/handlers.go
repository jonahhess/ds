package auth

import (
	"myapp/db"
	"net/http"
)

const sessionUserIDKey = "user_id"


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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, ok := SessionFromContext(r.Context())
	if !ok {
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}