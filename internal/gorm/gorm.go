package gorm

import (
	"boilerplate/internal/framework"
	"context"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormError struct {
	Base     framework.Error
	Previous error
}

func NewGormError(err error) *GormError {
	return &GormError{
		Base: framework.Error{
			Message: "Gorm error",
			Code:    "gorm.commonError",
		},
		Previous: err,
	}
}

func (e GormError) Error() string {
	return e.Base.Error()
}

var _ logger.Interface = Logger{}

func NewGorm(
	cfg *ModuleConfig,
	gormLogger *Logger,
	errorHandler *framework.ErrorHandler,
) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	if cfg.PreferSimpleProtocol {
		gormConfig.PrepareStmt = false
	}

	gormConfig.Logger = gormLogger

	var dialector gorm.Dialector
	switch cfg.Dialect {
	case "pgsql":
		dialector = postgres.New(
			postgres.Config{
				DSN:                  cfg.Dsn,
				PreferSimpleProtocol: cfg.PreferSimpleProtocol,
			},
		)
	case "mysql":
		dialector = mysql.New(
			mysql.Config{
				DSN: cfg.Dsn,
			},
		)
	default:
		return nil, errors.New(
			cfg.Dialect + " is unsupported DB dialect. " +
				"Use mysql or pgsql in the GORM_DIALECT environment variable",
		)
	}

	db, err := gorm.Open(
		dialector,
		gormConfig,
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

	if db != nil {
		sqlDB, err := db.DB()

		if err != nil {
			return nil, err
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		if cfg.ConnMaxLifetime.Seconds() == 0 {
			cfg.ConnMaxLifetime = time.Hour
		}
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		if cfg.LoggingEnabled {
			db = db.Debug()
		}
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
