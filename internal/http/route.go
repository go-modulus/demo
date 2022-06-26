package http

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler RequestHandler
}

type Routes struct {
	routes []Route
}

func NewRoutes() *Routes {
	return &Routes{routes: make([]Route, 0)}
}

func (r *Routes) Get(path string, handler RequestHandler) {
	r.routes = append(
		r.routes,
		Route{
			Method:  http.MethodGet,
			Path:    path,
			Handler: handler,
		},
	)
}
func (r *Routes) Post(path string, handler RequestHandler) {
	r.routes = append(
		r.routes,
		Route{
			Method:  http.MethodPost,
			Path:    path,
			Handler: handler,
		},
	)
}
func (r *Routes) Put(path string, handler RequestHandler) {
	r.routes = append(
		r.routes,
		Route{
			Method:  http.MethodPut,
			Path:    path,
			Handler: handler,
		},
	)
}
func (r *Routes) Patch(path string, handler RequestHandler) {
	r.routes = append(
		r.routes,
		Route{
			Method:  http.MethodPatch,
			Path:    path,
			Handler: handler,
		},
	)
}
func (r *Routes) Delete(path string, handler RequestHandler) {
	r.routes = append(
		r.routes,
		Route{
			Method:  http.MethodDelete,
			Path:    path,
			Handler: handler,
		},
	)
}
func (r *Routes) List() []Route {
	return r.routes
}
