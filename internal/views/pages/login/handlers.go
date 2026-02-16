package login

import (
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/views/layouts"
)

func Page(w http.ResponseWriter, r *http.Request) {

	// extract error if exists
	errMsg := ""
	if r.URL.Query().Get("error") == "1" {
		errMsg = "Invalid Credentials"
	}

	csrfToken := auth.CSRFToken(r)
	err := layouts.
		Base("Login", Login(errMsg, csrfToken)).
		Render(r.Context(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
