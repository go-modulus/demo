package service

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/storage"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
)

const emailExists application.ErrorIdentifier = "emailExists"

type Registration struct {
	finder  *dao.UserFinder
	saver   *dao.UserSaver
	queries *storage.Queries
	logger  framework.Logger
}

func NewRegistration(
	finder *dao.UserFinder,
	saver *dao.UserSaver,
	queries *storage.Queries,
	logger framework.Logger,
) *Registration {
	return &Registration{finder: finder, saver: saver, queries: queries, logger: logger}
}

// Register returns emailExists error
func (r Registration) Register(ctx context.Context, request storage.CreateUserParams) (*storage.User, error) {
	if r.emailExist(ctx, request.Email) {
		return nil, application.NewCommonError(emailExists, "not unique email")
	}
	id, _ := uuid.NewV6()
	request.ID = guid.UUID(id)

	user, err := r.queries.CreateUser(ctx, request)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	return &user, nil
}

func (r Registration) emailExist(ctx context.Context, email string) bool {
	query := r.finder.CreateQuery(ctx)
	query.Email(email)
	user, _ := r.finder.OneByQuery(query)

	return user != nil
}
