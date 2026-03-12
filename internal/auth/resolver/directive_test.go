package resolver_test

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/auth/resolver"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/graph/model"
	"boilerplate/internal/infra/test"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	storage2 "boilerplate/internal/user/storage"
	"context"
	"github.com/99designs/gqlgen/graphql"
	"testing"
)

func TestAuthGuard_CheckRoles(t *testing.T) {
	ctx := test.GetHttpHandlerContext()

	identity, _, rb, _, givenCredentials := authFixture.CreateDirectLoginCreds("test", "test_pwd", storage.IdentityTypeUsername)
	defer rb()

	t.Run(
		"Success", func(t *testing.T) {

			currentUser := framework.CurrentUser{
				Id:    identity.UserID,
				Roles: []string{string(storage2.UserRoleUser)},
			}

			var graphqlResolver graphql.Resolver
			graphqlResolver = func(ctx context.Context) (res interface{}, err error) {
				return nil, nil
			}

			ctx = framework.SetCurrentUser(ctx, &currentUser)
			_, err := resolver.AuthGuard(ctx, nil, graphqlResolver, []model.Role{model.RoleUser})

			spec.Given(t, givenCredentials)
			spec.When(t, "Check authGuard")
			spec.Then(t, "Error is empty", expect.Nil(err))

		},
	)

	t.Run(
		"ErrorNotAllow", func(t *testing.T) {
			currentUser := framework.CurrentUser{
				Id:    identity.UserID,
				Roles: []string{string(storage2.UserRoleUser)},
			}

			var graphqlResolver graphql.Resolver
			graphqlResolver = func(ctx context.Context) (res interface{}, err error) {
				return nil, nil
			}

			ctx = framework.SetCurrentUser(ctx, &currentUser)
			_, err := resolver.AuthGuard(ctx, nil, graphqlResolver, []model.Role{model.RoleManager})

			spec.Given(t, givenCredentials)
			spec.When(t, "Check authGuard")
			spec.Then(t, "Error is not empty", expect.NotNil(err))

			spec.HasCommonError(
				t,
				"Should return the error NotAuthorized",
				err,
				error2.NotAuthorized,
			)

		},
	)
}
