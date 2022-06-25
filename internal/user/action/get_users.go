package action

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/dao"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/go-chi/chi/v5"
)

type GetUsersRequest struct {
	Count int `in:"query=count"`
}

type UsersResponse struct {
	List []UserResponse `json:"list"`
}

type GetUsersAction struct {
	finder *dao.UserFinder
}

func NewGetUsersAction(finder *dao.UserFinder) *GetUsersAction {
	return &GetUsersAction{finder: finder}
}

func (a *GetUsersAction) Register(chi chi.Router, errorHandler *framework.HttpErrorHandler) error {
	getUsers, err := framework.WrapHandler[*GetUsersRequest](errorHandler, a)

	if err != nil {
		return err
	}
	chi.Get("/users", getUsers)

	return nil
}

func (a *GetUsersAction) Handle(ctx context.Context, req *GetUsersRequest) (*application.ActionResponse, error) {
	query := a.finder.CreateQuery(ctx)
	query.NewerFirst()
	users, err := a.finder.ListByQuery(query, req.Count)
	if err != nil {
		return nil, err
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{Id: user.Id, Name: user.Name}
	}

	r := application.NewSuccessResponse(
		UsersResponse{
			List: response,
		},
	)
	return &r, nil
}
