package chi

import (
	"context"
	"demo/internal/http"
	"demo/internal/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	oHttp "net/http"
)

type Config struct {
	Address string `mapstructure:"HTTP_ADDRESS"`
}

func NewConfig(viper *viper.Viper) (*Config, error) {
	config := &Config{
		Address: "127.0.0.1:8080",
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode chi config: %w", err)
	}

	return config, nil
}

type ChiParams struct {
	fx.In

	Lc           fx.Lifecycle
	Config       *Config
	Routes       *http.Routes
	ErrorHandler *http.ErrorHandler
	ErrChannel   chan<- error `name:"errors.channel"`
	Logger       logger.Logger
}

func NewChi(params ChiParams) (chi.Router, error) {
	router := chi.NewRouter()

	server := oHttp.Server{
		Addr:    params.Config.Address,
		Handler: router,
		// TODO: Add logger
		//ErrorLog: logger,
	}

	params.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, route := range params.Routes.List() {
				handler := route.Handler

				router.Method(
					route.Method,
					route.Path,
					params.ErrorHandler.Wrap(handler),
				)
			}

			params.Logger.Info(
				ctx,
				"Starting http-server",
				logger.Field("address", params.Config.Address),
			)

			go func() {
				err := server.ListenAndServe()
				params.ErrChannel <- err
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			params.Logger.Info(ctx, "Stopping http-server")
			return server.Shutdown(ctx)
		},
	})

	return router, nil
}

type ModuleParams struct {
	Configure interface{}
}

func Module(params ModuleParams) fx.Option {
	opts := make([]fx.Option, 0)
	opts = append(opts, fx.Provide(NewConfig, NewChi))

	if params.Configure != nil {
		opts = append(
			opts,
			fx.Invoke(params.Configure),
		)
	}

	return http.Module("chi", opts...)
}
