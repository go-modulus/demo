package http

import (
	"context"
	"demo/internal/errors"
	"encoding/json"
	"fmt"
	oHttp "net/http"
)

type ErrorTransformer func(ctx context.Context, err error) error

type ErrorHandler struct {
	errorHandler *errors.ErrorHandler
	transformers []ErrorTransformer
}

func NewErrorHandler(errorHandler *errors.ErrorHandler) *ErrorHandler {
	return &ErrorHandler{
		errorHandler: errorHandler,
	}
}

func (m *ErrorHandler) AttachTransformer(transformer ErrorTransformer) {
	m.transformers = append(m.transformers, transformer)
}

func (m *ErrorHandler) Handle(w oHttp.ResponseWriter, req *oHttp.Request, err error) {
	defer func() {
		if p := recover(); p != nil {
			m.errorHandler.Handle(req.Context(), errors.FromPanic(p))

			m.sendErrors(w, req, []map[string]interface{}{m.formatInternalServerError()})
		}
	}()

	ctx := req.Context()
	for _, transformer := range m.transformers {
		err = transformer(ctx, err)
	}

	uErrors := make([]map[string]interface{}, 0)

	ignoreHandling := false
	if ep, ok := err.(errors.UserErrorProvider); ok {
		uErr := ep.ToUserError()

		ignoreHandling = uErr.DontHandle

		uErrors = append(
			uErrors,
			m.formatError(uErr.Code, uErr.Message, uErr.Extra),
		)
	} else {
		if len(uErrors) == 0 {
			uErrors = append(uErrors, m.formatInternalServerError())
		}
	}

	if !ignoreHandling {
		m.errorHandler.Handle(ctx, err)
	}

	m.sendErrors(w, req, uErrors)
}

func (m *ErrorHandler) sendErrors(w oHttp.ResponseWriter, req *oHttp.Request, errors []map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(oHttp.StatusInternalServerError)

	err := json.NewEncoder(w).Encode(
		map[string]interface{}{
			"data":   nil,
			"errors": errors,
		},
	)

	if err == nil {
		return
	}

	m.errorHandler.Handle(
		req.Context(),
		fmt.Errorf("can`t send error to client: %w", err),
	)
}

func (m *ErrorHandler) formatInternalServerError() map[string]interface{} {
	return m.formatError(
		string(errors.InternalServerErrorCode),
		errors.InternalServerErrorMessage,
		map[string]interface{}{},
	)
}

func (m *ErrorHandler) formatError(code string, message string, extra map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"extra":   extra,
	}
}

func (m *ErrorHandler) Wrap(handler RequestHandler) oHttp.Handler {
	return oHttp.HandlerFunc(
		func(w oHttp.ResponseWriter, req *oHttp.Request) {
			defer func() {
				if p := recover(); p != nil {
					m.Handle(w, req, errors.FromPanic(p))
				}
			}()

			err := handler.Handle(w, req)
			if err != nil {
				m.Handle(w, req, err)
			}
		},
	)
}
