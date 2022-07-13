package gorm

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"time"
)

type ModuleConfig struct {
	Dsn                  string        `mapstructure:"GORM_DSN"`
	Dialect              string        `mapstructure:"GORM_DIALECT"`
	PreferSimpleProtocol bool          `mapstructure:"GORM_PREFER_SIMPLE_PROTOCOL"`
	SlowQueryLimit       int           `mapstructure:"GORM_SLOW_QUERY_LOGGING_LIMIT"`
	MaxIdleConns         int           `mapstructure:"GORM_MAX_IDLE_CONNECTIONS"`
	MaxOpenConns         int           `mapstructure:"GORM_MAX_OPEN_CONNECTIONS"`
	ConnMaxLifetime      time.Duration `mapstructure:"GORM_CONN_MAX_LIFETIME"`
	LoggingEnabled       bool          `mapstructure:"GORM_LOGGING_ENABLED"`
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"gorm",
		fx.Provide(
			NewGorm,
			NewGormLogger,

			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
		),
		fx.Invoke(NewGorm),
	)
}
