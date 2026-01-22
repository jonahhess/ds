package login

import (
	"myapp/layouts"
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
