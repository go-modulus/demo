package cli

import (
	"boilerplate/internal/framework"
	"context"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"os"
	"sort"
)

type ModuleConfig struct {
}

func ModuleHooks(
	lc fx.Lifecycle,
	commands *framework.Commands,
	logger framework.Logger,
) error {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				appCLI := &cli.App{
					Usage:                "Run console commands",
					Commands:             commands.GetAll(),
					EnableBashCompletion: true,
				}

				sort.Sort(cli.FlagsByName(appCLI.Flags))
				sort.Sort(cli.CommandsByName(appCLI.Commands))

				err := appCLI.Run(os.Args)
				if err != nil {
					logger.Panic(context.Background(), "failed to run app", err.Error())
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info(context.Background(), "Stopping http-server")
				return nil
			},
		},
	)

	return nil
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"cli",
		fx.Provide(
			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
		),
		fx.Invoke(
			ModuleHooks,
		),
	)
}
