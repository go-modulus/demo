package local

import (
	"context"
	"golang.org/x/crypto/bcrypt"
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

func (p Provider) Register(ctx context.Context, account LocalAccount) error {
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
