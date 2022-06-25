package framework

import (
	"net/http"
	"net/url"
)

type Router interface {
	RouteParams(r *http.Request) url.Values
}
