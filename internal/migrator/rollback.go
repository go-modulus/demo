package migrator

import (
	framework2 "boilerplate/internal/framework"
	"context"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"net/url"
)

type Rollback struct {
	cfg    *ModuleConfig
	logger *zap.Logger
}

func NewRollback(cfg *ModuleConfig, logger *zap.Logger) *Rollback {
	return &Rollback{cfg: cfg, logger: logger}
}

func RegisterRollbackCommand(
	command *Rollback,
	commands *framework2.Commands,
) error {
	rootCommand := commands.GetCommandByName("migrator")
	rootCommand.Subcommands = append(
		rootCommand.Subcommands,
		&cli.Command{
			Name:  "rollback",
			Usage: "Rollback the last migration",
			Action: func(cliContext *cli.Context) error {
				return command.Invoke(
					cliContext.Context,
				)
			},
		},
	)

	return nil
}

func (c *Rollback) Invoke(ctx context.Context) error {
	u, _ := url.Parse(c.cfg.GetDbUrl())
	db := dbmate.New(u)
	db.FS = c.cfg.FS
	db.MigrationsDir = "."

	err := db.Rollback()
	if err != nil {
		panic(err)
	}
	return nil
}
