package migrator

import (
	framework2 "boilerplate/internal/framework"
	"context"
	"errors"
	"io/fs"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type RollbackAll struct {
	cfg    *ModuleConfig
	logger *zap.Logger
}

func NewRollbackAll(cfg *ModuleConfig, logger *zap.Logger) *RollbackAll {
	return &RollbackAll{cfg: cfg, logger: logger}
}

func RegisterRollbackAllCommand(
	command *RollbackAll,
	commands *framework2.Commands,
) error {
	rootCommand := commands.GetCommandByName("migrator")
	rootCommand.Subcommands = append(
		rootCommand.Subcommands,
		&cli.Command{
			Name:  "rollback-all",
			Usage: "Rollback all migrations",
			Action: func(cliContext *cli.Context) error {
				return command.Invoke(
					cliContext.Context,
				)
			},
		},
	)

	return nil
}

func (c *RollbackAll) Invoke(ctx context.Context) error {
	u, _ := url.Parse(c.cfg.GetDbUrl())
	db := dbmate.New(u)
	db.FS = c.cfg.FS
	migrationsDir, err := fs.Glob(c.cfg.FS, "internal/*/storage/migration")
	if err != nil {
		panic(err)
	}
	db.MigrationsDir = migrationsDir
	db.AutoDumpSchema = false

	for {
		err = db.Rollback()
		if err != nil {
			if errors.Is(err, dbmate.ErrNoRollback) {
				break
			}
			panic(err)
		}
	}

	return nil
}
