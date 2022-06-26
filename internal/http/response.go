package http

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       any
}
