package middlewares

import (
	"log"
	"net/http"

	"myapp/utils"
)

func SessionSaver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		sess, ok := utils.SessionFromContext(r.Context())
		if ok {
			if err := sess.Save(r, w); err != nil {
				log.Printf("session save error: %v", err)
			}
		}
	})
}
