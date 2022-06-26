package http

import (
	"context"
	"demo/internal/errors"
	"demo/internal/validator"
	"github.com/ggicci/httpin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/fx"
	"net/http"
)

type Handler[R any] interface {
	Handle(w http.ResponseWriter, req R) error
}

type RequestHandler Handler[*http.Request]

type RequestHandlerFunc func(w http.ResponseWriter, req *http.Request) error

func (h RequestHandlerFunc) Handle(w http.ResponseWriter, req *http.Request) error {
	return h(w, req)
}

type HandlerRegistrarResult struct {
	fx.Out

	Handler HandlerRegistrar `group:"httpHandlerRegistrars"`
}

type HandlerRegistrar interface {
	Register(routes *Routes) error
}

type RegisterHandlersParams struct {
	fx.In
	Routes     *Routes
	Registrars []HandlerRegistrar `group:"httpHandlerRegistrars"`
}

func RegisterHandlers(params RegisterHandlersParams) error {
	for _, registrar := range params.Registrars {
		err := registrar.Register(params.Routes)

		if err != nil {
			return err
		}
	}

	return nil
}

func FromHttpHandler(handler http.Handler) RequestHandler {
	return RequestHandlerFunc(
		func(w http.ResponseWriter, req *http.Request) error {
			handler.ServeHTTP(w, req)

			return nil
		},
	)
}

func WrapHandler[R any](handler Handler[R]) (RequestHandler, error) {
	var req R
	engine, err := httpin.New(req)
	if err != nil {
		return nil, err
	}

	return RequestHandlerFunc(
		func(w http.ResponseWriter, httpReq *http.Request) error {
			req, err := engine.Decode(httpReq)
			if err != nil {
				return err
			}

			if validatable, ok := req.(validator.Validatable); ok {
				err := validatable.Validate(httpReq.Context())
				if err != nil {
					return err
				}
			}

			return handler.Handle(w, req.(R))
		},
	), nil
}

func UnwrapHttpinErrors(ctx context.Context, err error) error {
	var fErr *httpin.InvalidFieldError
	if !errors.As(err, &fErr) {
		return err
	}

	var vErr validation.ErrorObject
	if !errors.As(err, &vErr) {
		return err
	}

	return validator.NewValidationError(
		[]validator.FieldValidationError{
			validator.NewFieldFromOzzoErrorObject(fErr.Field, vErr),
		},
	)
}
