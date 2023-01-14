package temporal

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	HostPort string `mapstructure:"TEMPROAL_HOSTPORT"`
}

func NewConfig(viper *viper.Viper) (*Config, error) {
	config := &Config{
		HostPort: "localhost:7233",
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode temporal config: %w", err)
	}

	return config, nil
}
