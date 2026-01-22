package auth

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

func UserIDFromContext(ctx context.Context) (int, bool) {
	uid, ok := ctx.Value(userIDKey).(int)
	return uid, ok
}

