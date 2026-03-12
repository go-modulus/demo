package service

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/utils"
	"context"
	"strings"
	"time"
)

// token will have symbols count = oneTimePasswordTokenLength * 2
const oneTimePasswordTokenLength = 20

type OtpConfig interface {
	GetOneTimePasswordTtl() time.Duration
	GetOneTimePasswordAfterPurchasingTtl() time.Duration
	GetOneTimePasswordResendTimeout() time.Duration
	GetFrontendHost() string
}

type OtpGenerator interface {
	GenerateOtpCode(ctx context.Context, email string) (string, *time.Time, error)
	GenerateOtpCodeAfterPurchasing(ctx context.Context, email string) (string, *time.Time, error)
}

type OneTimePassword struct {
	config  OtpConfig
	queries *storage.Queries
}

func NewOneTimePassword(
	config OtpConfig,
	queries *storage.Queries,
) *OneTimePassword {
	return &OneTimePassword{
		config:  config,
		queries: queries,
	}
}

func (p *OneTimePassword) GenerateOtpCode(ctx context.Context, email string) (string, *time.Time, error) {
	return p.generateCode(ctx, email, p.config.GetOneTimePasswordTtl())
}

func (p *OneTimePassword) GenerateOtpCodeAfterPurchasing(ctx context.Context, email string) (string, *time.Time, error) {
	return p.generateCode(ctx, email, p.config.GetOneTimePasswordAfterPurchasingTtl())
}

func (p *OneTimePassword) generateCode(ctx context.Context, email string, ttl time.Duration) (string, *time.Time, error) {
	err := p.queries.DeleteOneTimePasswordByEmail(ctx, email)
	if err != nil {
		return "", nil, error2.CannotCreateOneTimePassword
	}

	email = strings.TrimSpace(email)
	token := utils.GenerateSecureToken(oneTimePasswordTokenLength)

	expiresAt := time.Now().Add(ttl)
	canResendAt := time.Now().Add(p.config.GetOneTimePasswordResendTimeout())
	_, err = p.queries.CreateOneTimePassword(ctx, storage.CreateOneTimePasswordParams{
		Token:       token,
		Email:       email,
		ExpiresAt:   expiresAt,
		CanResendAt: canResendAt,
	})
	if err != nil {
		return "", nil, error2.CannotCreateOneTimePassword
	}

	return token, &canResendAt, nil
}
