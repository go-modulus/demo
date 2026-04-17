package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	"github.com/go-modulus/demo/internal/blog/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type PostLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.Post]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.Post]
}

func NewPostLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.Post],
) *PostLoader {
	if cache == nil {
		cache = &dataloader.NoCache[uuid.UUID, storage.Post]{}
	}
	return &PostLoader{
		db:    db,
		cache: cache,
	}
}

func (l *PostLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.Post] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.Post] {
				postMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.Post], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.Post]{Data: storage.Post{}, Error: err}
						continue
					}

					if loadedItem, ok := postMap[key]; ok {
						result[i] = &dataloader.Result[storage.Post]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.Post]{Data: storage.Post{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *PostLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.Post, error) {
	res := make(map[uuid.UUID]storage.Post, len(keys))

	query := `SELECT id, title, preview, content, status, created_at, updated_at, published_at, deleted_at, author_id FROM "blog"."post" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.Post
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Preview,
			&result.Content,
			&result.Status,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.PublishedAt,
			&result.DeletedAt,
			&result.AuthorID,
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *PostLoader) Load(ctx context.Context, postKey uuid.UUID) (storage.Post, error) {
	return l.getInnerLoader().Load(ctx, postKey)()
}
