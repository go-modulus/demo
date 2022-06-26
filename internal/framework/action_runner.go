package framework

import (
	"context"
	"demo/internal/errors"
	"demo/internal/validator"
	"encoding/json"
	"github.com/ggicci/httpin"
	"net/http"
)

type ActionRunner struct {
	logger       Logger
	errorHandler *HttpErrorHandler
	router       Router
}

type ActionResponse struct {
	StatusCode      int
	Response        any
	Error           *errors.Error
	IsLoggingErrors bool
}

func NewSuccessResponse(response any) *ActionResponse {
	return &ActionResponse{
		StatusCode: 200,
		Response:   response,
	}
}

func NewSuccessCreationResponse(response any) *ActionResponse {
	return &ActionResponse{
		StatusCode: 201,
		Response:   response,
	}
}

func NewServerErrorResponse(err errors.Error) *ActionResponse {
	return &ActionResponse{
		StatusCode: 500,
		Error:      &err,
	}
}

func NewActionRunner(logger Logger, errorHandler *HttpErrorHandler, router Router) *ActionRunner {
	return &ActionRunner{logger: logger, errorHandler: errorHandler, router: router}
}

func (j *ActionRunner) Run(
	w http.ResponseWriter,
	rReq *http.Request,
	action func(ctx context.Context, request any) ActionResponse,
	req any,
) {
	ctx := rReq.Context()

	engine, err := httpin.New(req)
	if err != nil {
		j.errorHandler.Handle(err, w, rReq)
		return
	}

	dReq, err := engine.Decode(rReq)
	if err != nil {
		j.errorHandler.Handle(err, w, rReq)
		return
	}

	if validatable, ok := dReq.(validator.Validatable); ok {
		err := validatable.Validate(ctx)
		if err != nil {
			j.errorHandler.Handle(err, w, rReq)
			return
		}
	}

	res := action(ctx, dReq)
	if res.Error != nil {
		j.errorHandler.Handle(res.Error, w, rReq)
		return
	}

	w.WriteHeader(res.StatusCode)

	err = json.NewEncoder(w).Encode(res.Response)
	if err != nil {
		j.errorHandler.Handle(err, w, rReq)
		return
	}
}
