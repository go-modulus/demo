package graphql

import (
	"context"
	"demo/internal/errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type PlaygroundConfig struct {
	Enabled bool `mapstructure:"GRAPHQL_PLAYGROUND_ENABLED"`
	Path    string
}

type Config struct {
	Path       string
	Playground PlaygroundConfig `mapstructure:",squash"`
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
	schema graphql.ExecutableSchema,
	loadersInitializer *LoadersInitializer,
	errorHandler *errors.ErrorHandler,
) *handler.Server {
	srv := handler.New(schema)

	srv.Use(loadersInitializer)

	srv.SetRecoverFunc(func(ctx context.Context, p any) error {
		return errors.FromPanic(p)
	})

	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		originalErr := errors.Unwrap(err)

		var uErr *errors.UserError
		if ep, ok := originalErr.(errors.UserErrorProvider); ok {
			uErr = ep.ToUserError()
		} else {
			uErr = errors.FromError(err).ToUserError()
		}

		var gqlErr *gqlerror.Error
		if errors.As(err, &gqlErr) {
			if gqlErr.Path == nil {
				gqlErr.Path = graphql.GetPath(ctx)
			}

			return gqlErr
		}

		if !uErr.DontHandle {
			errorHandler.Handle(ctx, err)
		}

		return &gqlerror.Error{
			Message:    uErr.Message,
			Path:       gqlErr.Path,
			Extensions: uErr.Extra,
		}
	})

	return srv
}
