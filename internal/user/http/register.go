package http

import (
	"context"
	"demo/internal/http"
	"demo/internal/user/dto"
	"demo/internal/user/service"
	"demo/internal/user/storage"
	"demo/internal/validator"
	"encoding/json"
	"github.com/ggicci/httpin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	oHttp "net/http"
)

type RegisterRequest struct {
	httpin.JSONBody
	Ctx   context.Context `in:"ctx"`
	Name  string
	Email string
}

func (r *RegisterRequest) Validate(ctx context.Context) *validator.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		r,
		validation.Field(&r.Name, dto.NameRules()...),
		validation.Field(&r.Email, dto.EmailRules()...),
	)

	if err != nil {
		return validator.FromOzzoError(err)
	}

	return nil
}

type RegisterResponse struct {
	Id string
}

type RegisterAction struct {
	registration *service.Registration
}

func NewRegisterAction(registration *service.Registration) http.HandlerRegistrarResult {
	return http.HandlerRegistrarResult{Handler: &RegisterAction{registration: registration}}
}

func (a *RegisterAction) Register(routes *http.Routes) error {
	registerUser, err := http.WrapHandler[*RegisterRequest](a)

	if err != nil {
		return err
	}
	routes.Post("/users", registerUser)

	return nil
}

func (a *RegisterAction) Handle(w oHttp.ResponseWriter, req *RegisterRequest) error {
	user := storage.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	}
	result, err := a.registration.Register(req.Ctx, user)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(RegisterResponse{Id: result.ID.String()})
}
