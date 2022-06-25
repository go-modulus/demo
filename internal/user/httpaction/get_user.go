package httpaction

import (
	"boilerplate/internal/user/dto"
	actionError "boilerplate/internal/user/httpaction/errors"
	"boilerplate/internal/user/storage"
	"context"
	"errors"
	application "github.com/debugger84/modulus-application"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const DbError application.ErrorIdentifier = "DbError"

func (u *GetUserRequest) Validate(ctx context.Context) []application.ValidationError {
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

type GetUser struct {
	finder *storage.Queries
}

func NewGetUserProcessor(finder *storage.Queries) GetUserProcessor {
	return &GetUser{finder: finder}
}

func (a *GetUser) Process(ctx context.Context, request *GetUserRequest) application.ActionResponse {
	id, _ := uuid.Parse(request.Id)
	user, err := a.finder.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return actionError.UserNotFound(ctx, request.Id)
		} else {
			return application.NewServerErrorResponse(ctx, DbError, err)
		}
	}
	var response dto.User
	response.Id = request.Id
	response.Name = user.Name
	return application.NewSuccessResponse(response)
}
