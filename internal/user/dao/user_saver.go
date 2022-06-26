package dao

import (
	"context"
	"demo/internal/framework"
	"demo/internal/user/dao/query"
	"demo/internal/user/dto"
	"gorm.io/gorm"
)

type UserSaver struct {
	db *gorm.DB
}

func NewUserSaver(db *gorm.DB) *UserSaver {
	return &UserSaver{db: db}
}

func (f *UserSaver) Create(ctx context.Context, user dto.User) error {
	result := f.db.Table(query.UserTable).WithContext(ctx).Create(&user)
	if result.Error != nil {
		return framework.NewGormError(result.Error)
	}

	return nil
}

func (f *UserSaver) Update(ctx context.Context, user dto.User) error {
	result := f.db.Table(query.UserTable).WithContext(ctx).Save(&user)
	if result.Error != nil {
		return framework.NewGormError(result.Error)
	}

	return nil
}
