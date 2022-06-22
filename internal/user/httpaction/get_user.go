package httpaction

import (
	"context"
	"demo/internal/framework"
	validator "demo/internal/ozzo-validator"
	"demo/internal/user/dto"
	actionError "demo/internal/user/httpaction/errors"
	"demo/internal/user/storage"
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

type GetUser struct {
	finder *storage.Queries
}

func NewGetUserProcessor(finder *storage.Queries) GetUserProcessor {
	return &GetUser{finder: finder}
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
	response.Id = request.Id
	response.Name = user.Name
	return framework.NewSuccessResponse(response)
}
