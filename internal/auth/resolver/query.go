package resolver

import (
	authStorage "boilerplate/internal/auth/storage"
)

type QueryResolver struct {
	authQueries *authStorage.Queries
}

func NewQueryResolver(
	authQueries *authStorage.Queries,
) *QueryResolver {
	return &QueryResolver{
		authQueries: authQueries,
	}
}
