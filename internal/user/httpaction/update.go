package httpaction

import (
	"context"
	"demo/internal/errors"
	"demo/internal/framework"
	"demo/internal/user/dao"
	"demo/internal/user/dto"
	userErrors "demo/internal/user/errors"
	"demo/internal/validator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u *UpdateRequest) Validate(ctx context.Context) *validator.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Id,
			dto.IdRules()...,
		),
		validation.Field(
			&u.Name,
			dto.NameRules()...,
		),
	)

	if err != nil {
		return validator.FromOzzoError(err)
	}

	return nil
}

type Update struct {
	finder *dao.UserFinder
	saver  *dao.UserSaver
	logger framework.Logger
}

func NewUpdateProcessor(
	finder *dao.UserFinder,
	saver *dao.UserSaver,
	logger framework.Logger,
) UpdateProcessor {
	return &Update{finder: finder, saver: saver, logger: logger}
}

func (a *Update) Process(ctx context.Context, request *UpdateRequest) framework.ActionResponse {
	user := a.getUser(ctx, request.Id)
	if user == nil {
		return *framework.NewServerErrorResponse(*userErrors.NewUserNotFound(request.Id))
	}
	user.Name = request.Name
	err := a.saver.Update(ctx, *user)
	if err != nil {
		return *framework.NewServerErrorResponse(*errors.FromError(err))
	}

	return *framework.NewSuccessResponse(
		dto.User{
			Id:   request.Id,
			Name: user.Name,
		},
	)
}

func (a *Update) getUser(ctx context.Context, id string) *dto.User {
	query := a.finder.CreateQuery(ctx)
	query.Id(id)
	user, _ := a.finder.OneByQuery(query)

	return user
}
