package migrator

import (
	framework2 "boilerplate/internal/framework"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"io/fs"
	"strconv"
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

	FS fs.FS

	SlowQueryMs int `mapstructure:"PG_SLOW_QUERY_LOGGING_LIMIT"`
}

func (c *ModuleConfig) GetDbUrl() string {
	return "postgres://" + c.DbUser + ":" + c.DbPassword + "@" + c.DbHost + ":" + strconv.Itoa(c.DbPort) + "/" + c.DbName + "?sslmode=" + c.DbSslMode
}

func registerRootCommand(
	commands *framework2.Commands,
) {
	commands.Add(
		&cli.Command{
			Name:  "migrator",
			Usage: "Migrate your database",
		},
	)
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"migrator",
		fx.Provide(
			NewMigrate,
			NewAdd,
			NewRollback,
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
		),
		fx.Invoke(
			registerRootCommand,
			RegisterMigrateCommand,
			RegisterAddCommand,
			RegisterRollbackCommand,
		),
	)
}
