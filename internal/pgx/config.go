package pgx

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	PgDsn       string `mapstructure:"PGX_DSN"`
	SlowQueryMs int    `mapstructure:"PGX_SLOW_QUERY_LOGGING_LIMIT"`
}

func ProvidedServices(config ModuleConfig) []interface{} {
	return []interface{}{
		NewPgxPool,
		NewPgxLogger,
		func(viper *viper.Viper) (*ModuleConfig, error) {
			err := viper.Unmarshal(&config)
			if err != nil {
				return nil, err
			}
			return &config, nil
		},
	}
}

func NewPgxPool(cfg *ModuleConfig, logger *PgxLogger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(cfg.PgDsn)
	if err != nil {
		panic("cannot parse pg dsn" + err.Error())
	}
	config.ConnConfig.Logger = logger
	//config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	//	conn.ConnInfo().RegisterDataType(
	//		pgtype.DataType{
	//			Value: &pgtypeuuid.UUID{},
	//			Name:  "uuid",
	//			OID:   pgtype.UUIDOID,
	//		},
	//	)
	//	return nil
	//}
	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		panic("cannot establish connection" + err.Error())
	}

	return dbPool
}

func PgxModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"user",
		fx.Provide(
			ProvidedServices(config)...,
		),
		fx.Invoke(NewPgxPool),
	)
}
