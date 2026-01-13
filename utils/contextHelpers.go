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

func FlashesFromContext(ctx context.Context) []any {
	sess, ok := SessionFromContext(ctx)
	if !ok {
		return nil
	}
	return sess.Flashes()
}

func SetFlashesInContext(ctx context.Context, flashes []any) context.Context {
	sess, ok := SessionFromContext(ctx)
	if !ok {
		return ctx
	}
	for _, flash := range flashes {
		sess.AddFlash(flash)
	}
	return context.WithValue(ctx, SessionContextKey, sess)
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

type flashContextKeyType struct{}

var FlashContextKey = flashContextKeyType{}

func Flashes(ctx context.Context) []any {
	f, _ := ctx.Value(FlashContextKey).([]any)
	return f
}
