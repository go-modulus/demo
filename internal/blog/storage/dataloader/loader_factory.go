package dataloader

import (
	"github.com/go-modulus/demo/internal/blog/storage"
)

type LoaderFactory struct {
	db         storage.DBTX
	postLoader *PostLoader
}

func NewLoaderFactory(db storage.DBTX) *LoaderFactory {
	return &LoaderFactory{
		db: db,
	}
}

func (f *LoaderFactory) PostLoader() *PostLoader {
	if f.postLoader == nil {
		f.postLoader = NewPostLoader(f.db, nil)
	}
	return f.postLoader
}
