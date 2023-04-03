package main

import (
	"boilerplate/internal/cli"
	"boilerplate/internal/framework"
	"boilerplate/internal/logger"
	"boilerplate/internal/migrator"
	"context"
	"github.com/yalue/merged_fs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

////go:embed db/migrations/*.sql
//var embedFs embed.FS

func main() {
	fsUser := os.DirFS("internal/user/storage/migration")
	fsAuth := os.DirFS("internal/auth/storage/migration")

	mergedFS := merged_fs.MergeMultiple(fsAuth, fsUser)

	app := fx.New(

		framework.NewModule(),
		logger.NewModule(logger.ModuleConfig{}),
		cli.NewModule(cli.ModuleConfig{}),
		migrator.NewModule(migrator.ModuleConfig{FS: mergedFS}),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				lg2 := logger.WithOptions(zap.IncreaseLevel(zapcore.WarnLevel))

				return &fxevent.ZapLogger{Logger: lg2}
			},
		),
	)

	_ = app.Start(context.Background())
}
