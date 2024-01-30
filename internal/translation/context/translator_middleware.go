package context

import (
	"boilerplate/internal/translation"
	"net/http"
)

// TranslationMiddleware is a middleware that adds to a context a translator
// that can be used to translate messages.
// It uses the locale from the Accept-Language header to save
// a language tag to the context to determine the language of the translation.
func TranslationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		locale := req.Header.Get("Accept-Language")
		ctx = SetLocale(ctx, locale)
		p := translation.NewTranslator(locale)
		ctx = SetTranslator(ctx, p)

		req = req.WithContext(ctx)
		next(w, req)
	}
}
