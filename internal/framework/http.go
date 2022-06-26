package framework

import (
	"context"
	"demo/internal/errors"
	"encoding/json"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

type HttpConfig struct {
	Address string `mapstructure:"HTTP_ADDRESS"`
}

type HttpErrorHandler struct {
	baseHandler *errors.ErrorHandler
}

func NewHttpErrorHandler(baseHandler *errors.ErrorHandler) *HttpErrorHandler {
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
	Handle(ctx context.Context, req R) (*ActionResponse, error)
}

func WrapHandler[R any](
	errorHandler *HttpErrorHandler,
	handler Handler[R],
) (http.HandlerFunc, error) {
	var request R
	engine, err := httpin.New(request)
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

		err = json.NewEncoder(w).Encode(res.Response)
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
