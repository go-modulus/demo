package page

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/html"
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/page/template"
)

func InitGetUsersPage(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	actionHandler *action.GetUsersAction,
	indexPage html.IndexPage,
	ajaxPage html.AjaxPage,
) error {
	ds, err := framework.NewPageDataSource[*action.GetUsersRequest, action.UsersResponse]("users", actionHandler)

	if err != nil {
		return err
	}
	layout := indexPage.WithWidget(
		framework.NewWidget(
			"users.gohtml",
			template.GetTplFolder(),
			ds,
		),
	)
	ajaxLayout := ajaxPage.WithWidget(
		framework.NewWidget(
			"users.gohtml",
			template.GetTplFolder(),
			ds,
		),
	)

	if err != nil {
		return err
	}
	routes.Get("/", layout.Handler(200, nil, nil))
	routes.Get("/ajax/users", ajaxLayout.Handler(200, nil, nil))

	return nil
}
