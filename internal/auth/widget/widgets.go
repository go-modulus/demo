package widget

import (
	"boilerplate/internal/auth/action"
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

func NewCurrentUserWidget(
	currentUserAction *action.CurrentUser,
) (CurrentUserWidget, error) {
	ds, err := framework.WrapPageDataSource[*action.CurrentUserRequest, framework.CurrentUser](nil, currentUserAction)
	if err != nil {
		return nil, err
	}
	return framework.NewWidget(
		CurrentUserTemplate,
		ds,
	), nil
}
