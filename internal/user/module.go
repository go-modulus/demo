package user

import (
	"demo/internal/user/dao"
	"demo/internal/user/http"
	"demo/internal/user/service"
	"demo/internal/user/storage"
	"demo/internal/user/storage/loader"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"user",
		fx.Provide(
			http.NewRegisterAction,
			http.NewGetUserAction,
			http.NewGetUsersAction,
			http.NewUpdateAction,

			dao.NewUserFinder,
			dao.NewUserSaver,

			service.NewRegistration,

			loader.NewUserLoaderConfig,
			loader.NewUserLoader,
			loader.NewUserLoaderCache,

			func(db *pgxpool.Pool) storage.DBTX {
				return db
			},
			func(db storage.DBTX) *storage.Queries {
				return storage.New(db)
			},
		),
	)
}
