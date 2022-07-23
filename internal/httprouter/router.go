package router

import (
	"boilerplate/internal/framework"
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"regexp"
)

var paramReg = regexp.MustCompile(`\{([a-zA-Z0-9]+)\}+`)
var wildcardReg = regexp.MustCompile(`\/\*$`)

type Router struct {
	router *httprouter.Router
	port   int
	logger framework.Logger
}

func NewRouter(config *ModuleConfig, logger framework.Logger) *Router {
	r := &httprouter.Router{
		RedirectTrailingSlash:  config.RedirectTrailingSlash,
		RedirectFixedPath:      config.RedirectFixedPath,
		HandleMethodNotAllowed: config.HandleMethodNotAllowed,
		HandleOPTIONS:          config.HandleOPTIONS,
		NotFound:               config.NotFound,
		MethodNotAllowed:       config.MethodNotAllowed,
		PanicHandler:           config.PanicHandler,
	}
	router := &Router{router: r, logger: logger}
	router.port = config.Port
	return router
}

func (r *Router) AddRoutes(routes []framework.RouteInfo) {
	for _, info := range routes {
		r.logger.Debug(context.Background(), info.Method()+": "+info.Path())
		path := info.Path()
		path = paramReg.ReplaceAllString(path, ":$1")
		path = wildcardReg.ReplaceAllString(path, "$1param1")
		r.router.Handler(info.Method(), path, info.Handler())
	}
}

func (r *Router) RouteParams(request *http.Request) url.Values {
	result := make(url.Values)
	params := httprouter.ParamsFromContext(request.Context())
	for _, param := range params {
		result.Add(param.Key, param.Value)
	}
	return result
}
