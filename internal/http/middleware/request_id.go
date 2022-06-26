package middleware

import (
	"context"
	"demo/internal/http"
	"demo/internal/logger"
	"github.com/gofrs/uuid"
	oHttp "net/http"
)

const (
	RequestIdKey    = "requestId"
	RequestIdHeader = "X-Request-ID"
)

type RequestIdMiddleware struct {
}

func NewRequestIdMiddleware(rootEnricher *logger.RootEnricher) *RequestIdMiddleware {
	rootEnricher.AttachEnricher(func(ctx context.Context) map[string]string {
		requestId, ok := ctx.Value(RequestIdKey).(string)

		if !ok || requestId == "" {
			return map[string]string{}
		}

		return map[string]string{
			"requestId": requestId,
		}
	})

	return &RequestIdMiddleware{}
}

func (RequestIdMiddleware) Next(next http.RequestHandler) http.RequestHandler {
	return http.RequestHandlerFunc(
		func(w oHttp.ResponseWriter, req *oHttp.Request) error {
			id, err := uuid.NewV4()
			if err != nil {
				return err
			}
			requestId := id.String()

			ctx := context.WithValue(
				req.Context(),
				RequestIdKey,
				requestId,
			)

			w.Header().Set(RequestIdHeader, requestId)

			return next.Handle(w, req.WithContext(ctx))
		},
	)
}
