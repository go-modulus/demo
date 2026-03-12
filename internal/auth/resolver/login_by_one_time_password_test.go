package resolver_test

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/infra/utils"
	"boilerplate/internal/test"
	"boilerplate/internal/test/expect"
	"boilerplate/internal/test/spec"
	"boilerplate/internal/user/errors"
	storage2 "boilerplate/internal/user/storage"
	"github.com/gofrs/uuid"
	"strings"
	"testing"
	"time"
)

func TestMutationResolver_LoginByOneTimePassword(t *testing.T) {
	ctx := test.GetHttpHandlerContext()

	t.Run(
		"Success", func(t *testing.T) {
			id, _ := uuid.NewV4()
			emailForOtp := "Test" + id.String() + "@gmail.com"
			email := strings.ToLower(emailForOtp)
			oneTimePassword, rb, given := otpFixture.CreateOneTimePassword(emailForOtp)
			defer rb()

			authPayload, err := mutation.LoginByOneTimePassword(
				ctx, oneTimePassword.Token,
			)
			if authPayload != nil {
				defer rtFixture.DeleteRefreshToken(authPayload.RefreshToken.Value)
			}

			identity := authFixture.GetIdentity(email)
			user := userFixture.GetById(identity.UserID)

			if identity != nil {
				defer userFixture.DeleteUser(identity.UserID)
			}
			contact := contactFixture.GetVerifiedContact(emailForOtp)

			oneTimePasswordAfterUsed := otpFixture.GetByToken(oneTimePassword.Token)

			w := framework.GetHttpResponseWriter(ctx)
			cookie := w.Header().Get("Set-Cookie")

			spec.Given(t, given)
			spec.When(t, "Try to login by valid token")
			spec.Then(t, "Should return auth payload", expect.NotNil(authPayload))
			spec.Then(t, "Error is empty", expect.Nil(err))
			spec.Then(
				t, "Refresh token is saved to cookies",
				expect.StringContains(cookie, "auth-id-rt="+authPayload.RefreshToken.Value),
				expect.NotNil(cookie),
				expect.NotEqual("", cookie),
			)

			spec.Then(t, "Should create identity", expect.NotNil(identity))
			spec.Then(t, "Identity is equal: "+email, expect.Equal(email, identity.Identity))
			spec.Then(t, "Identity is verified", expect.Equal(storage.IdentityStatusVerified, identity.Status))
			spec.Then(t, "Should create user", expect.NotNil(user))
			spec.Then(t, "User has verified email: "+email, expect.Equal(email, user.VerifiedEmail.String))
			spec.Then(
				t,
				"User has role: "+string(storage2.UserRoleUser),
				expect.True(utils.SliceContains[string](user.Roles, string(storage2.UserRoleUser))),
			)
			spec.Then(
				t,
				"User has contact with email: "+emailForOtp,
				expect.NotNil(contact),
				expect.Equal(user.ID, contact.UserID),
				expect.Equal(emailForOtp, contact.Value),
			)
			spec.Then(t, "Contact has verified status", expect.Equal(storage2.ContactStatusVerified, contact.Status))
			spec.Then(
				t,
				"One time password should be used by user",
				expect.True(oneTimePasswordAfterUsed.UsedAt.Valid),
				expect.Equal(user.ID, oneTimePasswordAfterUsed.UserID.UUID),
			)

		},
	)

	t.Run(
		"Success by exist account", func(t *testing.T) {
			_, _, rb, existingUser, _ := authFixture.CreateDirectLoginCreds(
				"test",
				"test_pwd",
				storage.IdentityTypeEmail,
			)

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

			identity := authFixture.GetIdentity(email)
			user := userFixture.GetById(identity.UserID)
			contact := contactFixture.GetVerifiedContact(email)

			spec.Given(t, given)
			spec.When(t, "Try to login by valid token for exist account for email: "+email)
			spec.Then(t, "Should return auth payload", expect.NotNil(authPayload))
			spec.Then(t, "Error is empty", expect.Nil(err))
			spec.Then(t, "Identity is not verified", expect.Equal(storage.IdentityStatusNotVerified, identity.Status))
			spec.Then(t, "User has verified email: "+email, expect.Equal(email, user.VerifiedEmail.String))
			spec.Then(
				t,
				"User has contact with email: "+email,
				expect.NotNil(contact),
				expect.Equal(user.ID, contact.UserID),
				expect.Equal(email, contact.Value),
			)
			spec.Then(t, "Contact has verified status", expect.Equal(storage2.ContactStatusVerified, contact.Status))

		},
	)

	t.Run(
		"Fail when email is invalid", func(t *testing.T) {
			email := "Test 333@gmail..com"
			oneTimePassword, rb, _ := otpFixture.CreateOneTimePassword(email)
			defer rb()

			authPayload, err := mutation.LoginByOneTimePassword(
				ctx, oneTimePassword.Token,
			)
			spec.When(t, "Send request with the invalid email")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
			spec.HasValidationError(
				t,
				"Should return the error EmailIsNotValid",
				err,
				errors.EmailIsNotValid,
			)
		},
	)

	t.Run(
		"Fail when token is invalid", func(t *testing.T) {
			authPayload, err := mutation.LoginByOneTimePassword(
				ctx, "test123",
			)
			spec.When(t, "Send request with the invalid token")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
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

			authPayload2, err := mutation.LoginByOneTimePassword(
				ctx, oneTimePassword.Token,
			)

			spec.Given(t, given)
			spec.When(t, "Send request with the used token")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload2))
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

			authPayload, err := mutation.LoginByOneTimePassword(
				ctx, oneTimePassword.Token,
			)
			if authPayload != nil {
				defer rtFixture.DeleteRefreshToken(authPayload.RefreshToken.Value)
			}

			spec.Given(t, given)
			spec.When(t, "Send request with the expired token")
			spec.Then(t, "Should return nil auth payload", expect.Nil(authPayload))
			spec.HasCommonError(
				t,
				"Should return the error TokenIssue",
				err,
				error2.TokenIssue,
			)
		},
	)
}
