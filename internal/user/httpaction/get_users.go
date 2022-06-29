package httpaction

import (
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"boilerplate/internal/user/dao"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type UserResponse struct {
	Id   string
	Name string
}

type GetUsersRequest struct {
	Count int `json:"count" validate:"required,gte=0,lte=10"`
}

func (u *GetUsersRequest) Validate(ctx context.Context) []framework.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		&u,
		validation.Field(
			&u.Count,
			validation.Required.Error("Count parameter is required"),
			validation.Min(0).Error("Count parameter should be positive"),
			validation.Max(10).Error("Count parameter should be less or equal to 10"),
		),
	)

	return validator.AsAppValidationErrors(err)
}

type UsersResponse struct {
	List []UserResponse `json:"list"`
}

type GetUsersAction struct {
	runner *framework.ActionRunner
	finder *dao.UserFinder
}

func NewGetUsersAction(runner *framework.ActionRunner, finder *dao.UserFinder) *GetUsersAction {
	return &GetUsersAction{runner: runner, finder: finder}
}

func (a *GetUsersAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(
		w, r, func(ctx context.Context, request any) framework.ActionResponse {
			return a.process(ctx, request.(*GetUsersRequest))
		}, &GetUsersRequest{},
	)
}

func (a *GetUsersAction) process(ctx context.Context, request *GetUsersRequest) framework.ActionResponse {
	query := a.finder.CreateQuery(ctx)
	query.NewerFirst()
	users, err := a.finder.ListByQuery(query, request.Count)
	if err != nil {
		return framework.NewServerErrorResponse(ctx, DbError, err)
	}
	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{Id: user.Id, Name: user.Name}
	}

	return framework.NewSuccessResponse(
		UsersResponse{
			List: response,
		},
	)
}
