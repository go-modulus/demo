package dao

import (
	gorm2 "boilerplate/internal/gorm"
	userQuery "boilerplate/internal/user/dao/query"
	"boilerplate/internal/user/dto"
	"context"
	"gorm.io/gorm"
)

type UserFinder struct {
	db *gorm.DB
}

func NewUserFinder(db *gorm.DB) *UserFinder {
	return &UserFinder{db: db}
}

func (f *UserFinder) One(ctx context.Context, id string) (*dto.User, error) {
	query := f.CreateQuery(ctx)
	query.Id(id)

	return f.OneByQuery(query)
}

func (f *UserFinder) OneByQuery(query *userQuery.UserQuery) (*dto.User, error) {
	var user *dto.User
	res := query.Build().Limit(1).Scan(&user)

	if res.Error != nil {
		return nil, gorm2.NewGormError(res.Error)
	}

	return user, nil
}

func (f *UserFinder) ListByQuery(query *userQuery.UserQuery, count int) ([]*dto.User, error) {
	var users []*dto.User

	res := query.Build().Limit(count).Scan(&users)

	if res.Error != nil {
		return nil, gorm2.NewGormError(res.Error)
	}

	return users, nil
}

func (f *UserFinder) CreateQuery(ctx context.Context) *userQuery.UserQuery {
	return userQuery.NewUserQuery(ctx, f.db)
}
