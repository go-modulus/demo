package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	"github.com/go-modulus/demo/internal/auth/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type IdentityLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.Identity]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.Identity]
}

func NewIdentityLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.Identity],
) *IdentityLoader {
	if cache == nil {
		cache = &dataloader.NoCache[uuid.UUID, storage.Identity]{}
	}
	return &IdentityLoader{
		db:    db,
		cache: cache,
	}
}

func (l *IdentityLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.Identity] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.Identity] {
				identityMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.Identity], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.Identity]{Data: storage.Identity{}, Error: err}
						continue
					}

					if loadedItem, ok := identityMap[key]; ok {
						result[i] = &dataloader.Result[storage.Identity]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.Identity]{Data: storage.Identity{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *IdentityLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.Identity, error) {
	res := make(map[uuid.UUID]storage.Identity, len(keys))

	query := `SELECT id, identity, account_id, status, data, updated_at, created_at, type FROM "auth"."identity" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.Identity
		err := rows.Scan(
			&result.ID,
			&result.Identity,
			&result.AccountID,
			&result.Status,
			&result.Data,
			&result.UpdatedAt,
			&result.CreatedAt,
			&result.Type,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *IdentityLoader) Load(ctx context.Context, identityKey uuid.UUID) (storage.Identity, error) {
	return l.getInnerLoader().Load(ctx, identityKey)()
}
