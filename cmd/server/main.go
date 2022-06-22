package main

import (
	"demo/internal/cache"
	"demo/internal/framework"
	router "demo/internal/httprouter"
	"demo/internal/logger"
	"demo/internal/pgx"
	"demo/internal/user"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		framework.ConfigModule(),
		framework.ErrorsModule(),
		logger.NewModule(),
		framework.HttpModule(),
		framework.GormModule(),
		pgx.PgxModule(pgx.ModuleConfig{}),
		cache.NewModule(cache.ModuleConfig{}),
		router.NewModule(router.ModuleConfig{}),
		user.UserPlugin(),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)

	app.Run()
}
