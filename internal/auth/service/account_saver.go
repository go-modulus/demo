package service

import (
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var identityIsNotUnique = framework.NewCommonError("identityIsNotUnique", "Identity is not unique")
var PhoneIsNotUnique = framework.NewCommonError("PhoneIsNotUnique", "Phone is not unique")
var EmailIsNotUnique = framework.NewCommonError("EmailIsNotUnique", "Email is not unique")
var UsernameIsNotUnique = framework.NewCommonError("UsernameIsNotUnique", "Username is not unique")

type Account struct {
	UserId   uuid.UUID
	Phone    *storage.Identity
	Email    *storage.Identity
	Username *storage.Identity
	Password *storage.Password
}

type AccountSaver struct {
	queries *storage.Queries
	db      *pgxpool.Pool
}

func NewAccountSaver(queries *storage.Queries, dbConn *pgxpool.Pool) *AccountSaver {
	return &AccountSaver{queries: queries, db: dbConn}
}

// SaveAccount saves account as set of identities (phone, email, username) and a password
// Errors:
// - UsernameIsNotUnique
// - EmailIsNotUnique
// - PhoneIsNotUnique
func (s *AccountSaver) SaveAccount(
	ctx context.Context,
	userId uuid.UUID,
	phone, email,
	username string,
	password string,
	tx pgx.Tx,
) (*Account, error) {
	result := &Account{UserId: userId}

	qtx := s.queries.WithTx(tx)

	pwd, err := s.savePassword(ctx, qtx, userId, password)
	if err != nil {
		return nil, err
	}
	result.Password = pwd

	phoneIdentity, err := s.saveIdentity(ctx, qtx, userId, phone, storage.IdentityTypePhone)
	if err != nil {
		if err == identityIsNotUnique {
			return nil, PhoneIsNotUnique
		}
		return nil, err
	}
	result.Phone = phoneIdentity

	emailIdentity, err := s.saveIdentity(ctx, qtx, userId, email, storage.IdentityTypeEmail)
	if err != nil {
		if err == identityIsNotUnique {
			return nil, EmailIsNotUnique
		}
		return nil, err
	}
	result.Email = emailIdentity

	usernameIdentity, err := s.saveIdentity(ctx, qtx, userId, username, storage.IdentityTypeUsername)
	if err != nil {
		if err == identityIsNotUnique {
			return nil, UsernameIsNotUnique
		}
		return nil, err
	}
	result.Username = usernameIdentity

	return result, nil
}

// DeleteAccount deletes account
func (s *AccountSaver) DeleteAccount(ctx context.Context, userId uuid.UUID) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)
	if err := qtx.DeleteUserIdentities(ctx, userId); err != nil {
		return err
	}

	if err := qtx.DeleteUserPasswords(ctx, userId); err != nil {
		return err
	}

	if err := qtx.DeleteRefreshTokensByUserId(ctx, userId); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// savePassword Saves password
func (s *AccountSaver) savePassword(
	ctx context.Context,
	qtx *storage.Queries,
	userId uuid.UUID,
	password string,
) (*storage.Password, error) {
	if password == "" {
		return nil, nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	pwd, err := qtx.CreatePassword(
		ctx, storage.CreatePasswordParams{
			ID:           uuid.Must(uuid.NewV6()),
			UserID:       userId,
			PasswordHash: string(hash),
		},
	)
	if err != nil {
		return nil, err
	}
	return &pwd, nil
}

// saveIdentity Saves identity
func (s *AccountSaver) saveIdentity(
	ctx context.Context,
	qtx *storage.Queries,
	userId uuid.UUID,
	identity string,
	identityType storage.IdentityType,
) (*storage.Identity, error) {
	if identity == "" {
		return nil, nil
	}
	_, err := qtx.SelectIdentity(ctx, identity)
	if err == nil {
		return nil, identityIsNotUnique
	}
	identityObj, err := qtx.CreateIdentity(
		ctx, storage.CreateIdentityParams{
			ID:       uuid.Must(uuid.NewV6()),
			UserID:   userId,
			Identity: identity,
			Type:     identityType,
		},
	)
	if err != nil {
		return nil, err
	}
	return &identityObj, nil
}
