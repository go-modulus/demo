package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	"github.com/go-modulus/demo/internal/auth/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type SessionLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.Session]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.Session]
}

func NewSessionLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.Session],
) *SessionLoader {
	if cache == nil {
		cache = &dataloader.NoCache[uuid.UUID, storage.Session]{}
	}
	return &SessionLoader{
		db:    db,
		cache: cache,
	}
}

func (l *SessionLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.Session] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.Session] {
				sessionMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.Session], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.Session]{Data: storage.Session{}, Error: err}
						continue
					}

					if loadedItem, ok := sessionMap[key]; ok {
						result[i] = &dataloader.Result[storage.Session]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.Session]{Data: storage.Session{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *SessionLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.Session, error) {
	res := make(map[uuid.UUID]storage.Session, len(keys))

	query := `SELECT id, account_id, identity_id, data, expires_at, created_at FROM "auth"."session" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.Session
		err := rows.Scan(
			&result.ID,
			&result.AccountID,
			&result.IdentityID,
			&result.Data,
			&result.ExpiresAt,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *SessionLoader) Load(ctx context.Context, sessionKey uuid.UUID) (storage.Session, error) {
	return l.getInnerLoader().Load(ctx, sessionKey)()
}
