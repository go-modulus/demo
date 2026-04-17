package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	loaderCache "github.com/debugger84/sqlc-dataloader/cache"
	"github.com/go-modulus/demo/internal/auth/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
	"time"
)

type UserInfoLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.UserInfo]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.UserInfo]
}

func NewUserInfoLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.UserInfo],
) *UserInfoLoader {
	if cache == nil {
		ttl, _ := time.ParseDuration("1m")
		cache = loaderCache.NewLRU[uuid.UUID, storage.UserInfo](100, ttl)
	}
	return &UserInfoLoader{
		db:    db,
		cache: cache,
	}
}

func (l *UserInfoLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.UserInfo] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.UserInfo] {
				userInfoMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.UserInfo], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.UserInfo]{Data: storage.UserInfo{}, Error: err}
						continue
					}

					if loadedItem, ok := userInfoMap[key]; ok {
						result[i] = &dataloader.Result[storage.UserInfo]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.UserInfo]{Data: storage.UserInfo{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *UserInfoLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.UserInfo, error) {
	res := make(map[uuid.UUID]storage.UserInfo, len(keys))

	query := `SELECT id, name, created_at, updated_at FROM "auth"."user_info" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.UserInfo
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *UserInfoLoader) Load(ctx context.Context, userInfoKey uuid.UUID) (storage.UserInfo, error) {
	return l.getInnerLoader().Load(ctx, userInfoKey)()
}
