package http

type Pipeline struct {
	middlewares []Middleware
}

func (p *Pipeline) Push(m Middleware) {
	p.middlewares = append(
		p.middlewares,
		m,
	)
}

func (p *Pipeline) Next(next RequestHandler) RequestHandler {
	chain := Chain(p.middlewares...)

	return chain.Next(next)
}
