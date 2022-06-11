package service

import (
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/dto"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/gofrs/uuid"
	"time"
)

const emailExists application.ErrorIdentifier = "emailExists"

type Registration struct {
	finder *dao.UserFinder
	saver  *dao.UserSaver
}

func NewRegistration(finder *dao.UserFinder, saver *dao.UserSaver) *Registration {
	return &Registration{finder: finder, saver: saver}
}

// Register returns emailExists error
func (r Registration) Register(ctx context.Context, user dto.User) (*dto.User, error) {
	if r.emailExist(ctx, user.Email) {
		return nil, application.NewCommonError(emailExists, "not unique email")
	}
	id, _ := uuid.NewV6()
	user.Id = id.String()
	user.RegisteredAt = time.Now()

	err := r.saver.Create(ctx, user)
	if err != nil {
		//r.logger.Error(ctx, err.Error())
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
