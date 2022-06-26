package http

type Middleware interface {
	Next(next RequestHandler) RequestHandler
}

type MiddlewareFunc func(RequestHandler) RequestHandler

func (m MiddlewareFunc) Next(next RequestHandler) RequestHandler {
	return m(next)
}

func Chain(m ...Middleware) Middleware {
	return MiddlewareFunc(
		func(next RequestHandler) RequestHandler {
			for i := len(m) - 1; i >= 0; i-- {
				next = m[i].Next(next)
			}
			return next
		},
	)
}
