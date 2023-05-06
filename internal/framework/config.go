package framework

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewViper() (*viper.Viper, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetConfigType("dotenv")
	v.SetConfigName(".env")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	env := v.Get("APP_ENV")
	currentEnv := "dev"
	if envS, ok := env.(string); ok {
		currentEnv = envS
	}

	v.SetConfigName(".env." + currentEnv)
	err = v.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("unable to read config: %w", err)
		}
	}

	return v, nil
}

func NewModule() fx.Option {
	return fx.Module(
		"framework",
		fx.Provide(
			NewViper,
			NewRoutes,
			NewActionRunner,
			NewJsonResponseWriter,
			NewCommands,
			NewAuthenticator,
			NewApiCallerLogger,
			NewErrorHandler,
		),
	)
}
