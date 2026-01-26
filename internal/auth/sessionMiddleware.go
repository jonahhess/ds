package auth

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

func SessionMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "myapp-session")
			if err != nil {
				// handle corrupt or invalid cookies safely
				session, _ = store.New(r, "myapp-session")
			}

			ctx := context.WithValue(r.Context(), SessionContextKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
