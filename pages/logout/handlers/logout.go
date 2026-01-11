package handlers

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "myapp")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
