package loader

import (
	"boilerplate/internal/cache"
	"boilerplate/internal/user/storage"
	"context"
	"github.com/google/uuid"
	"time"
)

var userFinder *storage.Queries

func NewUserLoaderCache(config cache.ModuleConfig) *cache.LoaderCache[uuid.UUID, storage.User] {
	baseCache := cache.NewCache[uuid.UUID, storage.User](
		&cache.Config{
			MaxCachedItems: 10,
			CacheEnabled:   config.CacheEnabled,
			LifeTime:       time.Hour,
		},
	)
	return cache.NewLoaderCache(baseCache)
}

func NewUserLoaderConfig(finder *storage.Queries, cache *cache.LoaderCache[uuid.UUID, storage.User]) UserLoaderConfig {
	if userFinder == nil {
		userFinder = finder
	}
	return UserLoaderConfig{
		Fetch:    fetchUserList,
		Wait:     500 * time.Microsecond,
		MaxBatch: 100,
		Cache:    cache,
	}
}

func fetchUserList(keys []uuid.UUID) ([]storage.User, []error) {
	users := make([]storage.User, len(keys))
	errors := make([]error, len(keys))

	usersMap, err := userFinder.GetUsersMap(context.Background(), keys)
	if err != nil {
		for i := range keys {
			errors[i] = err
		}
	}

	for i, key := range keys {
		if user, ok := usersMap[key.String()]; ok {
			users[i] = user
		}

	}
	return users, errors
}
