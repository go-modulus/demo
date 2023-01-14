package auth

import (
	"context"
)

const (
	PerformerKey = "performer"
)

func PerformerFromContext(ctx context.Context) NullPerformer {
	if performer, ok := ctx.Value(PerformerKey).(Performer); ok {
		return NullPerformer{
			Value: performer,
			Valid: true,
		}
	}
	return NullPerformer{}
}

func ContextWithPerformer(ctx context.Context, performer Performer) context.Context {
	return context.WithValue(ctx, PerformerKey, performer)
}
