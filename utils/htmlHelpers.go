package utils

import (
	"context"
	"fmt"
	"myapp/auth"
)

func UserLink(ctx context.Context, path string) string {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return "/notfound"
	}
	return fmt.Sprintf("/%d/%s", userID, path)
}