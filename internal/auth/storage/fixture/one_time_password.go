package fixture

import (
	"boilerplate/internal/auth/service"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/utils"
	"context"
	"fmt"
	"time"
)

type OneTimePasswordFixture struct {
	authDb        *storage.Queries
	sendOtpConfig service.OtpConfig
}

func NewOneTimePasswordFixture(
	authDb *storage.Queries,
	sendOtpConfig service.OtpConfig,
) *OneTimePasswordFixture {
	return &OneTimePasswordFixture{
		authDb:        authDb,
		sendOtpConfig: sendOtpConfig,
	}
}

func (f *OneTimePasswordFixture) DeleteByToken(token string) {
	_ = f.authDb.DeleteOneTimePasswordByToken(context.Background(), token)
}

func (f *OneTimePasswordFixture) GetTokenTtl() time.Duration {
	return f.sendOtpConfig.GetOneTimePasswordTtl()
}

func (f *OneTimePasswordFixture) GetLastToken() *storage.OneTimePassword {
	oneTimePassword, err := f.authDb.GetLastOneTimePassword(context.Background())
	if err != nil {
		return nil
	}

	return &oneTimePassword
}

func (f *OneTimePasswordFixture) GetByToken(token string) *storage.OneTimePassword {
	oneTimePassword, err := f.authDb.GetOneTimePasswordByToken(context.Background(), token)
	if err != nil {
		return nil
	}

	return &oneTimePassword
}

func (f *OneTimePasswordFixture) UpdateExpiresAt(token string, expiresAt time.Time) (storage.OneTimePassword, error) {
	return f.authDb.UpdateOneTimePasswordExpiresAt(context.Background(), storage.UpdateOneTimePasswordExpiresAtParams{
		ExpiresAt: expiresAt,
		Token:     token,
	})
}

func (f *OneTimePasswordFixture) CreateOneTimePassword(email string) (storage.OneTimePassword, func(), string) {
	token := utils.GenerateSecureToken(20)
	expired := time.Now().Add(f.sendOtpConfig.GetOneTimePasswordTtl())
	oneTimePassword, _ := f.authDb.CreateOneTimePassword(context.Background(),
		storage.CreateOneTimePasswordParams{
			Token:     token,
			Email:     email,
			ExpiresAt: expired,
		})

	return oneTimePassword, func() {
			f.DeleteByToken(oneTimePassword.Token)
		},
		fmt.Sprintf("The email '%s' with token '%s' which will be expired at '%s'", email, token, expired.Format(time.DateTime))
}
