package auth

import (
	"context"
	"net/http"
)

func OptionalUserMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        sess, _ := SessionFromContext(r.Context())
        if userID, ok := sess.Values[sessionUserIDKey].(int); ok {
            ctx := context.WithValue(r.Context(), userIDKey ,userID)
			r = r.WithContext(ctx)
        }
        next.ServeHTTP(w, r)
    })
}
