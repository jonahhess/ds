package handlers

import (
	database "myapp/db"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := database.GetUserByEmail(email)
	if err != nil || !CheckPassword(password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "myapp")
	session.Values["userID"] = user.ID
	session.Options.MaxAge = 86400 // 1 day

	_ = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
