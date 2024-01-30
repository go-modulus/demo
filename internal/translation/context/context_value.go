package context

import (
	"boilerplate/internal/translation"
	"context"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type contextKey string

func SetLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, contextKey("CurrentLocale"), locale)
}

func GetLocaleTag(ctx context.Context) language.Tag {
	locale := ""
	if value := ctx.Value(contextKey("CurrentLocale")); value != nil {
		strVal, ok := value.(string)
		if ok {
			locale = strVal
		}
	}
	return translation.GetSupportedLocale(locale)
}

func SetTranslator(ctx context.Context, translator *message.Printer) context.Context {
	return context.WithValue(ctx, contextKey("Translator"), translator)
}

func GetTranslator(ctx context.Context) *message.Printer {
	if value := ctx.Value(contextKey("Translator")); value != nil {
		return value.(*message.Printer)
	}
	return translation.NewTranslator("")
}
