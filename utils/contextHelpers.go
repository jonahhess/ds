package utils

import (
	"context"
	"myapp/types"

	"github.com/gorilla/sessions"
)

func GetUserTypeFromContext(ctx context.Context) types.UserType {
	user, ok := ctx.Value("user").(types.User)
	if !ok {
		return types.Guest
	}
	userType := user.Type
	return userType
}

func SessionFromContext(ctx context.Context) (*sessions.Session, bool) {
	session, ok := ctx.Value(types.CtxKey(0)).(*sessions.Session)
	return session, ok
}
