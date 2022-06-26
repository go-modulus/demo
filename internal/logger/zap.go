package logger

import (
	"context"
	"go.uber.org/zap"
)

type ZapLogger struct {
	zapLogger *zap.Logger
	enricher  Enricher
}

func NewZapLogger(zapLogger *zap.Logger, enricher *RootEnricher) Logger {
	return &ZapLogger{zapLogger: zapLogger, enricher: enricher}
}

func (l *ZapLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	i = l.unwrapAndEnrichFields(ctx, i...)

	l.zapLogger.Sugar().Warnw(s, i...)
}

func (l *ZapLogger) Info(ctx context.Context, s string, i ...interface{}) {
	i = l.unwrapAndEnrichFields(ctx, i...)

	l.zapLogger.Sugar().Infow(s, i...)
}

func (l *ZapLogger) Error(ctx context.Context, s string, i ...interface{}) {
	i = l.unwrapAndEnrichFields(ctx, i...)

	l.zapLogger.Sugar().Errorw(s, i...)
}

func (l *ZapLogger) Debug(ctx context.Context, s string, i ...interface{}) {
	i = l.unwrapAndEnrichFields(ctx, i...)

	l.zapLogger.Sugar().Debugw(s, i...)
}

func (l *ZapLogger) Panic(ctx context.Context, s string, i ...interface{}) {
	i = l.unwrapAndEnrichFields(ctx, i...)

	l.zapLogger.Sugar().Panicw(s, i...)
}

func (l *ZapLogger) unwrapAndEnrichFields(ctx context.Context, input ...interface{}) []interface{} {
	output := make([]interface{}, 0, len(input))

	for _, keyOrValue := range input {
		if m, ok := keyOrValue.(map[string]interface{}); ok {
			for k, v := range m {
				output = append(output, k, v)
			}

			continue
		}

		output = append(output, keyOrValue)
	}

	for k, v := range l.enricher.Enrich(ctx) {
		output = append(output, k, v)
	}

	return output
}
