package template

import (
	"boilerplate/internal/framework"
	"embed"
	"io/fs"
	"os"
)

//go:embed *
var tplFolder embed.FS

func GetTplFs(embedded bool) fs.FS {
	if embedded {
		return tplFolder
	}
	return os.DirFS("internal/html/template")
}

func GetErrors(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("errors.gohtml", tplFs)
}

func GetIndex(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("index.gohtml", tplFs)
}

func GetAjaxContent(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("ajax_content.gohtml", tplFs)
}
