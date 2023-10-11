package auth

import (
	"boilerplate/internal/auth/action"
	"boilerplate/internal/auth/provider/local"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/auth/storage/fixture"
	"boilerplate/internal/auth/widget"
	"boilerplate/internal/framework"
	logger2 "github.com/go-pkgz/auth/logger"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"github.com/volatiletech/authboss/v3"
	"github.com/wader/gormstore/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"time"
)

type ModuleConfig struct {
	LocalAccountTable   string `mapstructure:"AUTH_LOCAL_ACCOUNT_TABLE"`
	TokenTable          string `mapstructure:"AUTH_TOKEN_TABLE"`
	SessionIdCookieName string `mapstructure:"AUTH_SESSION_ID_COOKIE_NAME"`
}

func registerRoutes(
	auth *Auth,
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	loginAction *action.LoginAction,
) error {
	authHandler, avatarHandler := auth.service.Handlers()

	err := action.InitLoginAction(routes, errorHandler, loginAction)
	if err != nil {
		return err
	}

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
		local.NewSession,
		action.NewLoginAction,
		action.NewCurrentUser,
		local.NewGormStorage,
		local.NewProvider,
		NewAuthGuardMiddleware,

		widget.NewCurrentUserWidget,

		fixture.NewLocalAccountFixture,
		func(logger framework.Logger) authboss.Logger {
			return newLogger(logger)
		},
		func(cfg *ModuleConfig) *local.ProviderConfig {
			return &local.ProviderConfig{
				AccountTable:        cfg.LocalAccountTable,
				SessionIdCookieName: cfg.SessionIdCookieName,
			}
		},
		func(db *pgxpool.Pool) storage.DBTX {
			return db
		},
		func(db storage.DBTX) *storage.Queries {
			return storage.New(db)
		},
		func(logger framework.Logger) logger2.L {
			return newLogger(logger)
		},
		func(db *gorm.DB) sessions.Store {
			// initialize and setup cleanup
			store := gormstore.NewOptions(
				db,
				gormstore.Options{
					TableName:       "auth.storage",
					SkipCreateTable: false,
				},
				[]byte("secret-hash-key"),
			)
			// some more settings, see sessions.Options
			store.SessionOpts.Secure = false
			store.SessionOpts.HttpOnly = true
			store.SessionOpts.MaxAge = 60 * 60 * 24 * 60

			// db cleanup every hour
			// close quit channel to stop cleanup
			quit := make(chan struct{})
			go store.PeriodicCleanup(1*time.Hour, quit)
			return store
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
