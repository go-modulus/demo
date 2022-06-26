package http

import (
	"context"
	"demo/internal/http"
	"demo/internal/user/dao"
	"demo/internal/validator"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	oHttp "net/http"
)

type GetUsersRequest struct {
	Ctx   context.Context `in:"ctx"`
	First int             `in:"query=first;default=25"`
	After *string         `in:"query=after"`
}

func (u *GetUsersRequest) Validate(ctx context.Context) *validator.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.First,
			validation.Required,
			validation.Min(0),
			validation.Max(25),
		),
	)

	if err != nil {
		return validator.FromOzzoError(err)
	}

	return nil
}

type UsersResponse struct {
	List []UserResponse
}

type GetUsersAction struct {
	finder *dao.UserFinder
}

func NewGetUsersAction(finder *dao.UserFinder) http.HandlerRegistrarResult {
	return http.HandlerRegistrarResult{Handler: &GetUsersAction{finder: finder}}
}

func (a *GetUsersAction) Register(routes *http.Routes) error {
	getUsers, err := http.WrapHandler[*GetUsersRequest](a)

	if err != nil {
		return err
	}
	routes.Get("/users", getUsers)

	return nil
}

func (a *GetUsersAction) Handle(w oHttp.ResponseWriter, req *GetUsersRequest) error {
	query := a.finder.CreateQuery(req.Ctx)
	query.NewerFirst()
	users, err := a.finder.ListByQuery(query, req.First)
	if err != nil {
		return err
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{Id: user.Id, Name: user.Name}
	}

	return json.NewEncoder(w).Encode(UsersResponse{List: response})
}
