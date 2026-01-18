package utils

import (
	"context"

	"github.com/gorilla/sessions"
)

type sessionContextKeyType struct{}

var SessionContextKey = sessionContextKeyType{}

func SessionFromContext(ctx context.Context) (*sessions.Session, bool) {
	sess, ok := ctx.Value(SessionContextKey).(*sessions.Session)
	return sess, ok
}

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

func GetUserID(ctx context.Context) string {
	session, ok := SessionFromContext(ctx)
	if !ok {
		return ""
	}
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return ""
	}
	return userID
}
