package action

import (
	"boilerplate/internal/auth"
	context2 "boilerplate/internal/auth/context"
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"boilerplate/internal/user/dao"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetUsersRequest struct {
	Count int `in:"query=count"`
}

func (r *GetUsersRequest) Validate(ctx context.Context) *framework.ValidationErrors {
	err := validation.ValidateStructWithContext(
		ctx,
		r,
		validation.Field(
			&r.Count,
			//validation.Required.Error("Count is required"),
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

func InitGetUsersAction(
	auth *auth.Auth,
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	action *GetUsersAction,
) error {
	getUsers, err := framework.WrapHandler[*GetUsersRequest, UsersResponse](errorHandler, action, 200)

	if err != nil {
		return err
	}
	routes.Get("/api/users", getUsers)

	return nil
}

func (a *GetUsersAction) Handle(ctx context.Context, req *GetUsersRequest) (UsersResponse, error) {
	if req.Count == 0 {
		req.Count = 10
	}
	userId := context2.GetCurrentUserId(ctx)
	query := a.finder.CreateQuery(ctx)
	if userId != "" {
		query = query.NotInIds([]string{userId})
	}
	query.NewerFirst()
	users, err := a.finder.ListByQuery(query, req.Count)
	if err != nil {
		return UsersResponse{}, err
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{Id: user.Id, Name: user.Name}
	}

	r := UsersResponse{
		List: response,
	}
	return r, nil
}
