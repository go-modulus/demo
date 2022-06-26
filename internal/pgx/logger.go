package pgx

import (
	"context"
	"demo/internal/logger"
	"github.com/jackc/pgx/v4"
	"time"
)

type PgxLogger struct {
	logger logger.Logger
	config *ModuleConfig
}

func NewPgxLogger(logger logger.Logger, config *ModuleConfig) *PgxLogger {
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
	case 6, 5, 4:
		l.logger.Debug(ctx, msg, params...)
	case 3:
		l.logger.Info(ctx, msg, params...)
	case 2:
		l.logger.Warn(ctx, msg, params...)
	}
}
