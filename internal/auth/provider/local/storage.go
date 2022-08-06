package local

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GormStorage struct {
	db  *gorm.DB
	cfg *ProviderConfig
}

func NewGormStorage(db *gorm.DB, cfg *ProviderConfig) *GormStorage {
	return &GormStorage{db: db, cfg: cfg}
}

// Save the user
func (s GormStorage) Save(ctx context.Context, entity LocalAccount) error {
	result := s.db.WithContext(ctx).Table(s.cfg.AccountTable).Create(&entity)

	return result.Error
}

func (s GormStorage) LoadByUserId(ctx context.Context, uid string) (*LocalAccount, error) {
	var localAccount LocalAccount
	result := s.db.
		WithContext(ctx).
		Table(s.cfg.AccountTable).
		Where("user_id = ?", uid).
		Limit(1).
		Find(&localAccount)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &localAccount, nil
}

func (s GormStorage) LoadByIdentity(ctx context.Context, identity string) (*LocalAccount, error) {
	var localAccount LocalAccount
	result := s.db.
		WithContext(ctx).
		Table(s.cfg.AccountTable).
		Where("email = ? or phone = ? or nickname = ?", identity, identity, identity).
		Limit(1).
		Find(&localAccount)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &localAccount, nil
}
