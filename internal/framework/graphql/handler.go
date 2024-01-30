package graphql

import (
	"boilerplate/internal/framework"
	"context"
	oHttp "net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Handler struct {
	handler       *handler.Server
	authenticator framework.Authenticator
	tokenParser   framework.TokenParser
}

func NewHandler(
	handler *handler.Server,
	authenticator framework.Authenticator,
	tokenParser framework.TokenParser,
) *Handler {
	return &Handler{
		handler:       handler,
		authenticator: authenticator,
		tokenParser:   tokenParser,
	}
}

func InitHandler(
	routes *framework.Routes,
	handler *Handler,
	config *Config,
) error {
	routes.Get(config.Path, handler.Handle)
	routes.Post(config.Path, handler.Handle)

	return nil
}

func InitPlaygroundHandler(
	routes *framework.Routes,
	handler *PlaygroundHandler,
	config *Config,
) error {
	if config.Playground.Enabled {
		routes.Get(config.Playground.Path, handler.Handle)
	}

	return nil
}

func (e *Handler) Handle(w oHttp.ResponseWriter, req *oHttp.Request) {
	ctx := req.Context()

	// TODO: Maybe we should move locale and auth logic to separate middlewares?

	ctx = e.addCurrentUserToContext(ctx, req.Header.Get("Authorization"))
	ctx = framework.SetHttpRequest(ctx, req)
	ctx = framework.SetHttpResponseWriter(ctx, w)

	req = req.WithContext(ctx)

	e.handler.ServeHTTP(w, req)
}

func (e *Handler) addCurrentUserToContext(
	ctx context.Context,
	authHeader string,
) context.Context {
	var currentUser *framework.CurrentUser
	var err error

	authToken := e.tokenParser.ParseAccessToken(ctx, authHeader)
	currentUser, err = e.authenticator.Authenticate(ctx, authToken)

	if err != nil {
		return ctx
	}

	return framework.SetCurrentUser(ctx, currentUser)
}

type PlaygroundHandler struct {
	config *Config
}

func NewPlaygroundHandler(config *Config, handler *handler.Server) *PlaygroundHandler {
	return &PlaygroundHandler{config: config}
}

func (e *PlaygroundHandler) Handle(w oHttp.ResponseWriter, req *oHttp.Request) {
	playground.Handler("Graphql Playground", e.config.Path).ServeHTTP(w, req)
}
