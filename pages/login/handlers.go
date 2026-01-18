package login

import (
	"log"
	"myapp/auth"
	"myapp/db"
	"myapp/layouts"
	"myapp/utils"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	// extract error if exists
	errMsg := ""
	if r.URL.Query().Get("error") == "1" {
		errMsg = "Invalid Credentials"
	}

	err := layouts.
		Base("Login", Login(errMsg)).
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

	q := r.URL.Query()

	if err != nil {
		q.Add("error","1")
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		return
	}

	if err := auth.CheckPassword(password, hash); err != nil {
		http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
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
