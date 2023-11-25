package page

import (
	"boilerplate/internal/blog/action"
	"boilerplate/internal/blog/page/template"
	"boilerplate/internal/framework"
	"boilerplate/internal/html"
	"boilerplate/internal/html/config"
)

func InitGetPostsPage(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	actionHandler *action.GetPostsAction,
	indexPage html.IndexPage,
	ajaxPage html.AjaxPage,
	config config.HtmlConfig,
) error {
	ds, err := framework.NewPageDataSource[
		*action.GetPostsRequest,
		action.PostsResponse,
	]("posts", actionHandler)

	if err != nil {
		return err
	}
	postsWidget := framework.NewWidget(
		[]*framework.TemplatePath{
			template.GetPosts(config.IsEmbeddedTemplates()),
			template.GetPost(config.IsEmbeddedTemplates()),
		},
		ds,
	)
	layout := indexPage.WithWidget(
		postsWidget,
	)
	ajaxLayout := ajaxPage.WithWidget(
		postsWidget,
	)

	if err != nil {
		return err
	}
	routes.Get("/", layout.Handler(200, nil, nil))
	routes.Get("/ajax/posts", ajaxLayout.Handler(200, nil, nil))

	return nil
}
