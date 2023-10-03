package user

import (
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/page"
	"boilerplate/internal/user/service"
	"boilerplate/internal/user/storage"
	"boilerplate/internal/user/storage/fixture"
	"boilerplate/internal/user/storage/loader"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
}

func invoke() []any {
	return []any{
		action.InitGetUserAction,
		action.InitGetUsersAction,
		action.InitRegisterAction,
		action.InitUpdateAction,
		page.InitGetUsersPage,
		page.InitNewUserPage,
	}
}

func provide() []any {
	return []any{
		action.NewRegisterAction,
		action.NewGetUserAction,
		action.NewGetUsersAction,
		action.NewUpdateAction,

		page.NewNewUserPage,

		dao.NewUserFinder,
		dao.NewUserSaver,

		service.NewRegistration,

		loader.NewUserLoaderConfig,
		loader.NewUserLoader,
		loader.NewUserLoaderCache,

		fixture.NewUserFixture,

		func(db *pgxpool.Pool) storage.DBTX {
			return db
		},
		func(db storage.DBTX) *storage.Queries {
			return storage.New(db)
		},
	}

}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"user",
		fx.Provide(
			append(
				provide(),
				func(viper *viper.Viper) (*ModuleConfig, error) {
					err := viper.Unmarshal(&config)
					if err != nil {
						return nil, err
					}
					return &config, nil
				},
			)...,
		),
		fx.Invoke(
			invoke()...,
		),
	)
}
