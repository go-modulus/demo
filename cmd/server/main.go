package main

import (
	"boilerplate/internal/framework"
	router "boilerplate/internal/httprouter"
	"boilerplate/internal/logger"
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
		logger.NewZapLoggerModule(),
		framework.HttpModule(),
		framework.GormModule(),
		pgx.PgxModule(pgx.ModuleConfig{}),
		router.HttpRouterModule(router.ModuleConfig{}),
		user.UserPlugin(),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)

	app.Run()
}
