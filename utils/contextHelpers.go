package utils

import (
	"context"

	"github.com/gorilla/sessions"
)

func IsLoggedIn(ctx context.Context) bool {
	session, ok := SessionFromContext(ctx)
	if !ok {
		return false
	}
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return false
	}
	return userID != 0
}

func SessionFromContext(ctx context.Context) (*sessions.Session, bool) {
	session, ok := ctx.Value("session-key").(*sessions.Session)
	return session, ok
}
