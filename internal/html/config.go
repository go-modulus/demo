package html

import (
	"boilerplate/internal/framework"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"io/fs"
	"net/http"
	"os"
)

type ModuleConfig struct {
}

func invoke() []any {

	return []any{
		initStaticAction,
	}
}

func initStaticAction(
	routes *framework.Routes,
) error {
	fsys := os.DirFS("static")
	fs.WalkDir(
		fsys, ".", func(p string, d fs.DirEntry, err error) error {
			fmt.Println(p)
			return nil
		},
	)

	fsHandler := http.FileServer(http.FS(fsys))

	routes.Get("/static/*path", http.StripPrefix("/static/", fsHandler).(http.HandlerFunc))

	return nil
}

func providedServices() []interface{} {
	return []any{
		NewIndexPage,
		NewAjaxPage,
	}
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Options(
		fx.Module(
			"html",
			fx.Provide(
				append(
					providedServices(),
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
		),
	)
}
