package framework

import (
	"context"
	"encoding/json"
	"errors"
	application "github.com/debugger84/modulus-application"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

type HttpConfig struct {
	Address string `mapstructure:"HTTP_ADDRESS"`
}

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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(
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

func NewChi(
	lc fx.Lifecycle,
	viper *viper.Viper,
	logger *zap.Logger,
) (chi.Router, error) {
	httpConfig := &HttpConfig{
		Address: "127.0.0.1:8080",
	}

	err := viper.Unmarshal(&httpConfig)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	httpin.UseGochiURLParam("path", chi.URLParam)

	server := http.Server{
		Addr:    httpConfig.Address,
		Handler: router,
		// TODO: Add logger
		//ErrorLog: logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info(
				"Starting http-server",
				zap.String("address", httpConfig.Address),
			)
			// TODO: Send errors to channel
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping http-server")
			return server.Shutdown(ctx)
		},
	})

	return router, nil
}

type Handler[R any] interface {
	Handle(ctx context.Context, req R) (*application.ActionResponse, error)
}

func WrapHandler[R any](
	errorHandler *HttpErrorHandler,
	handler Handler[R],
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

		val, err := engine.Decode(req)
		if err != nil {
			errorHandler.Handle(
				err,
				w,
				req,
			)
			return
		}

		res, err := handler.Handle(ctx, val.(R))

		if err != nil {
			errorHandler.Handle(err, w, req)
			return
		}

		w.WriteHeader(res.StatusCode)

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

func HttpModule() fx.Option {
	return fx.Module(
		"chi",
		fx.Provide(
			NewChi,
			NewHttpErrorHandler,
		),
	)
}
