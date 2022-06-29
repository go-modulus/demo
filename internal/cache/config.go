package cache

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	CacheEnabled bool `mapstructure:"CACHE_ENABLED"`
}

func ProvidedServices(config ModuleConfig) []interface{} {
	return []interface{}{
		func(viper *viper.Viper) (*ModuleConfig, error) {
			err := viper.Unmarshal(&config)
			if err != nil {
				return nil, err
			}
			return &config, nil
		},
	}
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"http-router",
		fx.Provide(
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
		),
	)
}
