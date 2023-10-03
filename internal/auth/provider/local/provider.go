package local

import (
	"boilerplate/internal/framework"
	"context"
	"golang.org/x/crypto/bcrypt"
)

var ErrIdentityIsRequired = framework.NewCommonError(
	"ErrIdentityIsRequired",
	"Identity is required. Please provide email or nickname or phone.",
)

type Provider struct {
	storage *GormStorage
}

type ProviderConfig struct {
	AccountTable        string
	SessionIdCookieName string
}

func NewProvider(storage *GormStorage) *Provider {
	return &Provider{storage: storage}
}

// Register the account for a user
// If the account with the same email, nickname or phone already exists, it returns an error
// Errors:
// - ErrIdentityIsRequired - if the account has no email, nickname or phone
func (p Provider) Register(ctx context.Context, account LocalAccount) error {
	if account.Email == nil && account.Nickname == nil && account.Phone == nil {
		return ErrIdentityIsRequired
	}
	data, err := bcrypt.GenerateFromPassword([]byte(account.Password), 14)
	if err != nil {
		return err
	}
	account.Password = string(data)
	return p.storage.Save(ctx, account)
}

func (p Provider) Login(
	ctx context.Context,
	identity string,
	credential string,
) (userId string, err error) {
	localAccount, err := p.storage.LoadByIdentity(ctx, identity)
	userId = ""
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(localAccount.Password),
		[]byte(credential),
	)
	if err != nil {
		return
	}
	userId = localAccount.UserID
	return
}
