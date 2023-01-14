package graph

import (
	"demo/graph/generated"
	"demo/graph/resolver"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"go.uber.org/fx"
)

func NewSchema(r *resolver.Resolver) graphql.ExecutableSchema {
	c := generated.Config{Resolvers: r, Directives: r.GetDirectives()}

	return generated.NewExecutableSchema(c)
}

func ConfigureServer(srv *handler.Server) {
	var mb int64 = 1 << 20

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(
		transport.MultipartForm{
			MaxUploadSize: mb * 101,
			MaxMemory:     mb * 151,
		},
	)
	srv.Use(extension.Introspection{})
}

func Module() fx.Option {
	return fx.Module(
		"graph",
		fx.Provide(
			resolver.NewResolver,
			NewSchema,
		),
		fx.Invoke(ConfigureServer),
	)
}
