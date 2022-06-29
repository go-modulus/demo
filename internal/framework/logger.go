package framework

import (
	"context"
	"log"
)

type Logger interface {
	Debug(ctx context.Context, s string, i ...interface{})
	Info(ctx context.Context, s string, i ...interface{})
	Warn(ctx context.Context, s string, i ...interface{})
	Error(ctx context.Context, s string, i ...interface{})
	Panic(ctx context.Context, s string, i ...interface{})
}

type DefaultLogger struct {
}

func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}

func (d *DefaultLogger) Debug(ctx context.Context, s string, i ...interface{}) {
	log.Println("DEBUG: ", s, i)
}

func (d *DefaultLogger) Info(ctx context.Context, s string, i ...interface{}) {
	log.Println("INFO: ", s, i)
}

func (d *DefaultLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	log.Println("WARN: ", s, i)
}

func (d *DefaultLogger) Error(ctx context.Context, s string, i ...interface{}) {
	log.Fatalln("ERROR: ", s, i)
}

func (d *DefaultLogger) Panic(ctx context.Context, s string, i ...interface{}) {
	log.Fatalln("PANIC: ", s, i)
}
