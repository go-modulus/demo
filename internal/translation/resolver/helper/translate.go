package helper

import (
	"boilerplate/internal/framework/loader"
	translationContext "boilerplate/internal/translation/context"
	translationLoader "boilerplate/internal/translation/resolver/loader"
	"boilerplate/internal/translation/storage"
	"context"
)

func TranslatePointer(
	ctx context.Context,
	loaderFactory *translationLoader.TranslationLoaderFactory,
	key string,
	path storage.Path,
	defaultValue *string,
) (*string, error) {
	tLoader := loader.GetLoader[translationLoader.TranslationId, string](ctx, loaderFactory)

	tag := translationContext.GetLocaleTag(ctx)
	locale, err := storage.LocaleFromTag(tag)
	if err != nil {
		return defaultValue, nil
	}

	translation, err := tLoader.Load(
		ctx, translationLoader.TranslationId{
			Key:    key,
			Path:   path,
			Locale: locale,
		},
	)()
	if err != nil {
		return nil, err
	}
	if translation == "" {
		return defaultValue, nil
	}

	return &translation, nil
}

func TranslateString(
	ctx context.Context,
	loaderFactory *translationLoader.TranslationLoaderFactory,
	key string,
	path storage.Path,
	defaultValue string,
) (string, error) {
	res, err := TranslatePointer(ctx, loaderFactory, key, path, &defaultValue)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}
	return *res, nil
}
