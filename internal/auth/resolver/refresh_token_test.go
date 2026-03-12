package resolver_test

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/auth/service"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/test"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	"github.com/gofrs/uuid"
	"testing"
)

func TestMutationResolver_RefreshToken(t *testing.T) {
	ctx := test.GetHttpHandlerContext()

	t.Run(
		"Success", func(t *testing.T) {
			_, _, rb, existingUser, _ := authFixture.CreateDirectLoginCreds("test", "test_pwd", storage.IdentityTypeEmail)

			defer rb()

			email := existingUser.VerifiedEmail.String

			oneTimePassword, otpRb, given := otpFixture.CreateOneTimePassword(email)
			defer otpRb()

			authPayload, err := mutation.LoginByOneTimePassword(
				ctx, oneTimePassword.Token,
			)
			if authPayload != nil {
				defer rtFixture.DeleteRefreshToken(authPayload.RefreshToken.Value)
			}

			authPayloadAfterRef, err := mutation.RefreshToken(ctx, &authPayload.RefreshToken.Value)
			if authPayloadAfterRef != nil {
				defer rtFixture.DeleteRefreshToken(authPayloadAfterRef.RefreshToken.Value)
			}

			accessToken := authPayloadAfterRef.AccessToken.Value
			authUser, authErr := authenticator.Authenticate(ctx, accessToken)

			spec.Given(t, given)
			spec.When(t, "Get new access token by refresh token")
			spec.Then(t, "Should return new auth payload", expect.NotNil(authPayloadAfterRef))
			spec.Then(t, "Error is empty", expect.Nil(err))
			if authPayloadAfterRef != nil {
				spec.Then(t, "New auth payload has other access token", expect.NotEqual(authPayloadAfterRef.AccessToken.Value, authPayload.AccessToken.Value))
			}

			spec.Then(t, "New auth access token allow authorized", expect.NotNil(authUser), expect.Nil(authErr))
		},
	)

	t.Run(
		"Fail when token not found", func(t *testing.T) {

			authPayload, err := mutation.RefreshToken(ctx, nil)
			token := "asdasd"
			authPayload2, err2 := mutation.RefreshToken(ctx, &token)

			spec.When(t, "Send request with token is nil")
			spec.Then(t, "Should return nil", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error notAuthenticated",
				err,
				error2.NotAuthenticated,
			)

			spec.When(t, "Send request with not exist token")
			spec.Then(t, "Should return nil", expect.Nil(authPayload2))
			spec.HasCommonError(
				t,
				"Should return the error RefreshTokenNotFound",
				err2,
				service.RefreshTokenNotFound,
			)
		},
	)

	t.Run(
		"Fail when token is expired", func(t *testing.T) {
			_, _, rb, existingUser, _ := authFixture.CreateDirectLoginCreds("test", "test_pwd", storage.IdentityTypeEmail)
			defer rb()
			token := "123456"
			refreshToken, rb2, given := rtFixture.CreateRandomRefreshToken(token, existingUser.ID)
			defer rb2()

			rtFixture.SetExpiredRefreshToken(refreshToken.Hash)
			authPayload, err := mutation.RefreshToken(ctx, &token)

			spec.Given(t, given)
			spec.When(t, "Send request with the expired token")
			spec.Then(t, "Should return nil", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error RefreshTokenExpired",
				err,
				service.RefreshTokenExpired,
			)
		},
	)

	t.Run(
		"Fail when wrong token userId ", func(t *testing.T) {
			userId, _ := uuid.NewV6()
			token := "123456"
			_, rb, given := rtFixture.CreateRandomRefreshToken(token, userId)
			defer rb()

			authPayload, err := mutation.RefreshToken(ctx, &token)

			spec.Given(t, given)
			spec.When(t, "Send request with the wrong token'userId")
			spec.Then(t, "Should return nil", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error UserIsNotFound",
				err,
				service.UserIsNotFound,
			)
		},
	)

}
