package study

import (
	"myapp/auth"
	"myapp/layouts"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Redirect(w,r,"/login",http.StatusUnauthorized)
	}

	if err := layouts.
		Base("Study", Study(userID)).
		Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
