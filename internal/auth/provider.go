package auth

import "context"

type Provider interface {
	GetUser(ctx context.Context, token string) (NullPerformer, error)
}
