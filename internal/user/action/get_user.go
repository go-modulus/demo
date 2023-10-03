package action

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/framework"
	"boilerplate/internal/user/dao"
	"context"
)

var UserNotFound = framework.NewCommonError("UserNotFound", "User with id %s is not found")

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

func InitGetUserAction(
	auth *auth.Auth,
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	action *GetUserAction,
) error {
	getUser, err := framework.WrapHandler[*GetUserRequest, UserResponse](errorHandler, action, 200)

	if err != nil {
		return err
	}
	routes.Get("/api/users/{id}", auth.AuthGuard().Auth(getUser))

	return nil
}

func (a *GetUserAction) Handle(ctx context.Context, request *GetUserRequest) (UserResponse, error) {
	user, err := a.finder.One(ctx, request.Id)
	if err != nil {
		return UserResponse{}, err
	}
	if user == nil {
		return UserResponse{}, UserNotFound.WithTplVariables(request.Id)
	}
	var response UserResponse
	response.Id = request.Id
	response.Name = user.Name

	return response, nil
}
