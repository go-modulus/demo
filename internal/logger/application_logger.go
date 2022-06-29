package logger

import (
	"boilerplate/internal/framework"
	"context"
	"go.uber.org/zap"
)

type frameworkLogger struct {
	zapLogger *zap.Logger
}

func NewFrameworkLogger(zapLogger *zap.Logger) framework.Logger {
	return &frameworkLogger{zapLogger: zapLogger}
}

func (a *frameworkLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	a.zapLogger.Sugar().Warnw(s, i...)
}

func (a *frameworkLogger) Info(ctx context.Context, s string, i ...interface{}) {
	a.zapLogger.Sugar().Infow(s, i...)
}

func (a *frameworkLogger) Error(ctx context.Context, s string, i ...interface{}) {
	a.zapLogger.Sugar().Errorw(s, i...)
}

func (a *frameworkLogger) Debug(ctx context.Context, s string, i ...interface{}) {
	a.zapLogger.Sugar().Debugw(s, i...)
}

func (a *frameworkLogger) Panic(ctx context.Context, s string, i ...interface{}) {
	a.zapLogger.Sugar().Panicw(s, i...)
}
