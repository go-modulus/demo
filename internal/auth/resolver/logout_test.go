package resolver_test

import (
	"boilerplate/internal/auth/service"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/infra/test"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	"testing"
)

func TestMutationResolver_Logout(t *testing.T) {
	ctx := test.GetHttpHandlerContext()
	_, _, rb, _, givenCredentials := authFixture.CreateDirectLoginCreds(
		"test",
		"test_pwd",
		storage.IdentityTypeUsername,
	)
	defer rb()
	t.Run(
		"Success", func(t *testing.T) {
			authPayload, _ := mutation.DirectLogin(
				ctx, "test", "test_pwd",
			)
			defer rtFixture.DeleteRefreshToken(authPayload.RefreshToken.Value)

			ctx = framework.SetAuthToken(ctx, authPayload.AccessToken.Value)

			result, logoutErr := mutation.Logout(ctx)

			accessToken := authPayload.AccessToken.Value
			authUser, authErr := authenticator.Authenticate(ctx, accessToken)

			claims, _ := tokenParser.Parse(accessToken)
			defer revtFixture.DeleteByTokenJti(claims.ID)

			spec.Given(t, givenCredentials)
			spec.When(t, "Logout are valid")
			spec.Then(t, "Should return true", expect.True(*result))
			spec.Then(t, "Error is empty", expect.Nil(logoutErr))
			spec.Then(t, "Can not get again authorized user ", expect.Nil(authUser))
			spec.HasCommonError(
				t,
				"Should return the error TokenIsRevoked",
				authErr,
				service.TokenIsRevoked,
			)
		},
	)

}
