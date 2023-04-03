package framework

import (
	"context"
)

type ApiCallerLogger struct {
	logger Logger
}

func NewApiCallerLogger(logger Logger) *ApiCallerLogger {
	return &ApiCallerLogger{logger: logger}
}

func (a ApiCallerLogger) Error(msg string, keysAndValues ...interface{}) {
	a.logger.Error(context.Background(), msg, keysAndValues)
}

func (a ApiCallerLogger) Info(msg string, keysAndValues ...interface{}) {
	a.logger.Info(context.Background(), msg, keysAndValues)
}

func (a ApiCallerLogger) Debug(msg string, keysAndValues ...interface{}) {
	a.logger.Debug(context.Background(), msg, keysAndValues)
}

func (a ApiCallerLogger) Warn(msg string, keysAndValues ...interface{}) {
	a.logger.Warn(context.Background(), msg, keysAndValues)
}
