package utils

import (
	"context"
	"fmt"
)

func UserLink(ctx context.Context, path string) string {
	userID := GetUserID(ctx)
	return fmt.Sprintf("/%d/%s", userID, path)
}