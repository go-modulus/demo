package graphql

import (
	"context"

	authGraphql "github.com/go-modulus/auth/graphql"
	emailGraphql "github.com/go-modulus/auth/providers/email/graphql"
	"github.com/go-modulus/demo/internal/auth/graphql"
	"github.com/go-modulus/demo/internal/auth/storage"
	"github.com/go-modulus/demo/internal/graphql/model"
	"github.com/go-modulus/modulus/errors"
	"github.com/go-modulus/modulus/validator"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
)

type Resolver struct {
	emailResolver *emailGraphql.Resolver
	authQueries   *storage.Queries
}

func NewResolver(
	emailResolver *emailGraphql.Resolver,
	authQueries *storage.Queries,
) *Resolver {
	return &Resolver{
		emailResolver: emailResolver,
		authQueries:   authQueries,
	}
}

func (r *Resolver) EmailSignUp(
	ctx context.Context,
	input emailGraphql.EmailSignUpInput,
	userInfo model.UserInfo,
) (authGraphql.TokenPair, error) {
	err := validator.ValidateStructWithContext(
		ctx, &userInfo,
		validation.Field(
			&userInfo.Name,
			validation.Required.Error("Name is required"),
			is.Alpha.Error("Name must contain only letters"),
			validation.Length(2, 20).Error("Name must be between 2 and 20 characters"),
		),
	)
	if err != nil {
		return authGraphql.TokenPair{}, errors.WithTrace(err)
	}
	userID := uuid.Must(uuid.NewV6())
	_, err = r.authQueries.SaveUserInfo(
		ctx, storage.SaveUserInfoParams{
			ID:   userID,
			Name: userInfo.Name,
		},
	)
	if err != nil {
		return authGraphql.TokenPair{}, errors.WithTrace(err)
	}
	input.ID = userID
	input.Roles = []string{graphql.DefaultUserRole}
	tokenPair, err := r.emailResolver.EmailSignUp(ctx, input)
	if err != nil {
		_ = r.authQueries.DeleteUserInfo(ctx, userID)
		return authGraphql.TokenPair{}, errors.WithTrace(err)
	}

	return tokenPair, err
}

func (r *Resolver) EmailSignIn(ctx context.Context, input emailGraphql.EmailSignInInput) (
	authGraphql.TokenPair,
	error,
) {
	return r.emailResolver.EmailSignIn(ctx, input)
}
