package chi

import (
	"demo/internal/cli"
	"demo/internal/http"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
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

func NewChi() chi.Router {
	return chi.NewRouter()
}

type ModuleParams struct {
	Configure interface{}
}

func Module(params ModuleParams) fx.Option {
	opts := make([]fx.Option, 0)
	opts = append(
		opts,
		fx.Provide(
			NewConfig,
			NewChi,
			NewServe,

			cli.ProvideCommand(
				func(serve *Serve) *cobra.Command {
					return serve.Command()
				},
			),
		),
	)

	if params.Configure != nil {
		opts = append(
			opts,
			fx.Invoke(params.Configure),
		)
	}

	return http.Module("chi", opts...)
}
