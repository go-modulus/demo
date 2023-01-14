package cli

import (
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func ProvideRoot(root interface{}) interface{} {
	return fx.Annotate(root, fx.ResultTags(`name:"cli.root"`))
}

func ProvideCommand(command interface{}) interface{} {
	return fx.Annotate(command, fx.ResultTags(`group:"cli.commands"`))
}

type StartCliParams struct {
	fx.In

	Lc       fx.Lifecycle
	Root     *cobra.Command   `name:"cli.root"`
	Commands []*cobra.Command `group:"cli.commands"`
}

func NewRootCommand(params StartCliParams) *cobra.Command {
	for _, command := range params.Commands {
		params.Root.AddCommand(command)
	}

	params.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return params.Root.Execute()
			},
		},
	)

	return params.Root
}

func Start(_ *cobra.Command) {}

func Module() fx.Option {
	return fx.Module("cli", fx.Provide(NewRootCommand))
}
