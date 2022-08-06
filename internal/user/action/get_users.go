package action

import (
	"boilerplate/internal/auth"
	context2 "boilerplate/internal/auth/context"
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"boilerplate/internal/user/dao"
	"context"
	application "github.com/debugger84/modulus-application"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetUsersRequest struct {
	Count int `in:"query=count"`
}

func (r *GetUsersRequest) Validate(ctx context.Context) []framework.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		r,
		validation.Field(
			&r.Count,
			validation.Required.Error("Count is required"),
			validation.Min(1).Error("Count should be more than 0."),
			validation.Max(10).Error("Count should be less than or equal 10."),
		),
	)

	return validator.AsAppValidationErrors(err)
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

func (a *GetUsersAction) Register(
	auth *auth.Auth,
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
) error {
	getUsers, err := framework.WrapHandler[*GetUsersRequest](errorHandler, a)

	if err != nil {
		return err
	}
	routes.Get("/users", auth.AuthGuard().Auth(getUsers))

	return nil
}

func (a *GetUsersAction) Handle(ctx context.Context, req *GetUsersRequest) (*application.ActionResponse, error) {
	userId := context2.GetCurrentUserId(ctx)
	query := a.finder.CreateQuery(ctx)
	query.NotInIds([]string{userId}).NewerFirst()
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
