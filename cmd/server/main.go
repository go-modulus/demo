package main

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/cache"
	"boilerplate/internal/framework"
	router "boilerplate/internal/httprouter"
	"boilerplate/internal/logger"
	"boilerplate/internal/override"
	"boilerplate/internal/pgx"
	"boilerplate/internal/user"
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
		auth.NewModule(auth.ModuleConfig{}),

		user.UserPlugin(),

		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
		override.Overrides(),
	)

	app.Run()
}
