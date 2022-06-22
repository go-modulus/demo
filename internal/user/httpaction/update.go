package httpaction

import (
	"context"
	"demo/internal/framework"
	validator "demo/internal/ozzo-validator"
	"demo/internal/user/dao"
	"demo/internal/user/dto"
	"demo/internal/user/httpaction/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u *UpdateRequest) Validate(ctx context.Context) []framework.ValidationError {
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

	return validator.AsAppValidationErrors(err)
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
		return errors.UserNotFound(ctx, request.Id)
	}
	user.Name = request.Name
	err := a.saver.Update(ctx, *user)
	if err != nil {
		a.logger.Error(ctx, err.Error())
		return errors.CannotUpdateUser(ctx, request.Id)
	}
	return framework.NewSuccessResponse(
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
