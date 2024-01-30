package graphql

import (
	"boilerplate/internal/framework/loader"
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"sync"
)

type LoadersInitializer struct{}

func NewLoadersInitializer() *LoadersInitializer {
	return &LoadersInitializer{}
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.ResponseInterceptor
} = LoadersInitializer{}

func (LoadersInitializer) ExtensionName() string {
	return "LoadersInitializer"
}

func (LoadersInitializer) Validate(graphql.ExecutableSchema) error {
	return nil
}

// Init loaders before each operation
func (l LoadersInitializer) InterceptOperation(
	ctx context.Context,
	next graphql.OperationHandler,
) graphql.ResponseHandler {
	return next(loader.WithLoaders(ctx, &sync.Map{}))
}

// Init loaders before each event in subscription
func (l LoadersInitializer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	oc := graphql.GetOperationContext(ctx)
	operation := oc.Operation

	if operation == nil {
		return next(ctx)
	}

	if operation.Operation == ast.Subscription {
		ctx = loader.WithLoaders(ctx, &sync.Map{})
	}

	return next(ctx)
}
