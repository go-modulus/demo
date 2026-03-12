package error

import (
	"boilerplate/internal/framework"
)

var (
	NotAuthenticated = framework.NewCommonError("notAuthenticated", "Authentication is required")
	NotAuthorized    = framework.NewCommonError("notAuthorized", "You don't have permission to access this resource")
	CannotSendEmail  = framework.NewCommonError("CannotSendEmail", "Cannot send the email with magic link")

	VerificationCodeNotFound  = framework.NewCommonError("VerificationCodeNotFound", "Verification code is not found.")
	VerificationCodeIsUsed    = framework.NewCommonError("VerificationCodeIsUsed", "Verification code is used.")
	VerificationCodeIsExpired = framework.NewCommonError("VerificationCodeIsExpired", "Verification code is expired.")

	EmailIsNotValid = framework.NewCommonError("EmailIsNotValid", "Email is not valid.")

	TokenNotFound = framework.NewCommonError("TokenNotFound", "Token is not found.")

	CannotCreateOneTimePassword = framework.NewCommonError(
		"CannotCreateOneTimePassword",
		"Cannot create the one time password",
	)

	TokenIssue = framework.NewCommonError("TokenIssue", "There might be some issues with this link. Please try again.")
)
