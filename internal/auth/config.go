package auth

import (
	"boilerplate/internal/framework"
	"github.com/volatiletech/authboss/v3"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	UserTable  string
	TokenTable string
}

func ProvidedServices() []interface{} {
	return []interface{}{
		NewAuth,
		NewGormStorage,
		func(logger framework.Logger) authboss.Logger {
			return newLogger(logger)
		},
		func(storage *GormStorage) authboss.ServerStorer {
			return storage
		},
	}
}

func NewModule() fx.Option {
	return fx.Module(
		"auth",
		fx.Provide(
			ProvidedServices()...,
		),
	)
}
