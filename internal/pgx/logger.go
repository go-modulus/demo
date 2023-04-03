package pgx

import (
	"boilerplate/internal/framework"
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type PgxLogger struct {
	logger framework.Logger
	config *ModuleConfig
}

func NewPgxLogger(logger framework.Logger, config *ModuleConfig) *PgxLogger {
	return &PgxLogger{logger: logger, config: config}
}

func (l PgxLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	params := make([]interface{}, 0, 2*len(data))
	for key, val := range data {
		params = append(params, key, val)
	}
	if t, ok := data["time"]; ok {
		if t64, ok2 := t.(time.Duration); ok2 {
			if t64.Milliseconds() > int64(l.config.SlowQueryMs) {
				level = pgx.LogLevelWarn
			}
		}
	}

	switch level {
	case pgx.LogLevelTrace, pgx.LogLevelDebug, pgx.LogLevelInfo:
		l.logger.Debug(ctx, msg, params...)
	case pgx.LogLevelWarn:
		l.logger.Info(ctx, msg, params...)
	case pgx.LogLevelError:
		l.logger.Warn(ctx, msg, params...)
	}
}
