package resolver_test

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/infra/test"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	"github.com/gofrs/uuid"
	"testing"
	"time"
)

func TestQueryResolver_GetOneTimePasswordEmail(t *testing.T) {
	ctx := test.GetHttpHandlerContext()

	t.Run(
		"Success", func(t *testing.T) {
			id, _ := uuid.NewV4()
			emailForOtp := "Test" + id.String() + "@gmail.com"
			oneTimePassword, rb, given := otpFixture.CreateOneTimePassword(emailForOtp)
			defer rb()

			result, err := query.GetOneTimePasswordEmail(
				ctx, oneTimePassword.Token,
			)

			spec.Given(t, given)
			spec.When(t, "Try to get email from otp token")
			spec.Then(t, "Should return email="+emailForOtp, expect.Equal(emailForOtp, result))
			spec.Then(t, "Error is empty", expect.Nil(err))

		},
	)

	t.Run(
		"Fail when token is not found", func(t *testing.T) {

			result, err := query.GetOneTimePasswordEmail(
				ctx, "ssss",
			)

			spec.When(t, "Send request with the invalid token")
			spec.Then(t, "Should return empty string", expect.Equal("", result))
			spec.HasCommonError(
				t,
				"Should return the error TokenIssue",
				err,
				error2.TokenIssue,
			)

		},
	)

	t.Run(
		"Fail when token is used", func(t *testing.T) {
			id, _ := uuid.NewV4()
			email := "test" + id.String() + "@gmail.com"
			oneTimePassword, rb, given := otpFixture.CreateOneTimePassword(email)
			defer rb()

			authPayload1, err := mutation.LoginByOneTimePassword(
				ctx, oneTimePassword.Token,
			)
			if authPayload1 != nil {
				defer rtFixture.DeleteRefreshToken(authPayload1.RefreshToken.Value)
			}

			identity := authFixture.GetIdentity(email)
			if identity != nil {
				defer userFixture.DeleteUser(identity.UserID)
			}

			result, err := query.GetOneTimePasswordEmail(
				ctx, oneTimePassword.Token,
			)

			spec.Given(t, given)
			spec.When(t, "Send request with the used token")
			spec.Then(t, "Should return empty string", expect.Equal("", result))
			spec.HasCommonError(
				t,
				"Should return the error TokenIssue",
				err,
				error2.TokenIssue,
			)
		},
	)

	t.Run(
		"Fail when token is expired", func(t *testing.T) {
			id, _ := uuid.NewV4()
			email := "test" + id.String() + "@gmail.com"
			oneTimePassword, rb, given := otpFixture.CreateOneTimePassword(email)
			defer rb()

			_, _ = otpFixture.UpdateExpiresAt(oneTimePassword.Token, time.Now().Add(time.Duration(-20)*time.Minute))

			result, err := query.GetOneTimePasswordEmail(
				ctx, oneTimePassword.Token,
			)

			spec.Given(t, given)
			spec.When(t, "Send request with the expired token")
			spec.Then(t, "Should return empty string", expect.Equal("", result))
			spec.HasCommonError(
				t,
				"Should return the error TokenIssue",
				err,
				error2.TokenIssue,
			)
		},
	)
}
