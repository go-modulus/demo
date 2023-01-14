package resolver

import (
	"demo/graph/generated"
	messenger "demo/internal/messenger/api"
)

type Resolver struct {
	messengerResolver *messenger.Resolver
}

func NewResolver(
	messengerResolver *messenger.Resolver,
) *Resolver {
	return &Resolver{
		messengerResolver: messengerResolver,
	}
}

func (r Resolver) GetDirectives() generated.DirectiveRoot {
	return generated.DirectiveRoot{}
}
