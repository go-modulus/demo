package main

import (
	"github.com/go-modulus/auth"
	"github.com/go-modulus/auth/providers/email"
	auth2 "github.com/go-modulus/demo/internal/auth"
	"github.com/go-modulus/demo/internal/blog"
	graphql2 "github.com/go-modulus/demo/internal/graphql"
	"github.com/go-modulus/graphql"
	"github.com/go-modulus/modulus/captcha"
	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/config"
	"github.com/go-modulus/modulus/http"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
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
					Version: "0.1.0",
					Usage:   "Run project commands",
				},
			),
		),
		logger.NewModule(),
		pgx.NewModule(),
		migrator.NewModule(),
		http.NewModule(
			http.AddMiddlewareFactoryToPipeline[*auth.Middleware](500),
		),
		graphql.NewModule(),
		graphql2.NewModule(),
		blog.NewModule(),
		captcha.NewModule(),
		auth.NewModule(),
		auth2.NewModule(),
		email.NewModule(),
	}

	app := fx.New(
		module.BuildFx(modules...),
		logger.FxLoggerOption(),
		cli.InvokeStartCli(),
	)

	app.Run()
}
