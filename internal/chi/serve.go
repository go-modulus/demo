package chi

import (
	"context"
	"demo/internal/http"
	"demo/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	oHttp "net/http"
)

type Serve struct {
	lc           fx.Lifecycle
	config       *Config
	router       chi.Router
	routes       *http.Routes
	errorHandler *http.ErrorHandler
	errChannel   chan<- error
	logger       logger.Logger
}

type ServeParams struct {
	fx.In

	Lc           fx.Lifecycle
	Config       *Config
	Router       chi.Router
	Routes       *http.Routes
	ErrorHandler *http.ErrorHandler
	ErrChannel   chan<- error `name:"errors.channel"`
	Logger       logger.Logger
}

func NewServe(params ServeParams) *Serve {
	return &Serve{
		lc:           params.Lc,
		config:       params.Config,
		router:       params.Router,
		routes:       params.Routes,
		errorHandler: params.ErrorHandler,
		errChannel:   params.ErrChannel,
		logger:       params.Logger,
	}
}

func (s *Serve) Command() *cobra.Command {
	return &cobra.Command{
		Use:  "serve",
		Args: cobra.NoArgs,
		Run:  s.Run,
	}
}

func (s *Serve) Run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	server := &oHttp.Server{
		Addr:    s.config.Address,
		Handler: s.router,
		// TODO: Add logger
		//ErrorLog: logger,
	}

	s.lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			s.logger.Info(ctx, "Stopping http-server")
			return server.Shutdown(ctx)
		},
	})

	for _, route := range s.routes.List() {
		handler := route.Handler

		s.router.Method(
			route.Method,
			route.Path,
			s.errorHandler.Wrap(handler),
		)
	}

	s.logger.Info(
		ctx,
		"Starting http-server",
		logger.Field("address", s.config.Address),
	)

	go func() {
		err := server.ListenAndServe()
		s.errChannel <- err
	}()
}
