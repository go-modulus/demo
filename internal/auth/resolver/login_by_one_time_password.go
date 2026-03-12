package resolver

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	error2 "boilerplate/internal/auth/error"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/graph/model"
	pgx2 "boilerplate/internal/infra/pgx"
	"boilerplate/internal/infra/sendgrid"
	"boilerplate/internal/infra/utils"
	"boilerplate/internal/marketing/service"
)

const welcomeEmailTemplateId = "d-821526ad7cc54ec3a19ad761f3e2bedd"

// Login by one time password's token
// If account doesn't exist it will be created
// Errors:
// - TokenIssue - if the token is not found
// - TokenIssue - if the token is used
// - TokenIssue - if the token is expired
// - UserNotFound - if the user is not found
func (r *MutationResolver) LoginByOneTimePassword(
	ctx context.Context,
	token string,
) (*model.AuthPayload, error) {

	var (
		err             error
		oneTimePassword storage.OneTimePassword
		existIdentity   bool
	)
	oneTimePassword, err = r.authQueries.GetOneTimePasswordByToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, error2.TokenIssue
		}
		return nil, pgx2.DbIssues
	}

	if oneTimePassword.UsedAt.Valid {
		return nil, error2.TokenIssue
	}

	if time.Now().After(oneTimePassword.ExpiresAt) {
		return nil, error2.TokenIssue
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	existIdentity = true
	emailForRegistration := strings.TrimSpace(oneTimePassword.Email)
	errValidation := r.userRegistration.ValidateEmail(ctx, emailForRegistration)
	if errValidation != nil {
		return nil, errValidation
	}

	userIdentity, err := r.authQueries.SelectIdentity(ctx, emailForRegistration)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, pgx2.DbIssues
		} else {
			existIdentity = false
		}
	}

	if !existIdentity {
		userIdentityAfterReg, err := r.registrationUser(ctx, emailForRegistration, tx)
		if err != nil {
			return nil, err
		}
		userIdentity = *userIdentityAfterReg

		err = r.userRegistration.VerifyIdentity(ctx, emailForRegistration, userIdentity, tx)
		if err != nil {
			return nil, err
		}
		_ = r.sendWelcomeEmail(ctx, emailForRegistration)

		go func() {
			defer utils.RecoverPanic(r.logger)
			err := r.sendSignUpEvent(ctx, emailForRegistration, userIdentity.UserID)
			if err != nil {
				r.logger.Error("Cannot send sign up event", zap.Error(err))
			}
		}()
	}

	result, err := r.authorize(ctx, tx, oneTimePassword, userIdentity)
	if err != nil {
		return nil, err
	}

	r.tokenCookie.SetCookie(ctx, result.RefreshToken.Value)

	err = tx.Commit(ctx)
	if err != nil {
		return nil, pgx2.DbIssues
	}

	return result, nil
}

func (r *MutationResolver) registrationUser(
	ctx context.Context,
	email string,
	tx pgx.Tx,
) (*storage.Identity, error) {
	password := utils.RandomString(8)
	_, account, err := r.userRegistration.Registration(
		ctx, model.RegisterUserRequest{
			Email:    &email,
			Password: &password,
		}, tx,
	)
	if err != nil {
		return nil, err
	}

	return account.Email, nil
}

func (r *MutationResolver) sendWelcomeEmail(
	ctx context.Context,
	email string,
) error {

	host := r.sendOtpConfig.GetFrontendHost()

	emailData := sendgrid.EmailData{
		To:     email,
		ToName: "",
		TemplateVars: map[string]any{
			"viewMyAssetsLink":   host,
			"startExploringLink": host,
		},
		Subject: "🎉 Welcome Aboard! Your Authentic Journey with Digital Collectibles Begins Here 🚀",
	}
	return r.sender.SendOne(
		ctx,
		welcomeEmailTemplateId,
		emailData,
		[]string{},
	)
}

func (r *MutationResolver) sendSignUpEvent(
	ctx context.Context,
	userEmail string,
	userId uuid.UUID,
) error {
	input := r.eventSender.GetBuilder().SignUpEvent(userId, userEmail, "", service.EventChannelWebSite)

	_, err := r.eventSender.Run(ctx, input)
	return err
}

func (r *MutationResolver) authorize(
	ctx context.Context,
	tx pgx.Tx,
	oneTimePassword storage.OneTimePassword,
	userIdentity storage.Identity,
) (*model.AuthPayload, error) {
	var err error
	authQueriesWithTx := r.authQueries.WithTx(tx)

	userQueriesWithTx := r.userQueries.WithTx(tx)
	user, err := userQueriesWithTx.GetUser(ctx, userIdentity.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, pgx2.DbIssues
	}

	result, err := r.authToken.GenerateNewToken(ctx, userIdentity.UserID, user.Roles, tx)
	if err != nil {
		return nil, err
	}

	_, err = authQueriesWithTx.SetUsedOneTimePassword(
		ctx, storage.SetUsedOneTimePasswordParams{
			UserID: userIdentity.UserID,
			Token:  oneTimePassword.Token,
		},
	)
	if err != nil {
		return nil, pgx2.DbIssues
	}

	return result, nil
}
