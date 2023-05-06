package page

import (
	"embed"
	"html/template"
)

//go:embed template
var tplFolder embed.FS

var usersTemplate = template.Must(
	template.ParseFS(tplFolder, "template/users.gohtml"),
)
