package service

import (
	"context"
	"demo/internal/errors"
	"demo/internal/user/dao"
	"demo/internal/user/dto"
	"github.com/gofrs/uuid"
	"time"
)

func NewEmailAlreadyInUse() *errors.Error {
	err := errors.NewBusinessLogicError(
		"user.emailAlreadyInUse",
		"This email already in use",
	)

	return err.WithFlags(err.Flags | errors.ErrorDontHandle)
}

type Registration struct {
	finder *dao.UserFinder
	saver  *dao.UserSaver
}

func NewRegistration(finder *dao.UserFinder, saver *dao.UserSaver) *Registration {
	return &Registration{finder: finder, saver: saver}
}

func (r Registration) Register(ctx context.Context, user dto.User) (*dto.User, error) {
	inUse, err := r.isEmailAlreadyInUse(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if inUse {
		return nil, NewEmailAlreadyInUse()
	}

	id, _ := uuid.NewV6()
	user.Id = id.String()
	user.RegisteredAt = time.Now()

	err = r.saver.Create(ctx, user)
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
