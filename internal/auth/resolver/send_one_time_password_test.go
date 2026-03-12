package resolver_test

import (
	"boilerplate/internal/infra/test"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	"boilerplate/internal/user/errors"
	"fmt"
	"math"
	"testing"
)

func TestMutationResolver_SendOneTimePassword(t *testing.T) {
	ctx := test.GetHttpHandlerContext()

	t.Run(
		"Success", func(t *testing.T) {
			email := "qwert123@test.com"
			sendResult, err := mutation.SendOneTimePassword(
				ctx, email,
			)

			sendResultAgain, _ := mutation.SendOneTimePassword(
				ctx, email,
			)

			lastToken := otpFixture.GetLastToken()

			spec.When(t, "Sending one time password is valid")
			spec.Then(t, "Should return expires unix time", expect.NotNil(sendResult), expect.True(sendResult.ValidTill > 0))
			spec.Then(t, "Error is empty", expect.Nil(err))
			spec.Then(
				t, "Checking that token created in DB",
				expect.NotNil(lastToken),
			)
			spec.Then(
				t, "Checking that token email is correct",
				expect.Equal(email, lastToken.Email),
			)
			spec.Then(
				t, "Checking that token сanResendAt is equal with data from response",
				expect.Equal(sendResult.ValidTill, int(lastToken.CanResendAt.Unix())),
			)
			spec.Then(
				t, fmt.Sprintf("Checking that token expires is %s", otpFixture.GetTokenTtl().String()),
				expect.Equal(
					math.Round(lastToken.ExpiresAt.Sub(lastToken.CreatedAt).Seconds()),
					math.Round(otpFixture.GetTokenTtl().Seconds()),
				),
			)
			spec.Then(
				t, "Checking that after send again got previous result",
				expect.Equal(sendResult.ValidTill, sendResultAgain.ValidTill),
			)

			if lastToken != nil {
				otpFixture.DeleteByToken(lastToken.Token)
			}

		},
	)

	t.Run(
		"Fail when email is invalid", func(t *testing.T) {
			sendResult, err := mutation.SendOneTimePassword(
				ctx, "qwert123test.com",
			)
			spec.When(t, "Send request with the invalid email 'qwert123test.com'")
			spec.Then(t, "Should return nil", expect.Nil(sendResult))
			spec.HasValidationError(
				t,
				"Should return the error EmailIsNotValid",
				err,
				errors.EmailIsNotValid,
			)
		},
	)
}
