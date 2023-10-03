package template

import (
	"embed"
	"html/template"
)

//go:embed *
var tplFolder embed.FS

var Users = template.Must(
	template.ParseFS(tplFolder, "users.gohtml"),
)

var NewUser = template.Must(
	template.ParseFS(tplFolder, "new_user.gohtml"),
)
