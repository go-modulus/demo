package auth

import (
	"context"
	"demo/internal/http"
	"demo/internal/logger"
	oHttp "net/http"
)

type Middleware struct {
	provider Provider
}

func NewMiddleware(rootEnricher *logger.RootEnricher, provider Provider) *Middleware {
	rootEnricher.AttachEnricher(func(ctx context.Context) map[string]string {
		performer := PerformerFromContext(ctx)

		if !performer.Valid {
			return map[string]string{}
		}

		return map[string]string{
			"performerId": performer.Value.Id.String(),
		}
	})

	return &Middleware{provider: provider}
}

func (m *Middleware) Next(next http.RequestHandler) http.RequestHandler {
	f := func(w oHttp.ResponseWriter, req *oHttp.Request) error {
		header := req.Header.Get("Authorization")
		if header == "" {
			return next.Handle(w, req)
		}

		token := header[len("Bearer "):]
		performer, err := m.provider.GetUser(req.Context(), token)
		if err != nil {
			return err
		}

		if !performer.Valid {
			return next.Handle(w, req)
		}

		ctx := ContextWithPerformer(req.Context(), performer.Value)

		return next.Handle(w, req.WithContext(ctx))
	}

	return http.RequestHandlerFunc(f)
}
