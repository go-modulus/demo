package middleware

import (
	"log/slog"

	"github.com/go-modulus/auth"
	"github.com/go-modulus/modulus/http"
	"github.com/go-modulus/modulus/http/errhttp"
	"github.com/go-modulus/modulus/http/middleware"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/translation"
	"github.com/rs/cors"
)

type ModuleConfig struct {
	// Add your module configuration here
	// e.g. Var1 string `env:"MYMODULE_VAR1, default=test"`
}

func NewModule() *module.Module {
	return module.NewModule("middleware").
		// Add all dependencies of a module here
		AddDependencies(
			auth.NewModule(),
		).
		// Add all your services here. DO NOT DELETE AddProviders call. It is used for code generation
		AddProviders(
			middleware.NewCors,
			NewCorsMiddlewareFactory,
			NewTranslationMiddlewareFactory,
			NewErrorPipelineFactory,
		).
		// Add all your CLI commands here
		AddCliCommands().
		// Add all your configs here
		InitConfig(ModuleConfig{})
}

type ErrorPipelineFactory struct {
	logger       *slog.Logger
	loggerConfig errhttp.ErrorLoggerConfig
}

func NewErrorPipelineFactory(logger *slog.Logger, loggerConfig errhttp.ErrorLoggerConfig) *ErrorPipelineFactory {
	return &ErrorPipelineFactory{logger: logger, loggerConfig: loggerConfig}
}

func (p *ErrorPipelineFactory) New() *errhttp.ErrorPipeline {
	defPipeline := errhttp.NewDefaultErrorPipeline(p.logger, p.loggerConfig)
	defPipeline.SetProcessor(400, auth.AddHttpCode())

	return defPipeline
}

type CorsMiddlewareFactory struct {
	cors *cors.Cors
}

func NewCorsMiddlewareFactory(cors *cors.Cors) *CorsMiddlewareFactory {
	return &CorsMiddlewareFactory{cors: cors}
}

func (f *CorsMiddlewareFactory) HTTPMiddleware() http.Middleware {
	return f.cors.Handler
}

type TranslationMiddlewareFactory struct {
	tmd *translation.Middleware
}

func NewTranslationMiddlewareFactory(tmd *translation.Middleware) *TranslationMiddlewareFactory {
	return &TranslationMiddlewareFactory{tmd: tmd}
}

func (f *TranslationMiddlewareFactory) HTTPMiddleware() http.Middleware {
	return f.tmd.Middleware
}
