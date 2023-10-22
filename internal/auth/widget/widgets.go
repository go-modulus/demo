package widget

import (
	"boilerplate/internal/auth/action"
	"boilerplate/internal/auth/widget/template"
	"boilerplate/internal/framework"
	"boilerplate/internal/html/config"
)

type CurrentUserWidget interface {
	framework.Widget
}

func NewCurrentUserWidget(
	currentUserAction *action.CurrentUser,
	config config.HtmlConfig,
) (CurrentUserWidget, error) {
	ds, err := framework.NewPageDataSource[*action.CurrentUserRequest, framework.CurrentUser](
		"currentUser",
		currentUserAction,
	)
	if err != nil {
		return nil, err
	}
	return framework.NewWidget(
		"current_user.gohtml",
		template.GetTplFs(config.IsEmbeddedTemplates()),
		ds,
	), nil
}
