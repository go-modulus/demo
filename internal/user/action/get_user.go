package action

import (
	"context"
	"demo/internal/framework"
	"demo/internal/user/dao"
	"demo/internal/user/errors"
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

func (a *GetUserAction) Handle(ctx context.Context, request *GetUserRequest) (*framework.ActionResponse, error) {
	user, err := a.finder.One(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewUserNotFound(request.Id)
	}

	var response UserResponse
	response.Id = request.Id
	response.Name = user.Name
	return framework.NewSuccessResponse(response), nil
}
