package translation

import (
	"boilerplate/internal/translation/resolver/loader"
	"boilerplate/internal/translation/storage"
	"boilerplate/internal/translation/storage/fixture"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
}

func ProvidedServices() []interface{} {
	return []interface{}{
		loader.NewTranslationLoaderFactory,

		fixture.NewTranslation,
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
		"translation",

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
	)
}
