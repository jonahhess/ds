package logout

import (
	"myapp/layouts"
	"myapp/types"
	"net/http"

	"github.com/gorilla/sessions"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Logout", Logout()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {

	session := r.Context().Value(types.CtxKey(0)).(*sessions.Session)
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}
