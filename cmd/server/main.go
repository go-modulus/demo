package main

import (
	"boilerplate/internal"
	router "boilerplate/internal/httprouter"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		append(
			internal.Modules(),

			router.NewModule(router.ModuleConfig{}),
			fx.WithLogger(
				func(logger *zap.Logger) fxevent.Logger {
					return &fxevent.ZapLogger{Logger: logger}
				},
			),
		)...,
	)

	app.Run()
}
