package auth

import (
	"boilerplate/internal/framework"
	"context"
)

type logger struct {
	fLogger framework.Logger
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
