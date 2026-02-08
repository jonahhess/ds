package params

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ctxKey string

func Int(name string) func(http.Handler) http.Handler {
	key := ctxKey(name)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := chi.URLParam(r, name)
			if raw == "" {
				http.Error(w, name+" missing from URL", http.StatusBadRequest)
				return
			}

			val, err := strconv.Atoi(raw)
			if err != nil {
				http.Error(w, "invalid "+name, http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), key, val)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

