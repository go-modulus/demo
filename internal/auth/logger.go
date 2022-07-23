package auth

import (
	"boilerplate/internal/framework"
	"context"
	"fmt"
	"strings"
)

type logger struct {
	fLogger framework.Logger
}

func (l logger) Logf(format string, args ...interface{}) {
	if strings.HasPrefix(format, "[DEBUG]") {
		l.fLogger.Debug(context.Background(), fmt.Sprintf(format, args))
		return
	}
	if strings.HasPrefix(format, "[INFO]") {
		l.fLogger.Info(context.Background(), fmt.Sprintf(format, args))
		return
	}
	if strings.HasPrefix(format, "[WARN]") {
		l.fLogger.Warn(context.Background(), fmt.Sprintf(format, args))
		return
	}
	if strings.HasPrefix(format, "[ERROR]") {
		l.fLogger.Error(context.Background(), fmt.Sprintf(format, args))
		return
	}
	l.fLogger.Info(context.Background(), fmt.Sprintf(format, args))
}

func newLogger(fLogger framework.Logger) *logger {
	return &logger{fLogger: fLogger}
}

func (l logger) Info(s string) {
	l.fLogger.Debug(context.Background(), s)
}

func (l logger) Error(s string) {
	l.fLogger.Error(context.Background(), s)
}
