package main

import (
	"boilerplate/internal/cli"
	"boilerplate/internal/framework"
	"boilerplate/internal/logger"
	"boilerplate/internal/migrator"
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	openDir := os.DirFS(".")

	app := fx.New(

		framework.NewModule(),
		logger.NewModule(logger.ModuleConfig{}),
		cli.NewModule(cli.ModuleConfig{}),
		migrator.NewModule(migrator.ModuleConfig{FS: openDir}),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				lg2 := logger.WithOptions(zap.IncreaseLevel(zapcore.WarnLevel))

				return &fxevent.ZapLogger{Logger: lg2}
			},
		),
	)

	_ = app.Start(context.Background())
}
