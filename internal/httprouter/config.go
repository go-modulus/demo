package router

import (
	"context"
	"demo/internal/framework"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type ModuleConfig struct {
	Port                   int `mapstructure:"HTTP_APP_PORT"`
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	HandleOPTIONS          bool
	NotFound               http.Handler
	MethodNotAllowed       http.Handler
	PanicHandler           func(http.ResponseWriter, *http.Request, interface{})
	container              *dig.Container
}

func ModuleHooks(
	lc fx.Lifecycle,
	router *Router,
	routes *framework.Routes,
	logger *zap.Logger,
) error {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				router.AddRoutes(routes.GetRoutesInfo())
				logger.Info(fmt.Sprintf("Listen to the port: %d", router.port))
				go http.ListenAndServe(fmt.Sprintf(":%d", router.port), router.router)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Stopping http-server")
				return nil
			},
		},
	)

	return nil
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"http-router",
		fx.Provide(
			NewRouter,
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
			func(router *Router) framework.Router {
				return router
			},
		),
		fx.Invoke(
			ModuleHooks,
		),
	)
}
