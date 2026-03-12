package fixture

import (
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/utils"
	"context"
	"github.com/gofrs/uuid"
	"time"
)

type RefreshTokenFixture struct {
	authDb *storage.Queries
}

func NewRefreshTokenFixture(authDb *storage.Queries) *RefreshTokenFixture {
	return &RefreshTokenFixture{
		authDb: authDb,
	}
}

func (f *RefreshTokenFixture) CreateRandomRefreshToken(refreshToken string, userId uuid.UUID) (storage.RefreshToken, func(), string) {
	sessionId, _ := uuid.NewV6()

	rtHash := f.hashRefreshToken(refreshToken)
	rt, _ := f.authDb.CreateRefreshToken(
		context.Background(), storage.CreateRefreshTokenParams{
			Hash:      rtHash,
			UserID:    userId,
			SessionID: sessionId,
			ExpiresAt: time.Now().Add(time.Hour * 24),
		},
	)
	return rt, func() {
		_ = f.authDb.DeleteRefreshToken(context.Background(), rtHash)
	}, "The refresh token for the user " + userId.String()
}

func (f *RefreshTokenFixture) DeleteRefreshToken(refreshToken string) {
	rtHash := f.hashRefreshToken(refreshToken)
	_ = f.authDb.DeleteRefreshToken(context.Background(), rtHash)
}

func (f *RefreshTokenFixture) SetExpiredRefreshToken(refreshTokenHash string) {
	_, _ = f.authDb.UpdateRefreshTokenExpiresAt(context.Background(), storage.UpdateRefreshTokenExpiresAtParams{
		Hash:      refreshTokenHash,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	})
}

func (f *RefreshTokenFixture) UpdateRefreshTokenUserId(refreshTokenHash string, userId uuid.UUID) {
	_, _ = f.authDb.UpdateRefreshTokenUserId(context.Background(), storage.UpdateRefreshTokenUserIdParams{
		Hash:   refreshTokenHash,
		UserID: userId,
	})
}

func (f *RefreshTokenFixture) hashRefreshToken(rt string) string {
	return utils.HashString(rt)
}
