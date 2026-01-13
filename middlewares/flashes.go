package middlewares

import (
	"context"
	"myapp/utils"
	"net/http"
)

func FlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := utils.SessionFromContext(r.Context())
		if ok {
			flashes := sess.Flashes()
			ctx := context.WithValue(r.Context(), utils.FlashContextKey, flashes)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}
