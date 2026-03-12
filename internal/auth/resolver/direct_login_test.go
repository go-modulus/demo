package resolver_test

import (
	"boilerplate/internal/auth/resolver"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/infra/test"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	"context"
	"testing"
)

func TestMutationResolver_DirectLogin(t *testing.T) {
	ctx := test.GetHttpHandlerContext()
	identity, pwd, rb, _, givenCredentials := authFixture.CreateDirectLoginCreds(
		"test",
		"test_pwd",
		storage.IdentityTypeUsername,
	)

	defer rb()

	t.Run(
		"Success", func(t *testing.T) {
			authPayload, err := mutation.DirectLogin(
				ctx, "test", "test_pwd",
			)
			defer rtFixture.DeleteRefreshToken(authPayload.RefreshToken.Value)

			w := framework.GetHttpResponseWriter(ctx)
			cookie := w.Header().Get("Set-Cookie")

			spec.Given(t, "The user with the identity 'test' and password 'test_pwd' exists")
			spec.When(t, "Login credentials are valid")
			spec.Then(t, "Should return auth payload", expect.NotNil(authPayload))
			spec.Then(t, "Error is empty", expect.Nil(err))
			spec.Then(
				t, "Refresh token is saved to cookies",
				expect.StringContains(cookie, "auth-id-rt="+authPayload.RefreshToken.Value),
				expect.NotNil(cookie),
				expect.NotEqual("", cookie),
			)
		},
	)

	t.Run(
		"Fail when identity is invalid", func(t *testing.T) {
			authPayload, err := mutation.DirectLogin(
				ctx, "test1", "test_pwd",
			)
			spec.Given(t, givenCredentials)
			spec.When(t, "Login with the invalid identity 'test1'")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error LoginOrPasswordIsWrong",
				err,
				resolver.LoginOrPasswordIsWrong,
			)
		},
	)

	t.Run(
		"Fail when password is invalid", func(t *testing.T) {
			authPayload, err := mutation.DirectLogin(
				ctx, "test", "test_pwd1",
			)
			spec.Given(t, givenCredentials)
			spec.When(t, "Login with the invalid password 'test_pwd1'")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error LoginOrPasswordIsWrong",
				err,
				resolver.LoginOrPasswordIsWrong,
			)
		},
	)

	t.Run(
		"Fail when password is not active", func(t *testing.T) {
			authQueries.ChangePasswordStatus(
				context.Background(), storage.ChangePasswordStatusParams{
					Status: storage.PasswordStatusOld,
					ID:     pwd.ID,
				},
			)
			defer authQueries.ChangePasswordStatus(
				context.Background(), storage.ChangePasswordStatusParams{
					Status: storage.PasswordStatusActive,
					ID:     pwd.ID,
				},
			)
			authPayload, err := mutation.DirectLogin(
				ctx, "test", "test_pwd",
			)
			spec.Given(t, givenCredentials, "Password 'test_pwd' is marked as old")
			spec.When(t, "Login with the deactivated password 'test_pwd'")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error LoginOrPasswordIsWrong",
				err,
				resolver.LoginOrPasswordIsWrong,
			)
		},
	)

	t.Run(
		"Fail when the old password is used", func(t *testing.T) {
			authQueries.ChangePasswordStatus(
				context.Background(), storage.ChangePasswordStatusParams{
					Status: storage.PasswordStatusOld,
					ID:     pwd.ID,
				},
			)
			defer authQueries.ChangePasswordStatus(
				context.Background(), storage.ChangePasswordStatusParams{
					Status: storage.PasswordStatusActive,
					ID:     pwd.ID,
				},
			)
			_, rb2, givenNewPassword := authFixture.CreatePasswordForUser(identity.UserID, "test_pwd1")
			defer rb2()

			authPayload, err := mutation.DirectLogin(
				ctx, "test", "test_pwd",
			)

			spec.Given(
				t,
				givenCredentials,
				givenNewPassword,
				"Password 'test_pwd' is marked as old",
			)
			spec.When(t, "Login with the deactivated password 'test_pwd'")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error OldPasswordIsUsed",
				err,
				resolver.OldPasswordIsUsed,
			)
		},
	)
}
