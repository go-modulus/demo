package temporal

import (
	"context"
	"demo/internal/logger"
	"go.temporal.io/sdk/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	Logger logger.Logger
}

func NewLogger(logger logger.Logger) *Logger {
	return &Logger{Logger: logger}
}

func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.Logger.Debug(context.Background(), msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.Logger.Info(context.Background(), msg, keyvals...)
}

func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.Logger.Warn(context.Background(), msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.Logger.Error(context.Background(), msg, keyvals...)
}
