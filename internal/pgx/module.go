package pgx

import (
	"context"
	"database/sql"
	"demo/internal/logger"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	PgDsn       string `mapstructure:"PGX_DSN"`
	SlowQueryMs int    `mapstructure:"PGX_SLOW_QUERY_LOGGING_LIMIT"`
}

func NewDB(cfg *ModuleConfig, logger logger.Logger) (*sql.DB, error) {
	config, err := pgx.ParseConfig(cfg.PgDsn)
	if err != nil {
		return nil, fmt.Errorf("cannot parse pg config: %w", err)
	}

	config.Tracer = &tracelog.TraceLog{
		Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
			logger.Info(ctx, msg, data)
		}),
		LogLevel: tracelog.LogLevelTrace,
	}

	db := stdlib.OpenDB(*config)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	return db, nil
}

func Module(config ModuleConfig) fx.Option {
	return fx.Module(
		"pgx",
		fx.Provide(
			NewDB,
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
