package action

import (
	"context"
	"demo/internal/framework"
	"demo/internal/user/service"
	"demo/internal/user/storage"
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
	user := storage.CreateUserParams{
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
