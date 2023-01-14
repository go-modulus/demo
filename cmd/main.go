package main

import (
	"demo/graph"
	"demo/internal/auth"
	"demo/internal/chi"
	"demo/internal/cli"
	"demo/internal/errors"
	"demo/internal/framework"
	"demo/internal/graphql"
	"demo/internal/http"
	"demo/internal/http/middleware"
	"demo/internal/logger"
	"demo/internal/messenger"
	"demo/internal/pgx"
	"demo/internal/temporal"
	"github.com/ggicci/httpin"
	oChi "github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	oHttp "net/http"
	"reflect"
	"strings"
	"time"
)

func main() {
	chiCfg := chi.ModuleParams{
		Configure: func(
			router oChi.Router,
			errorHandler *http.ErrorHandler,
			requestIdMiddleware *middleware.RequestIdMiddleware,
			correlationIdMiddleware *middleware.CorrelationIdMiddleware,
			authMiddleware *auth.Middleware,
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
			router.Use(func(handler oHttp.Handler) oHttp.Handler {
				return errorHandler.Wrap(authMiddleware.Next(http.FromHttpHandler(handler)))
			})
			router.NotFound(func(w oHttp.ResponseWriter, req *oHttp.Request) {
				errorHandler.Handle(w, req, errors.NewNotFoundError("http.notFound", "Not Found"))
			})
		},
	}

	app := fx.New(
		framework.ConfigModule(),
		errors.Module(),
		logger.NewModule(),
		cli.Module(),
		chi.Module(chiCfg),
		pgx.Module(pgx.ModuleConfig{}),
		graph.Module(),
		graphql.Module(),
		temporal.Module(),
		auth.Module(),
		messenger.Module(),
		fx.Provide(
			middleware.NewRequestIdMiddleware,
			middleware.NewCorrelationIdMiddleware,
			cli.ProvideRoot(
				func(shutdowner fx.Shutdowner) *cobra.Command {
					root := &cobra.Command{
						Use:   "modulus",
						Short: "Modulus is modulus",
						Long:  `Modulus is GoLang framework for building modular applications.`,
						Args:  cobra.NoArgs,
					}

					root.InitDefaultHelpCmd()
					root.InitDefaultHelpFlag()
					root.InitDefaultVersionFlag()

					return root
				},
			),
			cli.ProvideCommand(
				func() *cobra.Command {
					return &cobra.Command{
						Use:   "echo [string to echo]",
						Short: "Echo anything to the screen",
						Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
						Args: cobra.MinimumNArgs(1),
						Run: func(cmd *cobra.Command, args []string) {
							cmd.Println("Print: " + strings.Join(args, " "))
						},
					}
				},
			),
		),
		fx.Invoke(cli.Start),
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
