package page

import (
	"boilerplate/internal/blog/action"
	"boilerplate/internal/blog/page/template"
	"boilerplate/internal/blog/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/html"
	"boilerplate/internal/html/config"
	"context"
	"net/http"
)

type AddPostRequest struct {
	Title string `in:"form=title"`
	Body  string `in:"form=body"`
}
type AddPostResponse struct {
	Request       AddPostRequest
	ErrorMessages map[string]string
	Post          *storage.Post
}

type AddPostPage struct {
	addPostAction *action.AddPostAction
}

func NewAddPostPage(registerAction *action.AddPostAction) *AddPostPage {
	return &AddPostPage{addPostAction: registerAction}
}

func InitAddPostPage(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	actionHandler *AddPostPage,
	indexPage html.IndexPage,
	ajaxPage html.AjaxPage,
	config config.HtmlConfig,
) error {
	ds, err := framework.NewPageDataSource[*AddPostRequest, AddPostResponse]("newPost", actionHandler)

	if err != nil {
		return err
	}
	newPostWidget := framework.NewWidget(
		[]*framework.TemplatePath{
			template.GetNewPost(config.IsEmbeddedTemplates()),
			template.GetPost(config.IsEmbeddedTemplates()),
		},
		ds,
	)
	layout := indexPage.WithWidget(
		newPostWidget,
	)

	ajaxLayout := ajaxPage.WithWidget(
		newPostWidget,
	)

	if err != nil {
		return err
	}
	routes.Get("/blog/posts/new", layout.Handler(200, nil, nil))

	headers := http.Header{}
	errorHeaders := http.Header{}
	errorHeaders.Set("Content-Type", "text/vnd.turbo-stream.html")
	routes.Post("/ajax/blog/posts/new", ajaxLayout.Handler(201, headers, errorHeaders))
	routes.Get("/ajax/blog/posts/new", ajaxLayout.Handler(200, headers, errorHeaders))

	return nil
}

func (a *AddPostPage) Handle(ctx context.Context, req *AddPostRequest) (AddPostResponse, error) {
	httpReq := framework.GetHttpRequest(ctx)
	if httpReq.Method == "GET" {
		return AddPostResponse{}, nil
	}

	addPostReq := &action.AddPostRequest{
		Title:              req.Title,
		Body:               req.Body,
		PublishImmediately: true,
	}
	defResponse := AddPostResponse{
		Request:       *req,
		ErrorMessages: make(map[string]string),
		Post:          nil,
	}
	errors := addPostReq.Validate(ctx)
	if errors != nil {
		defResponse.ErrorMessages = errors.ErrorMessages()
		return defResponse, nil
	}
	post, err := a.addPostAction.Handle(ctx, addPostReq)
	if err != nil {
		return defResponse, err
	}
	defResponse.Post = &post

	return defResponse, nil
}
