package action

import (
	"boilerplate/internal/framework"
	actionError "boilerplate/internal/user/action/errors"
	"boilerplate/internal/user/dao"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/go-chi/chi/v5"
)

type GetUserRequest struct {
	Id string `in:"path=id;required"`
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

func (a *GetUserAction) Register(chi chi.Router, errorHandler *framework.HttpErrorHandler) error {
	getUser, err := framework.WrapHandler[*GetUserRequest](errorHandler, a)

	if err != nil {
		return err
	}
	chi.Get("/users/{id}", getUser)

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

	return application.NewSuccessResponse(response), nil
}
