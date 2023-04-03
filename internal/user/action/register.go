package action

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/service"
	"context"
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

func InitRegisterAction(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	action *RegisterAction,
) error {
	registerUser, err := framework.WrapHandler[*RegisterRequest, RegisterResponse](errorHandler, action, 201)

	if err != nil {
		return err
	}
	routes.Post("/users", registerUser)

	return nil
}

func (a *RegisterAction) Handle(ctx context.Context, req *RegisterRequest) (RegisterResponse, error) {
	user := service.RegisterUserRequest{
		Name:  req.Name,
		Email: req.Email,
	}
	result, err := a.registration.Register(ctx, user)
	if err != nil {
		return RegisterResponse{}, err
	}
	r := RegisterResponse{
		Id: result.ID.String(),
	}
	return r, nil
}
