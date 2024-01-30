package loader

import (
	"context"
	"fmt"
	"github.com/graph-gophers/dataloader/v7"
	"sync"
)

type contextKey string

const (
	LoadersKey contextKey = "Loaders"
)

type LoaderFactory[K comparable, T any] interface {
	Create() *dataloader.Loader[K, T]
}

func WithLoaders(ctx context.Context, loaders *sync.Map) context.Context {
	return context.WithValue(ctx, LoadersKey, loaders)
}

func getLoaders(ctx context.Context) *sync.Map {
	if loaders, ok := ctx.Value(LoadersKey).(*sync.Map); ok {
		return loaders
	}
	return nil
}

func GetLoader[K comparable, T any](ctx context.Context, factory LoaderFactory[K, T]) *dataloader.Loader[K, T] {
	loaders := getLoaders(ctx)
	key := fmt.Sprintf("%T", factory)
	existLoader, ok := loaders.Load(key)
	if ok {
		if loader, ok := existLoader.(*dataloader.Loader[K, T]); ok {
			return loader
		}

		var (
			loaderKey   K
			loaderValue T
		)
		panic(fmt.Sprintf("loader with key %s is not of type *dataloader.Loader[%T, %T]", key, loaderKey, loaderValue))
	}

	loader := factory.Create()
	loaders.Store(key, loader)
	return loader
}
