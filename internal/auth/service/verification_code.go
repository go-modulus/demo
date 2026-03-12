package service

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/errors"
	pgx2 "boilerplate/internal/infra/pgx"
	"boilerplate/internal/infra/utils"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"time"
)

const (
	codeLength = 6
)

type VerificationCodeConfig interface {
	GetVerificationCodeTtl() time.Duration
	GetVerificationCodeForDanaUserTtl() time.Duration
	GetOneTimePasswordResendTimeout() time.Duration
}

type VerificationCode struct {
	config  VerificationCodeConfig
	queries *storage.Queries
}

type VerificationCodePayload struct {
	NftId         string  `json:"nft_id"`
	WalletId      string  `json:"wallet_id"`
	NftTransferId *string `json:"nft_transfer_id"`
	DanaUserId    *string `json:"dana_user_id"`
	ToAddress     *string `json:"to_address"`
}

func NewVerificationCode(
	config VerificationCodeConfig,
	queries *storage.Queries,
) *VerificationCode {
	return &VerificationCode{
		config:  config,
		queries: queries,
	}
}

func (v *VerificationCode) GetLastActiveForTransferNft(ctx context.Context, userId uuid.UUID, nftId, walletId string) (*storage.VerificationCode, error) {

	verificationCode, err := v.queries.GetLastActiveVerificationCodeForTransferNft(ctx, storage.GetLastActiveVerificationCodeForTransferNftParams{
		UserID:   uuid.NullUUID{UUID: userId, Valid: true},
		NftID:    nftId,
		WalletID: walletId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &verificationCode, nil
}

func (v *VerificationCode) GetLastActiveForTransferNftByAddress(ctx context.Context, userId uuid.UUID, nftId, toAddress string) (*storage.VerificationCode, error) {

	verificationCode, err := v.queries.GetLastActiveVerificationCodeForTransferNftByAddress(ctx, storage.GetLastActiveVerificationCodeForTransferNftByAddressParams{
		UserID:    uuid.NullUUID{UUID: userId, Valid: true},
		NftID:     nftId,
		ToAddress: toAddress,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &verificationCode, nil
}

func (v *VerificationCode) GenerateForTransferNft(ctx context.Context, userId uuid.UUID, email, nftId, walletId string, nftTransferId *string, toAddress *string) (*storage.VerificationCode, error) {
	code := utils.RandomNumberString(codeLength)

	payload := VerificationCodePayload{
		NftId:         nftId,
		WalletId:      walletId,
		NftTransferId: nftTransferId,
		ToAddress:     toAddress,
	}
	payloadJson := pgtype.JSONB{}
	err := payloadJson.Set(payload)
	if err != nil {
		return nil, err
	}

	if toAddress != nil {
		err = v.queries.DeleteVerificationCodeForTransferNftByAddress(ctx, storage.DeleteVerificationCodeForTransferNftByAddressParams{
			UserID: uuid.NullUUID{
				UUID:  userId,
				Valid: true,
			},
			NftID:     nftId,
			ToAddress: *toAddress,
		})
	} else {
		err = v.queries.DeleteVerificationCodeForTransferNft(ctx, storage.DeleteVerificationCodeForTransferNftParams{
			UserID: uuid.NullUUID{
				UUID:  userId,
				Valid: true,
			},
			NftID:    nftId,
			WalletID: walletId,
		})
	}

	if err != nil {
		return nil, err
	}

	verificationCode, err := v.queries.CreateVerificationCode(ctx, storage.CreateVerificationCodeParams{
		Code:   code,
		Email:  email,
		Action: storage.VerificationActionTransferNft,
		UserID: uuid.NullUUID{
			UUID:  userId,
			Valid: true,
		},
		Payload:     payloadJson,
		ExpiresAt:   time.Now().Add(v.config.GetVerificationCodeTtl()),
		CanResendAt: time.Now().Add(v.config.GetOneTimePasswordResendTimeout()),
	})

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}

func (v *VerificationCode) GenerateForConfirmDanaUser(ctx context.Context, email, danaUserId string) (*storage.VerificationCode, error) {
	code := utils.RandomNumberString(codeLength)

	payload := VerificationCodePayload{
		DanaUserId: &danaUserId,
	}
	payloadJson := pgtype.JSONB{}
	err := payloadJson.Set(payload)
	if err != nil {
		return nil, err
	}

	err = v.queries.DeleteVerificationCodeForConfirmDanaUser(ctx, danaUserId)
	if err != nil {
		return nil, err
	}

	verificationCode, err := v.queries.CreateVerificationCode(ctx, storage.CreateVerificationCodeParams{
		Code:        code,
		Email:       email,
		Action:      storage.VerificationActionConfirmDanaUser,
		UserID:      uuid.NullUUID{},
		Payload:     payloadJson,
		ExpiresAt:   time.Now().Add(v.config.GetVerificationCodeForDanaUserTtl()),
		CanResendAt: time.Now().Add(v.config.GetOneTimePasswordResendTimeout()),
	})

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}

func (v *VerificationCode) GetByDanaUserId(ctx context.Context, danaUserId string) (*storage.VerificationCode, error) {
	verificationCode, err := v.queries.GetVerificationCodeByDanaUserId(ctx, danaUserId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, error2.VerificationCodeNotFound
		}
		return nil, pgx2.DbIssues
	}

	if verificationCode.UsedAt.Valid {
		return nil, error2.VerificationCodeIsUsed
	}

	if time.Now().After(verificationCode.ExpiresAt) {
		return nil, error2.VerificationCodeIsExpired
	}

	return &verificationCode, err
}

func (v *VerificationCode) GetByCode(ctx context.Context, code string, action storage.VerificationAction) (*storage.VerificationCode, error) {
	verificationCode, err := v.queries.GetVerificationCodeByCodeAndAction(ctx, storage.GetVerificationCodeByCodeAndActionParams{
		Code:   code,
		Action: action,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, error2.VerificationCodeNotFound
		}
		return nil, pgx2.DbIssues
	}

	if verificationCode.UsedAt.Valid {
		return nil, error2.VerificationCodeIsUsed
	}

	if time.Now().After(verificationCode.ExpiresAt) {
		return nil, error2.VerificationCodeIsExpired
	}

	return &verificationCode, err
}

func (v *VerificationCode) MarkUsed(ctx context.Context, userId uuid.UUID, code string) error {
	_, err := v.queries.SetUsedVerificationCode(ctx, storage.SetUsedVerificationCodeParams{
		UserID: userId,
		Code:   code,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return error2.VerificationCodeNotFound
		}
		return pgx2.DbIssues
	}

	return nil
}
