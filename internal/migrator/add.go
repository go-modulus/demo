package migrator

import (
	framework2 "boilerplate/internal/framework"
	"context"
	"fmt"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"net/url"
)

type Add struct {
	cfg    *ModuleConfig
	logger *zap.Logger
}

func NewAdd(cfg *ModuleConfig, logger *zap.Logger) *Add {
	return &Add{cfg: cfg, logger: logger}
}

func RegisterAddCommand(
	command *Add,
	commands *framework2.Commands,
) error {
	rootCommand := commands.GetCommandByName("migrator")
	rootCommand.Subcommands = append(
		rootCommand.Subcommands,
		&cli.Command{
			Name:  "add",
			Usage: "Add a new migration to the module",
			Action: func(cliContext *cli.Context) error {
				return command.Invoke(
					cliContext.Context,
					cliContext.String("module"),
					cliContext.String("name"),
				)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "module",
					Usage:    "The module to add the migration to",
					Required: true,
					Aliases:  []string{"m"},
				},
				&cli.StringFlag{
					Name:     "name",
					Usage:    "The name of migration to add",
					Required: true,
					Aliases:  []string{"n"},
				},
			},
		},
	)

	return nil
}

func (c *Add) Invoke(ctx context.Context, module string, name string) error {
	u, _ := url.Parse(c.cfg.GetDbUrl())
	db := dbmate.New(u)
	db.MigrationsDir = "internal/" + module + "/storage/migration"

	fmt.Println("Add a migration to the dir:" + db.MigrationsDir)
	err := db.NewMigration(name)
	if err != nil {
		return err
	}

	fmt.Println("\nMigration is created.")

	return nil
}
