package framework

import (
	"context"
	"github.com/gofrs/uuid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/http"
)

type contextKey string

func SetHttpRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, contextKey("HttpRequest"), r)
}

func GetHttpRequest(ctx context.Context) *http.Request {
	if value := ctx.Value(contextKey("HttpRequest")); value != nil {
		return value.(*http.Request)
	}
	return nil
}

func SetHttpResponseWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, contextKey("HttpResponseWriter"), w)
}

func GetHttpResponseWriter(ctx context.Context) http.ResponseWriter {
	if value := ctx.Value(contextKey("HttpResponseWriter")); value != nil {
		return value.(http.ResponseWriter)
	}
	return nil
}

func SetLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, contextKey("CurrentLocale"), locale)
}

func GetLocale(ctx context.Context) *string {
	if value := ctx.Value(contextKey("CurrentLocale")); value != nil {
		strVal := value.(string)
		return &strVal
	}
	return nil
}

func SetTranslator(ctx context.Context, translator *message.Printer) context.Context {
	return context.WithValue(ctx, contextKey("Translator"), translator)
}

func GetTranslator(ctx context.Context) *message.Printer {
	if value := ctx.Value(contextKey("Translator")); value != nil {
		return value.(*message.Printer)
	}
	locale := GetLocale(ctx)
	if locale != nil {
		return NewTranslator(*locale)
	}
	return NewTranslator(language.English.String())
}

func SetCurrentUser(ctx context.Context, user *CurrentUser) context.Context {
	if user != nil {
		ctx = SetCurrentUserId(ctx, user.Id)
	}
	return context.WithValue(ctx, contextKey("CurrentUser"), user)
}

func GetCurrentUser(ctx context.Context) *CurrentUser {
	if value := ctx.Value(contextKey("CurrentUser")); value != nil {
		val := value.(*CurrentUser)
		return val
	}
	return nil
}

func SetCurrentUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, contextKey("CurrentUserId"), userId)
}

func GetCurrentUserId(ctx context.Context) *string {
	if value := ctx.Value(contextKey("CurrentUserId")); value != nil {
		strVal := value.(string)
		return &strVal
	}
	return nil
}

func GetCurrentUserUuid(ctx context.Context) uuid.NullUUID {
	if value := ctx.Value(contextKey("CurrentUserId")); value != nil {
		return uuid.NullUUID{UUID: value.(uuid.UUID), Valid: true}
	}
	return uuid.NullUUID{}
}
