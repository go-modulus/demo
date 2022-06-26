package logger

import (
	"context"
	"demo/internal/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewOriginalZapLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func AddErrorListener(
	errorHandler *errors.ErrorHandler,
	logger Logger,
) {
	errorHandler.AttachListener(
		func(ctx context.Context, err error) error {
			logger.Error(ctx, "Unhandled error", Field("error", err))

			return nil
		},
	)
}

func NewModule() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(
			NewOriginalZapLogger,
			NewZapLogger,
			NewRootEnricher,
		),
		fx.Invoke(AddErrorListener),
	)
}
