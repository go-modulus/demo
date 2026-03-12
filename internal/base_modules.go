package internal

import (
	"boilerplate/internal/auth_old"
	"boilerplate/internal/blog"
	"boilerplate/internal/cache"
	"boilerplate/internal/framework"
	"boilerplate/internal/gorm"
	"boilerplate/internal/html"
	"boilerplate/internal/logger"
	"boilerplate/internal/pgx"
	"boilerplate/internal/translation"
	"boilerplate/internal/user"
	"go.uber.org/fx"
)

// returns modules that don't listen to some port or don't call cli commands
func Modules() []fx.Option {
	return []fx.Option{
		framework.NewModule(),
		logger.NewModule(logger.ModuleConfig{}),
		framework.HttpModule(),
		pgx.NewModule(pgx.ModuleConfig{}),
		gorm.NewModule(gorm.ModuleConfig{}),
		cache.NewModule(cache.ModuleConfig{}),
		auth_old.NewModule(auth_old.ModuleConfig{}),
		user.NewModule(user.ModuleConfig{}),
		html.NewModule(html.ModuleConfig{}),
		blog.NewModule(blog.ModuleConfig{}),
		translation.NewModule(translation.ModuleConfig{}),
	}
}
