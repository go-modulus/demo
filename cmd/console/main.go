package main

import (
	"github.com/go-modulus/auth"
	"github.com/go-modulus/auth/providers/email"
	auth2 "github.com/go-modulus/demo/internal/auth"
	"github.com/go-modulus/demo/internal/blog"
	"github.com/go-modulus/demo/internal/common/middleware"
	graphql2 "github.com/go-modulus/demo/internal/graphql"
	"github.com/go-modulus/graphql"
	"github.com/go-modulus/modulus/captcha"
	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/config"
	"github.com/go-modulus/modulus/http"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/otel"
	"github.com/go-modulus/pgx"
	"github.com/go-modulus/pgx/migrator"
	"go.uber.org/fx"
)

func main() {
	config.LoadDefaultEnv()

	// DO NOT Remove. It will be edited by the `mtools module create` CLI command.
	modules := []*module.Module{
		cli.NewModule(
			cli.SetConfig(
				cli.ModuleConfig{
					Version:        "0.1.0",
					Usage:          "Run project commands",
					DefaultCommand: "serve",
				},
			),
		),
		logger.NewModule(
			logger.AddMiddlewareFactoryToPipeline[*otel.LogMiddlewareFactory](400),
		),
		pgx.NewModule(),
		migrator.NewModule(),
		http.NewModule(
			http.AddMiddlewareFactoryToPipeline[*auth.Middleware](500),
			http.AddMiddlewareToPipeline(600, otel.NewMiddleware("demo-http")),
			http.OverrideErrorPipeline[*middleware.ErrorPipelineFactory],
		),
		graphql.NewModule(
			graphql.AddInitFuncFactory[*auth.GraphQLInitFuncFactory](100),
		),
		graphql2.NewModule(),
		blog.NewModule(),
		captcha.NewModule(),
		auth.NewModule(),
		auth2.NewModule(),
		email.NewModule(),
		otel.NewModule(),
		middleware.NewModule(),
	}

	app := fx.New(
		module.BuildFx(modules...),
		logger.FxLoggerOption(),
		cli.InvokeStartCli(),
	)

	app.Run()
}
