package logger

import (
	"boilerplate/internal/framework"
	"context"
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ModuleConfig struct {
	Level             string `mapstructure:"LOGGER_LEVEL"`
	Type              string `mapstructure:"LOGGER_TYPE"`
	AppName           string `mapstructure:"LOGGER_APP"`
	IsProd            bool   `mapstructure:"LOGGER_IS_PROD"`
	DisableStacktrace bool   `mapstructure:"LOGGER_DISABLE_STACKTRACE"`
}

func NewLogger(errorHandler *framework.ErrorHandler, config *ModuleConfig) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, errors.New(
			"invalid logger level " + config.Level +
				". Use \"debug\", \"info\", \"warn\" or \"error\"",
		)
	}
	if config.Type != "json" && config.Type != "console" {
		return nil, errors.New(
			"invalid logger type " + config.Type +
				". Use \"json\" or \"console\"",
		)
	}

	cfg := zap.NewProductionConfig()
	if !config.IsProd {
		cfg.Development = true
		cfg.Sampling = nil
	}
	cfg.Level = level
	cfg.Encoding = config.Type
	if config.Type == "console" {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.InitialFields = map[string]interface{}{
		"app": config.AppName,
	}
	cfg.DisableStacktrace = config.DisableStacktrace

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	errorHandler.AttachListener(
		func(ctx context.Context, err error) error {
			logger.Error("Unhandled error", zap.Error(err))

			return nil
		},
	)

	return logger, nil
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"zap-logger", fx.Provide(
			NewLogger,
			NewFrameworkLogger,

			func(viper *viper.Viper) (*ModuleConfig, error) {
				err := viper.Unmarshal(&config)
				if err != nil {
					return nil, err
				}
				return &config, nil
			},
		),
	)
}
