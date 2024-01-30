package fixture

import (
	"boilerplate/internal/translation/storage"
	"context"
	"fmt"
)

type Translation struct {
	db *storage.Queries
}

func NewTranslation(db *storage.Queries) *Translation {
	return &Translation{db: db}
}

func (f *Translation) CreateTranslation(
	key string,
	path storage.Path,
	locale storage.Locale,
	value string,
) (storage.Translation, func(), string) {
	ctx := context.Background()
	c, _ := f.db.SaveTranslation(
		ctx, storage.SaveTranslationParams{
			Key:    key,
			Path:   path,
			Value:  value,
			Locale: locale,
		},
	)

	return c, func() {
		f.DeleteTranslations(key)
	}, fmt.Sprintf("The translation with key %s", key)
}

func (f *Translation) DeleteTranslations(key string) {
	ctx := context.Background()
	_ = f.db.DeleteTransactionsByKey(ctx, key)
}
