package action

import (
	"boilerplate/internal/framework"
	"context"
)

type CurrentUserRequest struct {
}
type CurrentUser struct {
}

func NewCurrentUser() *CurrentUser {
	return &CurrentUser{}
}

func (a *CurrentUser) Handle(ctx context.Context, req *CurrentUserRequest) (framework.CurrentUser, error) {
	cu := framework.GetCurrentUser(ctx)
	if cu == nil {
		return framework.CurrentUser{}, nil
	}
	return *cu, nil
}
