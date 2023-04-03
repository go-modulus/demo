package framework

import (
	"context"
)

type CurrentUser struct {
	Id          string
	Roles       []string
	Permissions []string
	AccessToken string
}

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*CurrentUser, error)
}

type DefaultAuthenticator struct {
}

func NewAuthenticator() Authenticator {
	return &DefaultAuthenticator{}
}

func (a *DefaultAuthenticator) Authenticate(ctx context.Context, token string) (*CurrentUser, error) {
	return nil, nil
}
