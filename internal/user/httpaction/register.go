package httpaction

import (
	"boilerplate/internal/framework"
	validator "boilerplate/internal/ozzo-validator"
	"boilerplate/internal/user/dto"
	"boilerplate/internal/user/service"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *RegisterRequest) Validate(ctx context.Context) []framework.ValidationError {
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
	runner       *framework.ActionRunner
	registration *service.Registration
}

func NewRegisterAction(runner *framework.ActionRunner, registration *service.Registration) *RegisterAction {
	return &RegisterAction{runner: runner, registration: registration}
}

func (a *RegisterAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(
		w, r, func(ctx context.Context, request any) framework.ActionResponse {
			return a.process(ctx, request.(*RegisterRequest))
		}, &RegisterRequest{},
	)
}

func (a *RegisterAction) process(ctx context.Context, request *RegisterRequest) framework.ActionResponse {
	user := service.RegisterUserRequest{
		Name:  request.Name,
		Email: request.Email,
	}
	result, err := a.registration.Register(ctx, user)
	if err != nil {
		return framework.NewUnprocessableEntityResponse(ctx, err)
	}
	return framework.NewSuccessCreationResponse(
		RegisterResponse{
			Id: result.ID.String(),
		},
	)
}
