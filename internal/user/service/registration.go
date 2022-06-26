package service

import (
	"context"
	"demo/internal/errors"
	"demo/internal/user/dao"
	"demo/internal/user/storage"
	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
)

func NewEmailAlreadyInUse() *errors.Error {
	err := errors.NewBusinessLogicError(
		"user.emailAlreadyInUse",
		"This email already in use",
	)

	return err.WithFlags(err.Flags | errors.ErrorUserFriendly | errors.ErrorDontHandle)
}

type Registration struct {
	finder  *dao.UserFinder
	saver   *dao.UserSaver
	queries *storage.Queries
}

func NewRegistration(
	finder *dao.UserFinder,
	saver *dao.UserSaver,
	queries *storage.Queries,
) *Registration {
	return &Registration{finder: finder, saver: saver, queries: queries}
}

// Register returns emailExists error
func (r Registration) Register(ctx context.Context, request storage.CreateUserParams) (*storage.User, error) {
	inUse, err := r.isEmailAlreadyInUse(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if inUse {
		return nil, NewEmailAlreadyInUse()
	}

	id, _ := uuid.NewV6()
	request.ID = guid.UUID(id)

	user, err := r.queries.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r Registration) isEmailAlreadyInUse(ctx context.Context, email string) (bool, error) {
	query := r.finder.CreateQuery(ctx)
	query.Email(email)
	user, err := r.finder.OneByQuery(query)

	return user != nil, err
}
