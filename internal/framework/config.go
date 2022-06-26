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

	return v, nil
}

func ConfigModule() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(NewViper),
	)
}
