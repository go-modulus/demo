package auth

import (
	"context"
	"github.com/gofrs/uuid"
)

type Auth0Provider struct {
}

func NewAuth0Provider() *Auth0Provider {
	return &Auth0Provider{}
}

func (Auth0Provider) GetUser(ctx context.Context, token string) (NullPerformer, error) {
	return NullPerformer{
		Value: Performer{Id: uuid.Must(uuid.FromString(token))},
		Valid: true,
	}, nil
}
