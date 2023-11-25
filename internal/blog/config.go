package blog

import (
	"boilerplate/internal/blog/action"
	"boilerplate/internal/blog/page"
	"boilerplate/internal/blog/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
}

func invoke() []any {
	return []any{
		action.InitGetPostsAction,
		page.InitGetPostsPage,
		page.InitAddPostPage,
	}
}

func provide() []any {
	return []any{

		page.NewAddPostPage,
		action.NewGetPostsAction,
		action.NewAddPostAction,

		storage.New,
		func(db *pgxpool.Pool) storage.DBTX {
			return db
		},
	}

}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Module(
		"blog",
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
