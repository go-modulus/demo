package graphql

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"graphql",
		fx.Provide(
			NewConfig,
			NewGraphqlServer,
			NewHandler,
			NewPlaygroundHandler,
			NewLoadersInitializer,
		),
		fx.Invoke(
			InitHandler,
			InitPlaygroundHandler,
		),
	)
}
