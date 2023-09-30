package page

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/html"
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/page/template"
	"boilerplate/internal/user/service"
	"context"
	"net/http"
)

type NewUserRequest struct {
	Name  string `in:"form=name"`
	Email string `in:"form=email"`
}
type NewUserResponse struct {
	Name          string
	Email         string
	ErrorMessages map[string]string
	IsRegistered  bool
}

type NewUserPage struct {
	registerAction *action.RegisterAction
}

func NewNewUserPage(registerAction *action.RegisterAction) *NewUserPage {
	return &NewUserPage{registerAction: registerAction}
}

func InitNewUserPage(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	actionHandler *NewUserPage,
	indexPage html.IndexPage,
	ajaxPage html.AjaxPage,
) error {
	ds, err := framework.WrapPageDataSource[*NewUserRequest, NewUserResponse](errorHandler, actionHandler)

	if err != nil {
		return err
	}
	layout := indexPage.WithWidget(
		framework.NewWidget(
			template.NewUser,
			ds,
			[]string{
				html.LayoutBlockContent.String(),
				html.LayoutBlockTitle.String(),
			},
		),
	)

	ajaxLayout := ajaxPage.WithWidget(
		framework.NewWidget(
			template.NewUser,
			ds,
			[]string{
				html.LayoutBlockContent.String(),
			},
		),
	)

	if err != nil {
		return err
	}
	routes.Get("/users/new", layout.Handler(200, nil, nil))

	headers := http.Header{}
	headers.Set("Location", "/ajax/users")
	headers.Set("Content-Type", "text/vnd.turbo-stream.html")
	errorHeaders := http.Header{}
	errorHeaders.Set("Content-Type", "text/vnd.turbo-stream.html")
	routes.Post("/ajax/users/new", ajaxLayout.Handler(201, headers, errorHeaders))

	return nil
}

func (a *NewUserPage) Handle(ctx context.Context, req *NewUserRequest) (NewUserResponse, error) {
	httpReq := framework.GetHttpRequest(ctx)
	if httpReq.Method == "GET" {
		return NewUserResponse{}, nil
	}
	defResponse := NewUserResponse{
		Name:          req.Name,
		Email:         req.Email,
		ErrorMessages: make(map[string]string),
		IsRegistered:  false,
	}
	registerReq := &action.RegisterRequest{
		Name:  req.Name,
		Email: req.Email,
	}
	errors := registerReq.Validate(ctx)
	if errors != nil {
		defResponse.ErrorMessages = errors.ErrorMessages()
		return defResponse, nil
	}
	_, err := a.registerAction.Handle(ctx, registerReq)
	if err != nil {
		if service.EmailExists.Is(err) {
			defResponse.ErrorMessages = map[string]string{
				"email": "Email already exists",
			}
			return defResponse, nil
		}
		return defResponse, err
	}
	defResponse.IsRegistered = true

	return defResponse, nil
}
