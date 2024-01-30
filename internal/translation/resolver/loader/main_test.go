package loader_test

import (
	"boilerplate/internal/test"
	loader2 "boilerplate/internal/translation/resolver/loader"
	"boilerplate/internal/translation/storage/fixture"
	"go.uber.org/fx"
	"testing"
)

var loaderFactory *loader2.TranslationLoaderFactory
var translationFixture *fixture.Translation

func TestMain(m *testing.M) {
	test.TestMain(
		m, fx.Populate(
			&loaderFactory,
			&translationFixture,
		),
	)
}
