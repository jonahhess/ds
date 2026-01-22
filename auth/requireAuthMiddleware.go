package auth

import (
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := UserIDFromContext(r.Context())
		if ok {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

