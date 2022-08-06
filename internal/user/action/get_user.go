package action

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/framework"
	actionError "boilerplate/internal/user/action/errors"
	"boilerplate/internal/user/dao"
	"context"
	application "github.com/debugger84/modulus-application"
)

type GetUserRequest struct {
	Id string `json:"id" in:"path=id;required"`
}

type UserResponse struct {
	Id   string
	Name string
}

type GetUserAction struct {
	finder *dao.UserFinder
}

func NewGetUserAction(finder *dao.UserFinder) *GetUserAction {
	return &GetUserAction{finder: finder}
}

func (a *GetUserAction) Register(
	auth *auth.Auth,
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
) error {
	getUser, err := framework.WrapHandler[*GetUserRequest](errorHandler, a)

	if err != nil {
		return err
	}
	routes.Get("/users/{id}", auth.AuthGuard().Auth(getUser))

	return nil
}

func (a *GetUserAction) Handle(ctx context.Context, request *GetUserRequest) (*application.ActionResponse, error) {
	user, err := a.finder.One(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return actionError.UserNotFound(ctx, request.Id), nil
	}
	var response UserResponse
	response.Id = request.Id
	response.Name = user.Name

	r := application.NewSuccessResponse(response)
	return &r, nil
}
