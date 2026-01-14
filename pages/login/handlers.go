package login

import (
	"log"
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
		return
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
		http.Error(w, "nvalid Credentials", http.StatusNoContent)
		return
	}

	if err := auth2.CheckPassword(password, hash); err != nil {
		http.Error(w, "Invalid Credentials", http.StatusNoContent)
		return
	}

	sess, ok := utils.SessionFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	sess.Values["user_id"] = userID

	if err := sess.Save(r, w); err != nil {
		log.Printf("session save error: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
