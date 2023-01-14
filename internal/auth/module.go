package auth

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module(
		"auth",
		fx.Provide(
			NewMiddleware,
			NewAuth0Provider,

			func(p *Auth0Provider) Provider {
				return p
			},
		),
	)
}
