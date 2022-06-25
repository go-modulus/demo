package framework

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type GormError struct {
	Base     Error
	Previous error
}

func NewGormError(err error) *GormError {
	return &GormError{
		Base: Error{
			Message: "Gorm error",
			Code:    "gorm.commonError",
		},
		Previous: err,
	}
}

func (e GormError) Error() string {
	return e.Base.Error()
}

type ZapLogger struct {
	logLevel      logger.LogLevel
	slowThreshold time.Duration
	*zap.Logger
	*zap.SugaredLogger
}

func NewZapLogger(
	originalLogger *zap.Logger,
	slowThreshold time.Duration,
) *ZapLogger {
	return &ZapLogger{
		Logger:        originalLogger,
		SugaredLogger: originalLogger.Sugar(),
		logLevel:      logger.Info,
		slowThreshold: slowThreshold,
	}
}

func (z ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &ZapLogger{
		logLevel:      level,
		slowThreshold: z.slowThreshold,
		SugaredLogger: z.SugaredLogger,
	}
}

func (z ZapLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if z.logLevel < logger.Info {
		return
	}

	z.SugaredLogger.Infow(s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
}

func (z ZapLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if z.logLevel < logger.Warn {
		return
	}

	z.SugaredLogger.Warnw(s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
}

func (z ZapLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if z.logLevel < logger.Error {
		return
	}

	z.SugaredLogger.Errorw(s, append([]interface{}{utils.FileWithLineNum()}, i...)...)
}

func (z ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if z.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && z.logLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		z.Logger.Error(
			"sql error",
			zap.String("file", utils.FileWithLineNum()),
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case z.slowThreshold != 0 && elapsed > z.slowThreshold && z.logLevel >= logger.Warn:
		sql, rows := fc()
		z.Logger.Warn(
			"sql slow execution",
			zap.String("file", utils.FileWithLineNum()),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case z.logLevel >= logger.Info:
		sql, rows := fc()
		z.Logger.Debug(
			"sql execution",
			zap.String("file", utils.FileWithLineNum()),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}

var _ logger.Interface = ZapLogger{}

type GormConfig struct {
	Dsn                  string `mapstructure:"GORM_DSN"`
	PreferSimpleProtocol bool   `mapstructure:"GORM_PREFER_SIMPLE_PROTOCOL"`
	SlowQueryMs          int    `mapstructure:"GORM_SLOW_QUERY_LOGGING_LIMIT"`
}

func NewGorm(
	viper *viper.Viper,
	originalLogger *zap.Logger,
	errorHandler *ErrorHandler,
) (*gorm.DB, error) {
	gormConfig := &GormConfig{
		PreferSimpleProtocol: true,
	}

	err := viper.Unmarshal(&gormConfig)

	if err != nil {
		return nil, fmt.Errorf("unable to decode gorm config: %w", err)
	}

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  gormConfig.Dsn,
				PreferSimpleProtocol: gormConfig.PreferSimpleProtocol,
			},
		),
		&gorm.Config{
			PrepareStmt: false,
			Logger: NewZapLogger(
				originalLogger,
				time.Duration(gormConfig.SlowQueryMs)*time.Millisecond,
			),
		},
	)

	if err != nil {
		return nil, err
	}

	errorHandler.AttachFilter(
		func(_ context.Context, err error) bool {
			_, ok := err.(*GormError)

			return ok == false
		},
	)

	return db, nil
}

func GormModule() fx.Option {
	return fx.Module(
		"gorm",
		fx.Provide(NewGorm),
		fx.Invoke(NewGorm),
	)
}
