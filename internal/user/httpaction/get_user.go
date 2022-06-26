package httpaction

import (
	"context"
	"demo/internal/errors"
	"demo/internal/framework"
	pgx2 "demo/internal/pgx"
	"demo/internal/user/dto"
	userErrors "demo/internal/user/errors"
	"demo/internal/user/storage"
	"demo/internal/validator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func (u *GetUserRequest) Validate(ctx context.Context) *validator.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(
			&u.Id,
			dto.IdRules()...,
		),
	)

	if err != nil {
		return validator.FromOzzoError(err)
	}

	return nil
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
			return *framework.NewServerErrorResponse(*userErrors.NewUserNotFound(request.Id))
		}
		return *framework.NewServerErrorResponse(*pgx2.NewPgxError(err))
	}

	var response dto.User
	response.Id = request.Id
	response.Name = user.Name
	return *framework.NewSuccessResponse(response)
}
