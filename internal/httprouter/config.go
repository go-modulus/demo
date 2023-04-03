package router

import (
	framework2 "boilerplate/internal/framework"
	"context"
	"fmt"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

var corsRegexp *regexp.Regexp

type ModuleConfig struct {
	Port                   int    `mapstructure:"HTTP_APP_PORT"`
	CorsHost               string `mapstructure:"HTTP_CORS_HOST"`
	RedirectTrailingSlash  bool
	RedirectFixedPath      bool
	HandleMethodNotAllowed bool
	HandleOPTIONS          bool
	NotFound               http.Handler
	MethodNotAllowed       http.Handler
	PanicHandler           func(http.ResponseWriter, *http.Request, interface{})
}

func ModuleHooks(
	lc fx.Lifecycle,
	router *Router,
	routes *framework2.Routes,
	logger *zap.Logger,
	corsMiddleware *cors.Cors,
) error {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				router.AddRoutes(routes.GetRoutesInfo())
				logger.Info(fmt.Sprintf("Listen to the port: %d", router.port))
				go func() {

					err := http.ListenAndServe(fmt.Sprintf(":%d", router.port), corsMiddleware.Handler(router.router))
					if err != nil {
						logger.Error(
							fmt.Sprintf(
								"Error while listening the port is occured: %s",
								err.Error(),
							),
						)
					}
				}()
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

func getCorsRegexp(config *ModuleConfig) *regexp.Regexp {
	host := config.CorsHost
	if host == "*" {
		host = ".+"
	}
	if corsRegexp == nil {
		corsRegexp = regexp.MustCompile(host)
	}
	return corsRegexp
}

func createCorsMiddleware(config *ModuleConfig) *cors.Cors {
	corsReg := getCorsRegexp(config)
	return cors.New(
		cors.Options{
			AllowOriginFunc: func(origin string) bool {
				if corsReg == nil {
					return false
				}
				return corsReg.Match([]byte(origin))
			},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodOptions,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"accept",
				"Accept-Encoding",
				"Accept-Language",
				"Authorization",
				"authorization",
				"Content-Type",
				"Content-Length",
				"Cache-Control",
				"Connection",
				"Pragma",
				"Cookie",
				"Access-Control-Allow-Origin",
				"User-Agent",
			},
			ExposedHeaders:     nil,
			MaxAge:             3600,
			AllowCredentials:   true,
			OptionsPassthrough: false,
			Debug:              false,
		},
	)
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"http-router",
		fx.Provide(
			NewRouter,
			createCorsMiddleware,
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
			func(router *Router) framework2.Router {
				return router
			},
		),
		fx.Invoke(
			ModuleHooks,
		),
	)
}
