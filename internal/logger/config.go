package logger

import (
	"boilerplate/internal/framework"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(errorHandler *framework.ErrorHandler) (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	errorHandler.AttachListener(
		func(ctx context.Context, err error) error {
			logger.Error("Unhandled error", zap.Error(err))

			return nil
		},
	)

	return logger, nil
}

func NewZapLoggerModule() fx.Option {
	return fx.Module(
		"zap-logger", fx.Provide(
			NewLogger,
			NewFrameworkLogger,
		),
	)
}
