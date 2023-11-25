package internal

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/blog"
	"boilerplate/internal/cache"
	"boilerplate/internal/framework"
	"boilerplate/internal/gorm"
	"boilerplate/internal/html"
	"boilerplate/internal/logger"
	"boilerplate/internal/pgx"
	"boilerplate/internal/user"
	"go.uber.org/fx"
)

// returns modules that don't listen to some port or don't call cli commands
func BaseModules() []fx.Option {
	return []fx.Option{
		framework.NewModule(),
		logger.NewModule(logger.ModuleConfig{}),
		framework.HttpModule(),
		pgx.NewModule(pgx.ModuleConfig{}),
		gorm.NewModule(gorm.ModuleConfig{}),
		cache.NewModule(cache.ModuleConfig{}),
		auth.NewModule(auth.ModuleConfig{}),
		user.NewModule(user.ModuleConfig{}),
		html.NewModule(html.ModuleConfig{}),
		blog.NewModule(blog.ModuleConfig{}),
	}
}
