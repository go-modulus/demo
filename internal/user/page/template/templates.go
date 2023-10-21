package template

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed *
var tplFolder embed.FS

func GetTplFolder() fs.FS {
	if os.Getenv("APP_ENV") == "dev" {
		return os.DirFS("internal/user/page/template")
	}
	return tplFolder
}

////go:embed *
//var TplFolder embed.FS
//
//var Users = template.Must(
//	template.ParseFS(TplFolder, "users.gohtml"),
//)
//
//var NewUser = template.Must(
//	template.ParseFS(TplFolder, "new_user.gohtml"),
//)
