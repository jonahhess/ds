package logout

import (
	"myapp/layouts"
	"myapp/utils"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {

	err := layouts.
		Base("Logout", Logout()).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, ok := utils.SessionFromContext(r.Context())
	if !ok {
		http.Error(w, "Failed to retrieve session", http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}
