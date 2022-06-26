package logger

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, s string, i ...interface{})
	Info(ctx context.Context, s string, i ...interface{})
	Warn(ctx context.Context, s string, i ...interface{})
	Error(ctx context.Context, s string, i ...interface{})
	Panic(ctx context.Context, s string, i ...interface{})
}

func Field(k string, v any) map[string]interface{} {
	return map[string]interface{}{k: v}
}
