package user

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/framework"
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/httpaction"
	"boilerplate/internal/user/service"
	"boilerplate/internal/user/storage"
	"boilerplate/internal/user/storage/loader"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
)

func registerRouters(
	chi chi.Router,
	errorHandler *framework.HttpErrorHandler,
	registerAction *action.RegisterAction,
	getUserAction *action.GetUserAction,
	getUsersAction *action.GetUsersAction,
	updateAction *action.UpdateAction,
	genActions *httpaction.ModuleActions,
	routes *framework.Routes,
	auth *auth.Auth,
) error {
	err := registerAction.Register(routes, errorHandler)
	if err != nil {
		return err
	}
	err = getUserAction.Register(auth, routes, errorHandler)
	if err != nil {
		return err
	}
	err = getUsersAction.Register(auth, routes, errorHandler)
	if err != nil {
		return err
	}
	err = updateAction.Register(chi, errorHandler)
	if err != nil {
		return err
	}

	routes.AddFromRoutes(genActions.Routes())
	return nil
}

func ProvidedServices() []interface{} {
	return append(
		httpaction.ServiceProviders(),
		[]interface{}{
			action.NewRegisterAction,
			action.NewGetUserAction,
			action.NewGetUsersAction,
			action.NewUpdateAction,

			dao.NewUserFinder,
			dao.NewUserSaver,

			service.NewRegistration,

			httpaction.NewRegisterAction,
			httpaction.NewGetUsersAction,

			httpaction.NewGetUserProcessor,
			httpaction.NewUpdateProcessor,

			loader.NewUserLoaderConfig,
			loader.NewUserLoader,
			loader.NewUserLoaderCache,

			func(db *pgxpool.Pool) storage.DBTX {
				return db
			},
			func(db storage.DBTX) *storage.Queries {
				return storage.New(db)
			},
			func() httpaction.TestOverride {
				return &httpaction.Override{}
			},
		}...,
	)
}

func UserPlugin() fx.Option {
	return fx.Module(
		"user",
		fx.Provide(
			ProvidedServices()...,
		),
		fx.Invoke(registerRouters),
	)
}
