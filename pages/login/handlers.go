package login

import (
	"myapp/auth2"
	"myapp/db"
	"myapp/layouts"
	"myapp/types"
	"net/http"

	"github.com/gorilla/sessions"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Login", Login()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := auth2.CheckPassword(password, hash); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sess := r.Context().Value(types.CtxKey(0)).(*sessions.Session)
	sess.Values["user_id"] = userID
	sess.Save(r, w)

	w.WriteHeader(http.StatusOK)
}
