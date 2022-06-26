package middleware

import (
	"context"
	"demo/internal/http"
	"demo/internal/logger"
	"github.com/gofrs/uuid"
	oHttp "net/http"
)

const (
	CorrelationIdKey    = "correlationId"
	CorrelationIdHeader = "X-Correlation-ID"
)

type CorrelationIdMiddleware struct{}

func NewCorrelationIdMiddleware(rootEnricher *logger.RootEnricher) *CorrelationIdMiddleware {
	rootEnricher.AttachEnricher(func(ctx context.Context) map[string]string {
		correlationId, ok := ctx.Value(CorrelationIdKey).(string)

		if !ok || correlationId == "" {
			return map[string]string{}
		}

		return map[string]string{
			"correlationId": correlationId,
		}
	})

	return &CorrelationIdMiddleware{}
}

func (CorrelationIdMiddleware) Next(next http.RequestHandler) http.RequestHandler {
	return http.RequestHandlerFunc(
		func(w oHttp.ResponseWriter, req *oHttp.Request) error {
			correlationId := req.Header.Get(CorrelationIdHeader)

			if correlationId == "" {
				id, err := uuid.NewV4()
				if err != nil {
					return err
				}

				correlationId = id.String()
			}

			ctx := context.WithValue(
				req.Context(),
				CorrelationIdKey,
				correlationId,
			)

			w.Header().Set(CorrelationIdHeader, correlationId)

			return next.Handle(w, req.WithContext(ctx))
		},
	)
}
