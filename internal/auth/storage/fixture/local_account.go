package fixture

import (
	"boilerplate/internal/auth/storage"
	"context"
	"github.com/gofrs/uuid"
)

type LocalAccountFixture struct {
	db *storage.Queries
}

func NewLocalAccountFixture(userDb *storage.Queries) *LocalAccountFixture {
	return &LocalAccountFixture{
		db: userDb,
	}
}

func (f *LocalAccountFixture) DeleteLocalAccount(userId uuid.UUID) int64 {
	count, _ := f.db.DeleteLocalAccount(context.Background(), userId)
	return count
}

func (f *LocalAccountFixture) DeleteLocalAccountByEmail(email string) int64 {
	count, _ := f.db.DeleteLocalAccountByEmail(context.Background(), email)
	return count
}
