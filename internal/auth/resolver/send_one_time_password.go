package resolver

import (
	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/graph/model"
	"boilerplate/internal/infra/sendgrid"
	"context"
	"fmt"
)

const oneTimePasswordEmailTemplateId = "d-800c8cb9b5c24add80c0ca16f418ca22"

type SendOneTimePasswordConfig interface {
	GetFrontendHost() string
}

// SendOneTimePassword generate token and send email with magic link for one time auth
// Errors:
// - EmailIsNotValid - if the email is not correct
// - CannotCreateOneTimePassword - if the one time password cannot be generated
// - CannotSendEmail - if the one time password cannot send
func (r *MutationResolver) SendOneTimePassword(ctx context.Context, email string) (*model.SendOneTimePasswordResult, error) {
	errValidation := r.userRegistration.ValidateEmail(ctx, email)
	if errValidation != nil {
		return nil, errValidation
	}

	existOtp, err := r.authQueries.GetLastActiveOneTimePasswordByEmail(ctx, email)
	if err == nil {
		return &model.SendOneTimePasswordResult{
			ValidTill: int(existOtp.CanResendAt.Unix()),
		}, nil
	}

	token, canResendAt, err := r.oneTimePassword.GenerateOtpCode(ctx, email)
	if err != nil {
		return nil, err
	}

	err = r.sendEmail(ctx, email, token)
	if err != nil {
		return nil, error2.CannotSendEmail
	}

	return &model.SendOneTimePasswordResult{
		ValidTill: int(canResendAt.Unix()),
	}, nil
}

// sendEmail sends an email with one time password token.
func (r *MutationResolver) sendEmail(
	ctx context.Context,
	email string,
	token string,
) error {
	host := r.sendOtpConfig.GetFrontendHost()
	magicLink := fmt.Sprintf("%s/auth?token=%s&boilerplate_source=web_site", host, token)

	emailData := sendgrid.EmailData{
		To:           email,
		ToName:       "",
		TemplateVars: map[string]any{"authLink": magicLink},
		Subject:      "Your magic link",
	}
	err := r.sender.SendOne(
		ctx,
		oneTimePasswordEmailTemplateId,
		emailData,
		[]string{},
	)
	if err != nil {
		return err
	}

	return nil
}
