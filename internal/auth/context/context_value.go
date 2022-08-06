package context

import (
	"context"
)

type contextKey string

func SetCurrentUserId(ctx context.Context, currentUserId string) context.Context {
	return context.WithValue(ctx, contextKey("CurrentUserId"), currentUserId)
}

func GetCurrentUserId(ctx context.Context) string {
	if value := ctx.Value(contextKey("CurrentUserId")); value != nil {
		return value.(string)
	}
	return ""
}
