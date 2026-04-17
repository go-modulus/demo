package dataloader

import (
	"github.com/go-modulus/demo/internal/auth/storage"
)

type LoaderFactory struct {
	db             storage.DBTX
	userInfoLoader *UserInfoLoader
}

func NewLoaderFactory(db storage.DBTX) *LoaderFactory {
	return &LoaderFactory{
		db: db,
	}
}

func (f *LoaderFactory) UserInfoLoader() *UserInfoLoader {
	if f.userInfoLoader == nil {
		f.userInfoLoader = NewUserInfoLoader(f.db, nil)
	}
	return f.userInfoLoader
}
