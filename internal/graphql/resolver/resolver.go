package resolver

import (
	"github.com/go-modulus/demo/internal/auth/graphql"
	authEmailGraphql "github.com/go-modulus/demo/internal/auth/providers/email/graphql"
	authDataloader "github.com/go-modulus/demo/internal/auth/storage/dataloader"
	blogGraphql "github.com/go-modulus/demo/internal/blog/graphql"
	"github.com/go-modulus/demo/internal/graphql/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// Place all dependencies here
	blogResolver      *blogGraphql.Resolver
	authResolver      *authEmailGraphql.Resolver
	authLoaderFactory *authDataloader.LoaderFactory
}

func NewResolver(
	blogResolver *blogGraphql.Resolver,
	authResolver *authEmailGraphql.Resolver,
	authLoaderFactory *authDataloader.LoaderFactory,
) *Resolver {
	return &Resolver{
		blogResolver:      blogResolver,
		authResolver:      authResolver,
		authLoaderFactory: authLoaderFactory,
	}
}

func (r Resolver) GetDirectives() generated.DirectiveRoot {
	return generated.DirectiveRoot{
		AuthGuard: graphql.AuthGuard,
	}
}
