package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"io"
	"net/http"
	"reflect"
)

var isPathSet = false

type HttpConfig struct {
}

type RequestBody struct {
	*bytes.Reader
}

func (RequestBody) Close() error { return nil }

type HttpErrorHandler struct {
	baseHandler *ErrorHandler
}

func NewHttpErrorHandler(baseHandler *ErrorHandler) *HttpErrorHandler {
	return &HttpErrorHandler{baseHandler: baseHandler}
}

func (h HttpErrorHandler) Handle(
	err error,
	w http.ResponseWriter,
	req *http.Request,
) {
	h.baseHandler.Handle(req.Context(), err)
	if cErr, ok := err.(*CommonError); ok {
		w.WriteHeader(http.StatusInternalServerError)
		p := GetTranslator(req.Context())
		_ = json.NewEncoder(w).Encode(
			map[string]interface{}{
				"data": nil,
				"errors": []map[string]interface{}{
					{
						"message": cErr.Message(p),
						"code":    cErr.Identifier(),
					},
				},
			},
		)
	}
	if cErr, ok := err.(*ValidationErrors); ok {
		w.WriteHeader(http.StatusInternalServerError)
		p := GetTranslator(req.Context())
		errors := make([]map[string]interface{}, len(cErr.Errors()))
		for i, e := range cErr.Errors() {
			errors[i] = map[string]interface{}{
				"message": e.Message(p),
				"code":    e.Identifier,
			}
		}
		_ = json.NewEncoder(w).Encode(
			map[string]interface{}{
				"data":   nil,
				"errors": errors,
			},
		)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	_ = json.NewEncoder(w).Encode(
		map[string]interface{}{
			"data": nil,
			"errors": []map[string]interface{}{
				{
					"message": "Internal Server Error",
					"code":    "server.internalError",
				},
			},
		},
	)
}

type Handler[Req any, Resp any] interface {
	Handle(ctx context.Context, req Req) (Resp, error)
}

type EmptyHandlerRequest struct {
}

type EmptyHandlerResponse struct {
}
type EmptyHandler struct {
}

func (e EmptyHandler) Handle(ctx context.Context, req EmptyHandlerRequest) (*ActionResponse, error) {
	r := NewSuccessResponse(EmptyHandlerResponse{})
	return &r, nil
}

func WrapHandler[Req any, Resp any](
	errorHandler *HttpErrorHandler,
	handler Handler[Req, Resp],
	successCode int,
) (http.HandlerFunc, error) {
	method, ok := reflect.TypeOf(handler).MethodByName("Handle")
	if !ok {
		return nil, errors.New("invalid handler: can`t find Handle method")
	}

	methodType := method.Type

	numArgs := methodType.NumIn()
	if methodType.IsVariadic() {
		return nil, errors.New("invalid handler: variadic functions isn`t supported")
	}
	if numArgs != 3 {
		return nil, errors.New("invalid handler: invalid signature")
	}

	ctxArg := methodType.In(1)
	if ctxArg.String() != "context.Context" {
		return nil, errors.New("invalid handler: first arg must be context.Context")
	}

	reqArg := methodType.In(2)
	if reqArg.Kind() == reflect.Ptr {
		reqArg = reqArg.Elem()
	}
	if reqArg.Kind() != reflect.Struct {
		return nil, errors.New("invalid handler: second arg must be struct")
	}

	engine, err := httpin.New(reflect.New(reqArg).Interface())
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var data []byte
		if req.Body != nil {
			data, _ = io.ReadAll(req.Body)
			req.Body = RequestBody{bytes.NewReader(data)}
		}
		ctx = SetHttpRequest(ctx, req)
		ctx = SetHttpResponseWriter(ctx, w)
		req = req.WithContext(ctx)

		val, err := engine.Decode(req)
		if err != nil {
			errorHandler.Handle(
				err,
				w,
				req,
			)
			return
		}

		if validatable, ok := val.(ValidatableStruct); ok {
			validationErrs := validatable.Validate(ctx)
			if validationErrs != nil {
				errorHandler.Handle(validationErrs, w, req)
				return
			}
		}

		res, err := handler.Handle(ctx, val.(Req))

		if err != nil {
			errorHandler.Handle(err, w, req)
			return
		}
		w.WriteHeader(successCode)

		r, err := json.Marshal(res)

		if err != nil {
			errorHandler.Handle(err, w, req)
			return
		}

		_, err = w.Write(r)
		if err != nil {
			errorHandler.Handle(err, w, req)
			return
		}
	}, nil
}

type PageDataSource func(w http.ResponseWriter, req *http.Request) (any, error)
type PageErrorHandler func(w http.ResponseWriter, req *http.Request, errors []error) []error

func WrapPageDataSource[Req any, Resp any](
	errorHandler *HttpErrorHandler,
	handler Handler[Req, Resp],
) (PageDataSource, error) {
	method, ok := reflect.TypeOf(handler).MethodByName("Handle")
	if !ok {
		return nil, errors.New("invalid handler: can`t find Handle method")
	}

	methodType := method.Type

	numArgs := methodType.NumIn()
	if methodType.IsVariadic() {
		return nil, errors.New("invalid handler: variadic functions isn`t supported")
	}
	if numArgs != 3 {
		return nil, errors.New("invalid handler: invalid signature")
	}

	ctxArg := methodType.In(1)
	if ctxArg.String() != "context.Context" {
		return nil, errors.New("invalid handler: first arg must be context.Context")
	}

	reqArg := methodType.In(2)
	if reqArg.Kind() == reflect.Ptr {
		reqArg = reqArg.Elem()
	}
	if reqArg.Kind() != reflect.Struct {
		return nil, errors.New("invalid handler: second arg must be struct")
	}

	reqInterface := reflect.New(reqArg).Interface()
	engine, err := httpin.New(reqInterface)
	return func(w http.ResponseWriter, req *http.Request) (any, error) {
		if err != nil {
			return nil, err
		}
		ctx := req.Context()
		data, _ := io.ReadAll(req.Body)
		req.Body = RequestBody{bytes.NewReader(data)}

		val, err := engine.Decode(req)
		if err != nil {
			//errorHandler.Handle(
			//	err,
			//	w,
			//	req,
			//)
			return nil, err
		}

		//if validatable, ok := val.(ValidatableStruct); ok {
		//	validationErrs := validatable.Validate(ctx)
		//	if validationErrs != nil {
		//		errorHandler.Handle(validationErrs[0], w, req)
		//		return nil
		//	}
		//}

		res, err := handler.Handle(ctx, val.(Req))

		if err != nil {
			//errorHandler.Handle(err, w, req)
			return res, err
		}

		return res, nil

	}, nil
}

func HttpModule() fx.Option {
	if !isPathSet {
		httpin.UseGochiURLParam("path", chi.URLParam)
		isPathSet = true
	}

	return fx.Module(
		"chi",
		fx.Provide(
			NewHttpErrorHandler,
		),
	)
}
