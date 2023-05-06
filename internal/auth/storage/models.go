// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package storage

import (
	"database/sql"

	uuid "github.com/gofrs/uuid"
	null "gopkg.in/guregu/null.v4"
)

type Account struct {
	ID                 uuid.UUID    `db:"id" json:"id"`
	Email              null.String  `db:"email" json:"email"`
	Password           null.String  `db:"password" json:"password"`
	ConfirmSelector    null.String  `db:"confirm_selector" json:"confirmSelector"`
	ConfirmVerifier    null.String  `db:"confirm_verifier" json:"confirmVerifier"`
	Confirmed          bool         `db:"confirmed" json:"confirmed"`
	AttemptCount       int32        `db:"attempt_count" json:"attemptCount"`
	LastAttemptAt      sql.NullTime `db:"last_attempt_at" json:"lastAttemptAt"`
	LockedAt           sql.NullTime `db:"locked_at" json:"lockedAt"`
	RecoverSelector    null.String  `db:"recover_selector" json:"recoverSelector"`
	RecoverVerifier    null.String  `db:"recover_verifier" json:"recoverVerifier"`
	RecoverTokenExpiry sql.NullTime `db:"recover_token_expiry" json:"recoverTokenExpiry"`
	Oauth2Uid          null.String  `db:"oauth2_uid" json:"oauth2Uid"`
	Oauth2Provider     null.String  `db:"oauth2_provider" json:"oauth2Provider"`
	Oauth2AccessToken  null.String  `db:"oauth2_access_token" json:"oauth2AccessToken"`
	Oauth2RefreshToken null.String  `db:"oauth2_refresh_token" json:"oauth2RefreshToken"`
	Oauth2Expiry       sql.NullTime `db:"oauth2_expiry" json:"oauth2Expiry"`
	TotpSecretKey      null.String  `db:"totp_secret_key" json:"totpSecretKey"`
	SmsPhoneNumber     null.String  `db:"sms_phone_number" json:"smsPhoneNumber"`
	SmsSeedPhoneNumber null.String  `db:"sms_seed_phone_number" json:"smsSeedPhoneNumber"`
	RecoveryCodes      null.String  `db:"recovery_codes" json:"recoveryCodes"`
}

type LocalAccount struct {
	UserID    uuid.UUID    `db:"user_id" json:"userID"`
	Email     null.String  `db:"email" json:"email"`
	Nickname  null.String  `db:"nickname" json:"nickname"`
	Phone     null.String  `db:"phone" json:"phone"`
	Password  null.String  `db:"password" json:"password"`
	CreatedAt sql.NullTime `db:"created_at" json:"createdAt"`
}

type RememberToken struct {
	ID        uuid.UUID `db:"id" json:"id"`
	AccountID uuid.UUID `db:"account_id" json:"accountID"`
	Token     string    `db:"token" json:"token"`
}