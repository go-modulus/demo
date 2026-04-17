package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	"github.com/go-modulus/demo/internal/auth/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type AccountLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.Account]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.Account]
}

func NewAccountLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.Account],
) *AccountLoader {
	if cache == nil {
		cache = &dataloader.NoCache[uuid.UUID, storage.Account]{}
	}
	return &AccountLoader{
		db:    db,
		cache: cache,
	}
}

func (l *AccountLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.Account] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.Account] {
				accountMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.Account], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.Account]{Data: storage.Account{}, Error: err}
						continue
					}

					if loadedItem, ok := accountMap[key]; ok {
						result[i] = &dataloader.Result[storage.Account]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.Account]{Data: storage.Account{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *AccountLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.Account, error) {
	res := make(map[uuid.UUID]storage.Account, len(keys))

	query := `SELECT id, status, roles, data, updated_at, created_at FROM "auth"."account" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.Account
		err := rows.Scan(
			&result.ID,
			&result.Status,
			&result.Roles,
			&result.Data,
			&result.UpdatedAt,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *AccountLoader) Load(ctx context.Context, accountKey uuid.UUID) (storage.Account, error) {
	return l.getInnerLoader().Load(ctx, accountKey)()
}
