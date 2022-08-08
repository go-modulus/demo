package action

import (
	"boilerplate/internal/auth/provider/local"
	"boilerplate/internal/framework"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/ggicci/httpin"
)

type LoginRequest struct {
	httpin.JSONBody
	Identity   string `json:"identity"`
	Credential string `json:"credential"`
}
type LoginResponse struct {
	Id string `json:"id"`
}

type LoginAction struct {
	provider     *local.Provider
	sessionStore *local.Session
}

func NewLoginAction(provider *local.Provider, sessionStore *local.Session) *LoginAction {
	return &LoginAction{provider: provider, sessionStore: sessionStore}
}

func (a *LoginAction) Register(routes *framework.Routes, errorHandler *framework.HttpErrorHandler) error {
	loginUser, err := framework.WrapHandler[*LoginRequest](errorHandler, a)

	if err != nil {
		return err
	}
	routes.Post(
		"/auth/local/login",
		loginUser,
	)

	return nil
}

func (a *LoginAction) Handle(ctx context.Context, req *LoginRequest) (*application.ActionResponse, error) {
	userId, err := a.provider.Login(ctx, req.Identity, req.Credential)
	if err != nil {
		return nil, err
	}
	err = a.saveSession(ctx, userId)
	if err != nil {
		return nil, err
	}

	r := application.NewSuccessCreationResponse(
		LoginResponse{
			Id: userId,
		},
	)
	return &r, nil
}

func (a *LoginAction) saveSession(ctx context.Context, userId string) error {
	request := framework.GetHttpRequest(ctx)
	writer := framework.GetHttpResponseWriter(ctx)

	return a.sessionStore.Save(writer, request, userId)
}
