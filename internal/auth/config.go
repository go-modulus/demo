package auth

import (
	"boilerplate/internal/framework"
	logger2 "github.com/go-pkgz/auth/logger"
	"github.com/spf13/viper"
	"github.com/volatiletech/authboss/v3"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	AccountTable string `mapstructure:"AUTH_ACCOUNT_TABLE"`
	TokenTable   string `mapstructure:"AUTH_TOKEN_TABLE"`
}

func registerRoutes(
	auth *Auth,
	routes *framework.Routes,
) error {
	authHandler, avatarHandler := auth.service.Handlers()

	routes.Get(
		"/auth/google/login",
		authHandler.ServeHTTP,
	)
	routes.Get(
		"/auth/google/callback",
		authHandler.ServeHTTP,
	)
	routes.Get(
		"/auth/google/logout",
		authHandler.ServeHTTP,
	)

	routes.Get(
		"/auth/local/login",
		authHandler.ServeHTTP,
	)

	routes.Get(
		"/avatar/*",
		avatarHandler.ServeHTTP,
	)

	return nil
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
		func(logger framework.Logger) logger2.L {
			return newLogger(logger)
		},
	}
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"auth",
		fx.Provide(
			append(
				ProvidedServices(),
				func(viper *viper.Viper) (*ModuleConfig, error) {
					err := viper.Unmarshal(&config)
					if err != nil {
						return nil, err
					}
					return &config, nil
				},
			)...,
		),
		fx.Invoke(registerRoutes),
	)
}
