package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	"github.com/go-modulus/demo/internal/auth/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type ResetPasswordRequestLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.ResetPasswordRequest]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.ResetPasswordRequest]
}

func NewResetPasswordRequestLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.ResetPasswordRequest],
) *ResetPasswordRequestLoader {
	if cache == nil {
		cache = &dataloader.NoCache[uuid.UUID, storage.ResetPasswordRequest]{}
	}
	return &ResetPasswordRequestLoader{
		db:    db,
		cache: cache,
	}
}

func (l *ResetPasswordRequestLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.ResetPasswordRequest] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.ResetPasswordRequest] {
				resetPasswordRequestMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.ResetPasswordRequest], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.ResetPasswordRequest]{Data: storage.ResetPasswordRequest{}, Error: err}
						continue
					}

					if loadedItem, ok := resetPasswordRequestMap[key]; ok {
						result[i] = &dataloader.Result[storage.ResetPasswordRequest]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.ResetPasswordRequest]{Data: storage.ResetPasswordRequest{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *ResetPasswordRequestLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.ResetPasswordRequest, error) {
	res := make(map[uuid.UUID]storage.ResetPasswordRequest, len(keys))

	query := `SELECT id, account_id, status, token, last_send_at, used_at, created_at FROM "auth"."reset_password_request" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.ResetPasswordRequest
		err := rows.Scan(
			&result.ID,
			&result.AccountID,
			&result.Status,
			&result.Token,
			&result.LastSendAt,
			&result.UsedAt,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *ResetPasswordRequestLoader) Load(ctx context.Context, resetPasswordRequestKey uuid.UUID) (storage.ResetPasswordRequest, error) {
	return l.getInnerLoader().Load(ctx, resetPasswordRequestKey)()
}
