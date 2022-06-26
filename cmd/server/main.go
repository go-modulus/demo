package main

import (
	"demo/internal/cache"
	"demo/internal/chi"
	"demo/internal/errors"
	"demo/internal/framework"
	"demo/internal/http"
	"demo/internal/http/middleware"
	"demo/internal/logger"
	"demo/internal/pgx"
	"demo/internal/user"
	"github.com/ggicci/httpin"
	oChi "github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	oHttp "net/http"
	"reflect"
	"time"
)

func main() {
	chiCfg := chi.ModuleParams{
		Configure: func(
			router oChi.Router,
			errorHandler *http.ErrorHandler,
			requestIdMiddleware *middleware.RequestIdMiddleware,
			correlationIdMiddleware *middleware.CorrelationIdMiddleware,
		) {
			m := http.Chain(requestIdMiddleware, correlationIdMiddleware)
			router.Use(func(handler oHttp.Handler) oHttp.Handler {
				return errorHandler.Wrap(m.Next(http.FromHttpHandler(handler)))
			})
			router.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))
			router.Use(chiMiddleware.CleanPath)
			router.Use(chiMiddleware.RealIP)
			router.Use(chiMiddleware.Timeout(5 * time.Second))
			router.Use(chiMiddleware.NewCompressor(4, "application/json").Handler)
			router.NotFound(func(w oHttp.ResponseWriter, req *oHttp.Request) {
				errorHandler.Handle(w, req, errors.NewNotFoundError("http.notFound", "Not Found"))
			})
		},
	}

	app := fx.New(
		framework.ConfigModule(),
		errors.Module(),
		logger.NewModule(),
		chi.Module(chiCfg),
		framework.GormModule(),
		pgx.PgxModule(pgx.ModuleConfig{}),
		cache.NewModule(cache.ModuleConfig{}),
		user.Module(),
		fx.Provide(
			middleware.NewRequestIdMiddleware,
			middleware.NewCorrelationIdMiddleware,
		),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)

	app.Run()
}

func init() {
	httpin.UseGochiURLParam("path", oChi.URLParam)
	httpin.RegisterTypeDecoder(
		reflect.TypeOf(uuid.UUID{}),
		httpin.ValueTypeDecoderFunc(
			func(s string) (interface{}, error) {
				if s == "" {
					return "", nil
				}

				err := validation.Validate(s, is.UUID)
				if err != nil {
					return nil, err
				}

				return uuid.FromString(s)
			},
		),
	)
	httpin.RegisterDirectiveExecutor(
		"ctx",
		httpin.DirectiveExecutorFunc(func(ctx *httpin.DirectiveContext) error {
			ctx.Value.Elem().Set(reflect.ValueOf(ctx.Context))

			return nil
		}),
		nil,
	)
}
