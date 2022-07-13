package httpaction

import (
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"boilerplate/internal/user/dto"
	actionError "boilerplate/internal/user/httpaction/errors"
	"boilerplate/internal/user/storage"
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const DbError framework.ErrorIdentifier = "DbError"

func (u *GetUserRequest) Validate(ctx context.Context) []framework.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Id,
			dto.IdRules()...,
		),
	)

	return validator.AsAppValidationErrors(err)
}

type TestOverride interface {
	Name() string
}

type Override struct {
}

func (o Override) Name() string {
	return "Base"
}

type GetUser struct {
	finder *storage.Queries
	name   TestOverride
}

func NewGetUserProcessor(finder *storage.Queries, name TestOverride) GetUserProcessor {
	return &GetUser{finder: finder, name: name}
}

func (a *GetUser) Process(ctx context.Context, request *GetUserRequest) framework.ActionResponse {
	id, _ := uuid.Parse(request.Id)
	user, err := a.finder.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return actionError.UserNotFound(ctx, request.Id)
		} else {
			return framework.NewServerErrorResponse(ctx, DbError, err)
		}
	}
	var response dto.User
	response.Id = a.name.Name()
	response.Name = user.Name
	return framework.NewSuccessResponse(response)
}
