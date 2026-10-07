package requestctx

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

var ErrUserIDNotFound = errors.New("user ID not found in request context")

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(userIDKey).(uint)
	if !ok || userID == 0 {
		return 0, ErrUserIDNotFound
	}

	return userID, nil
}
