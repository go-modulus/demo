package errors

import (
	"context"
	"go.uber.org/fx"
)

func NewErrChannel(errorHandler *ErrorHandler) chan<- error {
	channel := make(chan error)

	go func() {
		for err := range channel {
			errorHandler.Handle(context.Background(), err)
		}
	}()

	return channel
}

func Module() fx.Option {
	return fx.Module(
		"errors",
		fx.Provide(
			NewErrorHandler,
			fx.Annotate(
				NewErrChannel,
				fx.ResultTags(`name:"errors.channel"`),
			),
		),
	)
}
