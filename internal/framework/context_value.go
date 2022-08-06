package framework

import (
	"context"
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
