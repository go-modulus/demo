package http

import (
	"go.uber.org/fx"
)

func Module(name string, opts ...fx.Option) fx.Option {
	opts = append(
		opts,
		fx.Provide(NewRoutes, NewErrorHandler),
		fx.Invoke(
			func(errorHandler *ErrorHandler) {
				errorHandler.AttachTransformer(UnwrapHttpinErrors)
			},
			RegisterHandlers,
		),
	)

	return fx.Module(name, opts...)
}
