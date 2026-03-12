package fixture

import (
	"boilerplate/internal/auth/storage"
	"context"
)

type RevokedTokenFixture struct {
	authDb *storage.Queries
}

func NewRevokedTokenFixture(authDb *storage.Queries) *RevokedTokenFixture {
	return &RevokedTokenFixture{
		authDb: authDb,
	}
}

func (r *RevokedTokenFixture) DeleteByTokenJti(tokenJti string) error {
	return r.authDb.DeleteRevokedToken(context.Background(), tokenJti)
}
