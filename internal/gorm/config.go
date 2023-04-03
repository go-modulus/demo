package gorm

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"strconv"
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

	PgHost                  string        `mapstructure:"PG_HOST"`
	PgPort                  int           `mapstructure:"PG_PORT"`
	PgUser                  string        `mapstructure:"PG_USER"`
	PgPassword              string        `mapstructure:"PG_PASSWORD"`
	PgDbName                string        `mapstructure:"PG_DB_NAME"`
	PgMinIdleConnections    int           `mapstructure:"PG_MIN_OPEN_CONNECTIONS"`
	PgMaxOpenConnections    int           `mapstructure:"PG_MAX_OPEN_CONNECTIONS"`
	PgMaxConnectionLifetime time.Duration `mapstructure:"PG_MAX_CONNECTION_LIFETIME"`
	PgMaxConnectionIdleTime time.Duration `mapstructure:"PG_MAX_CONNECTION_IDLETIME"`
	PgSslMode               string        `mapstructure:"PG_SSL_MODE"`
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
				if (config.Dialect == "" || config.Dsn == "") && config.PgHost != "" {
					config.Dialect = "pgsql"
					config.Dsn = "host=" + config.PgHost + " port=" + strconv.Itoa(config.PgPort) + " user=" +
						config.PgUser + " dbname=" + config.PgDbName + " password=" + config.PgPassword +
						" sslmode=" + config.PgSslMode
				}
				return &config, nil
			},
		),
		fx.Invoke(NewGorm),
	)
}
