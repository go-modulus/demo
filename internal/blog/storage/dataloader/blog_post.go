package dataloader

import (
	"context"
	dl "github.com/debugger84/sqlc-dataloader"
	"github.com/go-modulus/demo/internal/blog/storage"
	uuid "github.com/gofrs/uuid"
	"github.com/graph-gophers/dataloader/v7"
)

type BlogPostLoader struct {
	innerLoader *dataloader.Loader[uuid.UUID, storage.BlogPost]
	db          storage.DBTX
	cache       dataloader.Cache[uuid.UUID, storage.BlogPost]
}

func NewBlogPostLoader(
	db storage.DBTX,
	cache dataloader.Cache[uuid.UUID, storage.BlogPost],
) *BlogPostLoader {
	if cache == nil {
		cache = &dataloader.NoCache[uuid.UUID, storage.BlogPost]{}
	}
	return &BlogPostLoader{
		db:    db,
		cache: cache,
	}
}

func (l *BlogPostLoader) getInnerLoader() *dataloader.Loader[uuid.UUID, storage.BlogPost] {
	if l.innerLoader == nil {
		l.innerLoader = dataloader.NewBatchedLoader(
			func(ctx context.Context, keys []uuid.UUID) []*dataloader.Result[storage.BlogPost] {
				blogPostMap, err := l.findItemsMap(ctx, keys)

				result := make([]*dataloader.Result[storage.BlogPost], len(keys))
				for i, key := range keys {
					if err != nil {
						result[i] = &dataloader.Result[storage.BlogPost]{Data: storage.BlogPost{}, Error: err}
						continue
					}

					if loadedItem, ok := blogPostMap[key]; ok {
						result[i] = &dataloader.Result[storage.BlogPost]{Data: loadedItem}
					} else {
						result[i] = &dataloader.Result[storage.BlogPost]{Data: storage.BlogPost{}, Error: dl.ErrNoRows}
					}
				}
				return result
			},
			dataloader.WithCache(l.cache),
		)
	}
	return l.innerLoader
}

func (l *BlogPostLoader) findItemsMap(ctx context.Context, keys []uuid.UUID) (map[uuid.UUID]storage.BlogPost, error) {
	res := make(map[uuid.UUID]storage.BlogPost, len(keys))

	query := `SELECT id, title, preview, content, status, created_at, updated_at, published_at, deleted_at FROM "blog"."post" WHERE id = ANY($1)`
	rows, err := l.db.Query(ctx, query, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result storage.BlogPost
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
		)
		if err != nil {
			return nil, err
		}
		res[result.ID] = result
	}
	return res, nil
}

func (l *BlogPostLoader) Load(ctx context.Context, blogPostKey uuid.UUID) (storage.BlogPost, error) {
	return l.getInnerLoader().Load(ctx, blogPostKey)()
}
