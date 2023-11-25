package migrator

import (
	framework2 "boilerplate/internal/framework"
	"context"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type Migrate struct {
	cfg    *ModuleConfig
	logger *zap.Logger
}

func NewMigrate(cfg *ModuleConfig, logger *zap.Logger) *Migrate {
	return &Migrate{cfg: cfg, logger: logger}
}

func RegisterMigrateCommand(
	command *Migrate,
	commands *framework2.Commands,
) error {
	rootCommand := commands.GetCommandByName("migrator")
	rootCommand.Subcommands = append(
		rootCommand.Subcommands,
		&cli.Command{
			Name:  "migrate",
			Usage: "Apply all migrations from the registered modules to the database",
			Action: func(cliContext *cli.Context) error {
				return command.Invoke(
					cliContext.Context,
				)
			},
		},
	)

	return nil
}

func (c *Migrate) Invoke(ctx context.Context) error {
	u, _ := url.Parse(c.cfg.GetDbUrl())
	db := dbmate.New(u)
	db.FS = c.cfg.FS
	db.AutoDumpSchema = false
	migrationsDir, err := fs.Glob(c.cfg.FS, "internal/*/storage/migration")
	if err != nil {
		panic(err)
	}
	db.MigrationsDir = migrationsDir

	fmt.Println("\nApplying...")
	err = db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
	return nil
}
