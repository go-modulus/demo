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
	return os.DirFS("internal/auth/widget/template")
}

func GetCurrentUser(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("current_user.gohtml", tplFs)
}
