package resolver_test

import (
	"boilerplate/internal/auth/resolver"
	service2 "boilerplate/internal/auth/service"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/auth/storage/fixture"
	framework2 "boilerplate/internal/framework"
	"boilerplate/internal/infra/test"
	fixtureUser "boilerplate/internal/user/storage/fixture"
	"go.uber.org/fx"
	"testing"
)

var mutation *resolver.MutationResolver
var query *resolver.QueryResolver
var tokenCookie *service2.TokenCookie
var authQueries *storage.Queries
var authFixture *fixture.AuthFixture
var rtFixture *fixture.RefreshTokenFixture
var revtFixture *fixture.RevokedTokenFixture
var otpFixture *fixture.OneTimePasswordFixture
var userFixture *fixtureUser.UserFixture
var contactFixture *fixtureUser.ContactFixture
var authenticator framework2.Authenticator
var tokenParser *service2.TokenParser

func TestMain(m *testing.M) {
	test.TestMain(
		m, fx.Invoke(
			func(
				r *resolver.MutationResolver,
				q *resolver.QueryResolver,
				t *service2.TokenCookie,
				a *storage.Queries,
				f *fixture.AuthFixture,
				rt *fixture.RefreshTokenFixture,
				revt *fixture.RevokedTokenFixture,
				otp *fixture.OneTimePasswordFixture,
				u *fixtureUser.UserFixture,
				uc *fixtureUser.ContactFixture,
				auth framework2.Authenticator,
				tp *service2.TokenParser,
			) {
				mutation = r
				query = q
				tokenCookie = t
				authQueries = a
				authFixture = f
				rtFixture = rt
				revtFixture = revt
				otpFixture = otp
				userFixture = u
				contactFixture = uc
				authenticator = auth
				tokenParser = tp
			},
		),
	)
}
