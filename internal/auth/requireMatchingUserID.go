package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ctxKey int

const (
	ctxUserIDFromURL ctxKey = iota
	ctxCourseID
)


func RequireMatchingUserID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authUserID, ok := UserIDFromContext(ctx)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userIDStr := chi.URLParam(r, "userID")
			userIDFromURL, err := strconv.Atoi(userIDStr)
			if err != nil {
				http.Error(w, "invalid user id", http.StatusBadRequest)
				return
			}

			if authUserID != userIDFromURL {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx = context.WithValue(ctx, ctxUserIDFromURL, userIDFromURL)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
