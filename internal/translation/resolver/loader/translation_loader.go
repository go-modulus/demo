package loader

import (
	"boilerplate/internal/translation/storage"
	"context"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/jackc/pgtype"
)

type TranslationResponse = []*Translation
type Translation = dataloader.Result[string]
type TranslationId struct {
	Key    string         `json:"key"`
	Path   storage.Path   `json:"path"`
	Locale storage.Locale `json:"locale"`
}

type TranslationLoaderFactory struct {
	queries *storage.Queries
}

func NewTranslationLoaderFactory(queries *storage.Queries) *TranslationLoaderFactory {
	return &TranslationLoaderFactory{queries: queries}
}

func (l *TranslationLoaderFactory) Create() *dataloader.Loader[TranslationId, string] {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, translationIds []TranslationId) TranslationResponse {
			result := make(TranslationResponse, len(translationIds))
			keys := pgtype.JSONB{}
			err := keys.Set(translationIds)
			if err != nil {
				return result
			}
			translations, err := l.queries.FindTranslations(ctx, keys)
			if err != nil {
				return result
			}
			translationMap := make(map[TranslationId]string)
			for _, translation := range translations {
				translationMap[TranslationId{
					Key:    translation.Key,
					Path:   translation.Path,
					Locale: translation.Locale,
				}] = translation.Value
			}

			for i, translationId := range translationIds {
				value := ""

				if translation, ok := translationMap[translationId]; ok {
					value = translation
				}
				result[i] = &Translation{Data: value}
			}
			return result
		},
	)
}
