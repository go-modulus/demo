package widget

import (
	"boilerplate/internal/framework"
	"embed"
	"html/template"
)

//go:embed template
var tplFolder embed.FS

var CurrentUserTemplate = template.Must(
	template.ParseFS(tplFolder, "template/current_user.gohtml"),
)

var AuthTemplate = template.Must(
	template.ParseFS(tplFolder, "template/auth.gohtml"),
)

type CurrentUserWidget interface {
	framework.Widget
}
