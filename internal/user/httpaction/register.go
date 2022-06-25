package httpaction

import (
	"boilerplate/internal/user/dto"
	"boilerplate/internal/user/service"
	"boilerplate/internal/user/storage"
	"context"
	application "github.com/debugger84/modulus-application"
	validator "github.com/debugger84/modulus-validator-ozzo"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *RegisterRequest) Validate(ctx context.Context) []application.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		r,
		validation.Field(
			&r.Name,
			dto.NameRules()...,
		),
		validation.Field(&r.Email, dto.EmailRules()...),
	)

	return validator.AsAppValidationErrors(err)
}

type RegisterResponse struct {
	Id string `json:"id"`
}

type RegisterAction struct {
	runner       *application.ActionRunner
	registration *service.Registration
}

func NewRegisterAction(runner *application.ActionRunner, registration *service.Registration) *RegisterAction {
	return &RegisterAction{runner: runner, registration: registration}
}

func (a *RegisterAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(
		w, r, func(ctx context.Context, request any) application.ActionResponse {
			return a.process(ctx, request.(*RegisterRequest))
		}, &RegisterRequest{},
	)
}

func (a *RegisterAction) process(ctx context.Context, request *RegisterRequest) application.ActionResponse {
	user := storage.CreateUserParams{
		Name:  request.Name,
		Email: request.Email,
	}
	result, err := a.registration.Register(ctx, user)
	if err != nil {
		return application.NewUnprocessableEntityResponse(ctx, err)
	}
	return application.NewSuccessCreationResponse(
		RegisterResponse{
			Id: result.ID.String(),
		},
	)
}
