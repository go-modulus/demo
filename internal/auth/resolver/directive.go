package resolver

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/framework"
	"boilerplate/internal/graph/model"
	"boilerplate/internal/infra/utils"
	"context"
	"github.com/99designs/gqlgen/graphql"
	"strings"
)

func AuthGuard(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (
	res interface{},
	err error,
) {
	user := framework.GetCurrentUser(ctx)
	if user == nil {
		return nil, error2.NotAuthenticated
	}
	allowRoleList := make([]string, 0, len(roles))
	for _, role := range roles {
		allowRoleList = append(allowRoleList, strings.ToLower(role.String()))
	}

	hasRole := false
	for _, userRole := range user.Roles {
		if utils.SliceContains[string](allowRoleList, userRole) {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return nil, error2.NotAuthorized
	}

	return next(ctx)
}
