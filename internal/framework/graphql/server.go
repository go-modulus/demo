package graphql

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/framework/errors"
	translationContext "boilerplate/internal/translation/context"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type PlaygroundConfig struct {
	Enabled bool   `mapstructure:"GQL_PLAYGROUND_ENABLED"`
	Path    string `mapstructure:"GQL_PLAYGROUND_URL"`
}

type Config struct {
	ComplexityLimit int              `mapstructure:"GQL_COMPLEXITY_LIMIT"`
	Path            string           `mapstructure:"GQL_API_URL"`
	TracingEnabled  bool             `mapstructure:"GQL_TRACING_ENABLED"`
	Playground      PlaygroundConfig `mapstructure:",squash"`
	AppEnv          string           `mapstructure:"APP_ENV"`
}

func NewConfig(viper *viper.Viper) (*Config, error) {
	config := &Config{
		Path: "/graphql",
		Playground: PlaygroundConfig{
			Enabled: true,
			Path:    "/playground",
		},
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode graphql config: %w", err)
	}

	return config, nil
}

type UserError interface {
	ToUserError() map[string]interface{}
}

func NewGraphqlServer(
	config *Config,
	schema graphql.ExecutableSchema,
	loadersInitializer *LoadersInitializer,
	errorHandler *framework.ErrorHandler,
) *handler.Server {
	var mb int64 = 1 << 20

	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(
		transport.MultipartForm{
			MaxUploadSize: mb * 5,
			MaxMemory:     mb * 5,
		},
	)
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New(1000)})
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(config.ComplexityLimit))
	srv.Use(loadersInitializer)
	srv.Use(otelgqlgen.Middleware())

	if config.TracingEnabled {
		srv.Use(apollotracing.Tracer{})
	}

	srv.SetRecoverFunc(
		func(ctx context.Context, p any) error {
			return fmt.Errorf("panic: %v", p)
		},
	)

	srv.SetErrorPresenter(
		func(ctx context.Context, err error) *gqlerror.Error {
			var gqlErr *gqlerror.Error
			path := graphql.GetPath(ctx)
			if errors.As(err, &gqlErr) {
				if gqlErr.Path == nil {
					gqlErr.Path = path
				} else {
					path = gqlErr.Path
				}

				originalErr := gqlErr.Unwrap()
				if originalErr == nil {
					return gqlErr
				}

				err = originalErr
			}

			code := errors.Code(err)
			message := errors.Message(translationContext.GetTranslator(ctx), err)
			extra := errors.Extra(err)
			if extra == nil {
				extra = make(map[string]any)
			}
			extra["code"] = code

			loggable := errors.Loggable(err)
			stack := errors.Stack(err)
			if stack != "" && config.AppEnv != "prod" {
				extra["stack"] = stack
			}

			// Handle only internal server errors (Log them, etc)
			if code == errors.InternalServerErrorCode || loggable {
				errorHandler.Handle(ctx, err)
			}

			return &gqlerror.Error{
				Message:    message,
				Path:       path,
				Extensions: extra,
			}
		},
	)

	return srv
}
