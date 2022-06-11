package action

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/dto"
	"boilerplate/internal/user/service"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
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

func (a *RegisterAction) Register(chi chi.Router, errorHandler *framework.HttpErrorHandler) error {
	registerUser, err := framework.WrapHandler[*RegisterRequest](errorHandler, a)

	if err != nil {
		return err
	}
	chi.Post("/users", registerUser)

	return nil
}

func (a *RegisterAction) Handle(ctx context.Context, req *RegisterRequest) (*application.ActionResponse, error) {
	user := dto.User{
		Name:  req.Name,
		Email: req.Email,
	}
	result, err := a.registration.Register(ctx, user)
	if err != nil {
		return application.NewUnprocessableEntityResponse(ctx, err), nil
	}
	return application.NewSuccessCreationResponse(
		RegisterResponse{
			Id: result.Id,
		},
	), nil
}
