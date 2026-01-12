package middlewares

import (
	"context"
	"net/http"

	"myapp/types"

	"github.com/gorilla/sessions"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sess := r.Context().Value(types.CtxKey(0)).(*sessions.Session)
		userID, ok := sess.Values["user_id"].(int)

		if !ok || userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
