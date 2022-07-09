package action

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/service"
	"boilerplate/internal/user/storage"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type RegisterRequest struct {
	httpin.JSONBody
	Email    string `json:"name"`
	Nickname string `json:"name"`
	Password string `json:"email"`
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
