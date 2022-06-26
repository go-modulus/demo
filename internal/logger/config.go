package logger

import (
	"context"
	"demo/internal/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(errorHandler *errors.ErrorHandler) (*zap.Logger, error) {
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

func NewModule() fx.Option {
	return fx.Module(
		"zap-logger", fx.Provide(
			NewLogger,
			NewFrameworkLogger,
		),
	)
}
