package html

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
}

func invoke() []any {
	return []any{}
}

func providedServices() []interface{} {
	return []any{
		NewIndexPage,
		NewAjaxPage,
	}
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Options(
		fx.Module(
			"auth",
			fx.Provide(
				append(
					providedServices(),
					func(viper *viper.Viper) (*ModuleConfig, error) {
						err := viper.Unmarshal(&config)
						if err != nil {
							return nil, err
						}
						return &config, nil
					},
				)...,
			),
			fx.Invoke(
				invoke()...,
			),
		),
	)
}
