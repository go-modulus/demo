package loader_test

import (
	"boilerplate/internal/test/expect"
	"boilerplate/internal/test/spec"
	loader2 "boilerplate/internal/translation/resolver/loader"
	"boilerplate/internal/translation/storage"
	"context"
	"testing"
)

func TestTranslationLoader_Load(t *testing.T) {
	t.Run(
		"should load translation", func(t *testing.T) {

			_, rb, givenTranslation := translationFixture.CreateTranslation(
				"test",
				storage.PathAdminconfigvalue,
				storage.LocaleEn,
				"test translation",
			)
			defer rb()

			loader := loaderFactory.Create()

			translation, err := loader.Load(
				context.Background(), loader2.TranslationId{
					Key:    "test",
					Path:   storage.PathAdminconfigvalue,
					Locale: storage.LocaleEn,
				},
			)()

			spec.Given(t, givenTranslation)
			spec.When(t, "translation is loaded")
			spec.NoError(t, "No error", err)
			spec.Then(
				t, "translation is returned",
				expect.Equal("test translation", translation),
			)
		},
	)
}
