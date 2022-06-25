package cache

import (
	application "github.com/debugger84/modulus-application"
	"time"
)

type ModuleConfig struct {
	MaxCacheSizeInMb int
	CacheEnabled     *bool
	LifeTime         time.Duration
}

func (s *ModuleConfig) InitConfig(config application.Config) error {
	if s.CacheEnabled == nil {
		val := config.GetEnvAsBool("CACHE_ENABLED")
		s.CacheEnabled = &val
	}

	return nil
}

func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{}
}

func (s *ModuleConfig) ProvidedServices() []interface{} {
	return []interface{}{
		func() *ModuleConfig {
			return s
		},
	}
}
