package graphql

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"graphql",
		fx.Provide(
			NewConfig,
			NewGraphqlServer,
			NewHandler,
			NewPlaygroundHandler,
			NewLoadersInitializer,
		),
	)
}
