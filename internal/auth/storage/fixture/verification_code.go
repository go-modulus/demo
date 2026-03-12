package fixture

import (
	"boilerplate/internal/auth/service"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/utils"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"strings"
	"time"
)

type VerificationCodeFixture struct {
	queries *storage.Queries
	config  service.VerificationCodeConfig
}

func NewVerificationCodeFixture(
	queries *storage.Queries,
	config service.VerificationCodeConfig,
) *VerificationCodeFixture {
	return &VerificationCodeFixture{
		queries: queries,
		config:  config,
	}
}

func (f *VerificationCodeFixture) DeleteByCode(code string) {
	_ = f.queries.DeleteVerificationCode(context.Background(), code)
}

func (f *VerificationCodeFixture) GetCodeTtl() time.Duration {
	return f.config.GetVerificationCodeTtl()
}

func (f *VerificationCodeFixture) GetCodeTtlForDanaUser() time.Duration {
	return f.config.GetVerificationCodeForDanaUserTtl()
}

func (f *VerificationCodeFixture) GetLastCode(action storage.VerificationAction) *storage.VerificationCode {
	verificationCode, err := f.queries.GetLastVerificationCodeByAction(context.Background(), action)
	if err != nil {
		return nil
	}

	return &verificationCode
}

func (f *VerificationCodeFixture) GetByCode(code string, action storage.VerificationAction) *storage.VerificationCode {
	verificationCode, err := f.queries.GetVerificationCodeByCodeAndAction(context.Background(), storage.GetVerificationCodeByCodeAndActionParams{
		Code:   code,
		Action: action,
	})
	if err != nil {
		return nil
	}

	return &verificationCode
}

func (f *VerificationCodeFixture) UpdateExpiresAt(code string, expiresAt time.Time) (storage.VerificationCode, error) {
	return f.queries.UpdateVerificationCodeExpires(context.Background(), storage.UpdateVerificationCodeExpiresParams{
		ExpiresAt: expiresAt,
		Code:      code,
	})
}

func (f *VerificationCodeFixture) MarkUsed(code string, userId uuid.UUID) (storage.VerificationCode, error) {
	return f.queries.SetUsedVerificationCode(context.Background(), storage.SetUsedVerificationCodeParams{
		UserID: userId,
		Code:   code,
	})
}

func (f *VerificationCodeFixture) CreateVerificationCode(action storage.VerificationAction, userId uuid.UUID, payload *service.VerificationCodePayload) (storage.VerificationCode, func(), string) {
	code := utils.GenerateSecureToken(6)
	expired := time.Now().Add(f.config.GetVerificationCodeTtl())
	email := "test_" + utils.RandomString(12) + "@test.com"
	email = strings.ToLower(email)
	payloadJson := pgtype.JSONB{}
	_ = payloadJson.Set(payload)
	verificationCode, _ := f.queries.CreateVerificationCode(context.Background(),
		storage.CreateVerificationCodeParams{
			Code:   code,
			Action: action,
			Email:  email,
			UserID: uuid.NullUUID{
				UUID:  userId,
				Valid: true,
			},
			Payload:   payloadJson,
			ExpiresAt: expired,
		})

	return verificationCode, func() {
			f.DeleteByCode(verificationCode.Code)
		},
		fmt.Sprintf("The verification code '%s' for action '%s' which will be expired at '%s'", code, action, expired.Format(time.DateTime))
}
