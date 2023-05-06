package pgx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"strconv"
	"strings"
	"time"
)

type ModuleConfig struct {
	DbHost                  string        `mapstructure:"PG_HOST"`
	DbPort                  int           `mapstructure:"PG_PORT"`
	DbUser                  string        `mapstructure:"PG_USER"`
	DbPassword              string        `mapstructure:"PG_PASSWORD"`
	DbName                  string        `mapstructure:"PG_DB_NAME"`
	DbMinIdleConnections    int           `mapstructure:"PG_MIN_OPEN_CONNECTIONS"`
	DbMaxOpenConnections    int           `mapstructure:"PG_MAX_OPEN_CONNECTIONS"`
	DbMaxConnectionLifetime time.Duration `mapstructure:"PG_MAX_CONNECTION_LIFETIME"`
	DbMaxConnectionIdleTime time.Duration `mapstructure:"PG_MAX_CONNECTION_IDLETIME"`
	DbSslMode               string        `mapstructure:"PG_SSL_MODE"`

	SlowQueryMs int `mapstructure:"PG_SLOW_QUERY_LOGGING_LIMIT"`
}

func (c *ModuleConfig) GetDsn() string {
	return strings.Join(
		[]string{
			"host=" + c.DbHost,
			"port=" + strconv.Itoa(c.DbPort),
			"user=" + c.DbUser,
			"dbname=" + c.DbName,
			"password=" + c.DbPassword,
			"pool_min_conns=" + strconv.Itoa(c.DbMinIdleConnections),
			"pool_max_conns=" + strconv.Itoa(c.DbMaxOpenConnections),
			"pool_max_conn_lifetime=" + c.DbMaxConnectionLifetime.String(),
			"pool_max_conn_idle_time=" + c.DbMaxConnectionIdleTime.String(),
			"prefer_simple_protocol=0",
			"sslmode=" + c.DbSslMode,
		},
		" ",
	)
}

func NewPgxPool(cfg *ModuleConfig, logger *PgxLogger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.GetDsn())
	if err != nil {
		return nil, fmt.Errorf("can`t create pgx pool: %w", err)
	}

	config.ConnConfig.Logger = logger

	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("can`t create pgx pool: %w", err)
	}

	return dbPool, nil
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"pgx",
		fx.Provide(
			NewPgxPool,
			NewPgxLogger,
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				if config.SlowQueryMs == 0 {
					config.SlowQueryMs = 100
				}
				return &config, nil
			},
		),
		fx.Invoke(NewPgxPool),
	)
}
