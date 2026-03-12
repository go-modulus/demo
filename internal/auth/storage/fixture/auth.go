package fixture

import (
	"boilerplate/internal/auth/storage"
	userStorage "boilerplate/internal/user/storage"
	"boilerplate/internal/user/storage/fixture"
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthFixture struct {
	authDb      *storage.Queries
	userFixture *fixture.UserFixture
}

func NewAuthFixture(authDb *storage.Queries, userFixture *fixture.UserFixture) *AuthFixture {
	return &AuthFixture{authDb: authDb, userFixture: userFixture}
}

func (f *AuthFixture) CreateDirectLoginCreds(login string, password string, identityType storage.IdentityType) (
	storage.Identity,
	storage.Password,
	func(),
	userStorage.User,
	string,
) {
	userId, _ := uuid.NewV6()
	id, _ := uuid.NewV6()

	user, rb, _ := f.userFixture.CreateParticularUser(userId, login)
	if identityType == storage.IdentityTypeEmail {
		login = user.VerifiedEmail.String
	}

	identity, _ := f.authDb.CreateIdentity(
		context.Background(), storage.CreateIdentityParams{
			ID:       id,
			UserID:   userId,
			Identity: login,
			Type:     identityType,
		},
	)

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pwd, _ := f.authDb.CreatePassword(
		context.Background(), storage.CreatePasswordParams{
			ID:           id,
			UserID:       userId,
			PasswordHash: string(hash),
		},
	)
	return identity, pwd, func() {
			rb()
			_ = f.authDb.DeleteIdentity(context.Background(), id)
			_ = f.authDb.DeletePassword(context.Background(), id)
		},
		user,
		fmt.Sprintf("The identity '%s' with the password '%s' for the user '%s'", login, password, userId.String())
}

func (f *AuthFixture) CreateIdentity(
	userId uuid.UUID,
	email string,
) (
	storage.Identity,
	func(),
	string,
) {
	id, _ := uuid.NewV6()

	identity, _ := f.authDb.CreateIdentity(
		context.Background(), storage.CreateIdentityParams{
			ID:       id,
			UserID:   userId,
			Identity: email,
			Type:     storage.IdentityTypeEmail,
		},
	)

	return identity, func() {
		_ = f.authDb.DeleteIdentity(context.Background(), id)
	}, fmt.Sprintf("The identity '%s' for the user '%s'", email, userId.String())
}

func (f *AuthFixture) CreatePasswordForUser(userId uuid.UUID, password string) (
	storage.Password,
	func(),
	string,
) {
	id, _ := uuid.NewV6()

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pwd, _ := f.authDb.CreatePassword(
		context.Background(), storage.CreatePasswordParams{
			ID:           id,
			UserID:       userId,
			PasswordHash: string(hash),
		},
	)
	return pwd, func() {
		_ = f.authDb.DeletePassword(context.Background(), id)
	}, fmt.Sprintf("The new password '%s' for the user '%s'", password, userId.String())
}

func (f *AuthFixture) GetIdentity(identity string) *storage.Identity {
	identityObject, err := f.authDb.SelectIdentity(context.Background(), identity)
	if err == nil {
		return &identityObject
	}

	return nil
}
