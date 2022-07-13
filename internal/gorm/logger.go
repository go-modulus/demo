package gorm

import (
	"context"
	application "github.com/debugger84/modulus-application"
	"gorm.io/gorm/logger"
	"runtime/debug"
	"time"
)

type Logger struct {
	application.Logger
	cfg *ModuleConfig
}

func NewGormLogger(cfg *ModuleConfig, logger application.Logger) *Logger {
	return &Logger{Logger: logger, cfg: cfg}
}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	if err != nil && err.Error() != "invalid value" {
		// it is inside the "if" statement only to avoid unnecessary calculations
		//of rows and execution time unlogged queries
		sql, rows := fc()
		if err.Error() == "context canceled" {
			l.Logger.Warn(
				ctx,
				err.Error(),
				"elapsedTime", elapsed,
				"sql", sql,
				"rows", rows,
				"trace", l.getTrace(),
			)
		} else {
			l.Logger.Error(
				ctx,
				err.Error(),
				"elapsedTime", elapsed,
				"sql", sql,
				"rows", rows,
				"trace", l.getTrace(),
			)
		}
	} else if elapsed > time.Duration(l.cfg.SlowQueryLimit)*time.Millisecond ||
		(err != nil && err.Error() == "context canceled") {
		sql, rows := fc()
		l.Logger.Warn(
			ctx,
			"Too long execution",
			"elapsedTime", elapsed,
			"sql", sql,
			"rows", rows,
			"trace", l.getTrace(),
		)
	} else if l.cfg.LoggingEnabled {
		sql, rows := fc()
		l.Logger.Debug(
			ctx,
			"SQL execution",
			"elapsedTime", elapsed,
			"sql", sql,
			"rows", rows,
			"trace", l.getTrace(),
		)
	}
}

func (l Logger) getTrace() string {
	return string(debug.Stack())
}
