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
	ds, err := framework.WrapPageDataSource[*action.GetUsersRequest, action.UsersResponse](errorHandler, actionHandler)

	if err != nil {
		return err
	}
	layout := indexPage.WithWidget(
		framework.NewWidget(
			template.Users,
			ds,
			[]string{
				html.LayoutBlockContent.String(),
				html.LayoutBlockTitle.String(),
			},
		),
	)
	ajaxLayout := ajaxPage.WithWidget(
		framework.NewWidget(
			template.Users,
			ds,
			[]string{
				html.LayoutBlockContent.String(),
			},
		),
	)

	if err != nil {
		return err
	}
	routes.Get("/", layout.Handler(200, nil, nil))
	routes.Get("/ajax/users", ajaxLayout.Handler(200, nil, nil))

	return nil
}