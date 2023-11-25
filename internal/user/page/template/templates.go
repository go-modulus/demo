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
	return os.DirFS("internal/user/page/template")
}

func GetNewUser(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("new_user.gohtml", tplFs)
}

func GetUsers(embedded bool) *framework.TemplatePath {
	tplFs := GetTplFs(embedded)
	return framework.NewTemplatePath("users.gohtml", tplFs)
}
