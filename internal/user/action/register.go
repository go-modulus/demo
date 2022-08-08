package action

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/service"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/ggicci/httpin"
)

type RegisterRequest struct {
	httpin.JSONBody
	Name  string `json:"name"  validate:"required,min=3,max=50,alphaunicode"`
	Email string `json:"email"  validate:"required,email,max=150"`
}
type RegisterResponse struct {
	Id string `json:"id"`
}

type RegisterAction struct {
	registration *service.Registration
}

func NewRegisterAction(registration *service.Registration) *RegisterAction {
	return &RegisterAction{registration: registration}
}

func (a *RegisterAction) Register(routes *framework.Routes, errorHandler *framework.HttpErrorHandler) error {
	registerUser, err := framework.WrapHandler[*RegisterRequest](errorHandler, a)

	if err != nil {
		return err
	}
	routes.Post("/users", registerUser)

	return nil
}

func (a *RegisterAction) Handle(ctx context.Context, req *RegisterRequest) (*application.ActionResponse, error) {
	user := service.RegisterUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}
	result, err := a.registration.Register(ctx, user)
	if err != nil {
		r := application.NewUnprocessableEntityResponse(ctx, err)
		return &r, nil
	}
	r := application.NewSuccessCreationResponse(
		RegisterResponse{
			Id: result.ID.String(),
		},
	)
	return &r, nil
}
