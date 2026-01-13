package utils

import (
	"context"

	"github.com/gorilla/sessions"
)

func IsLoggedIn(ctx context.Context) bool {
	userID, ok := ctx.Value("user-id").(int)
	if !ok {
		return false
	}
	return userID != 0
}

func SessionFromContext(ctx context.Context) (*sessions.Session, bool) {
	session, ok := ctx.Value("myapp-session").(*sessions.Session)
	return session, ok
}
