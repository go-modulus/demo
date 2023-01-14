package graphql

import (
	"demo/internal/http"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	oHttp "net/http"
)

type Handler struct {
	config  *Config
	handler *handler.Server
}

func NewHandler(config *Config, handler *handler.Server) http.HandlerRegistrarResult {
	return http.HandlerRegistrarResult{Handler: &Handler{config: config, handler: handler}}
}

func (e *Handler) Register(routes *http.Routes) error {
	routes.Post(e.config.Path, e)

	return nil
}

func (e *Handler) Handle(w oHttp.ResponseWriter, req *oHttp.Request) error {
	e.handler.ServeHTTP(w, req)

	return nil
}

type PlaygroundHandler struct {
	config  *Config
	handler *handler.Server
}

func NewPlaygroundHandler(config *Config, handler *handler.Server) http.HandlerRegistrarResult {
	return http.HandlerRegistrarResult{Handler: &PlaygroundHandler{config: config, handler: handler}}
}

func (e *PlaygroundHandler) Register(routes *http.Routes) error {
	if e.config.Playground.Enabled {
		routes.Get(e.config.Playground.Path, e)
	}

	return nil
}

func (e *PlaygroundHandler) Handle(w oHttp.ResponseWriter, req *oHttp.Request) error {
	playground.Handler("Graphql Playground", e.config.Path).ServeHTTP(w, req)

	return nil
}
