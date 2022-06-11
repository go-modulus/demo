package framework

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(errorHandler *ErrorHandler) (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	errorHandler.AttachListener(func(ctx context.Context, err error) error {
		logger.Error("Unhandled error", zap.Error(err))

		return nil
	})

	return logger, nil
}

func LoggerModule() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(NewLogger),
	)
}
