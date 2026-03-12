package resolver

import (
	error2 "boilerplate/internal/auth/error"
	pgx2 "boilerplate/internal/infra/pgx"
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

// Get One Time Password Email
// Errors:
// - TokenIssue - if the token is not found
// - TokenIssue - if the token is used
// - TokenIssue - if the token is expired
// - UserNotFound - if the user is not found
func (r *QueryResolver) GetOneTimePasswordEmail(ctx context.Context, token string) (string, error) {

	oneTimePassword, err := r.authQueries.GetOneTimePasswordByToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", error2.TokenIssue
		}
		return "", pgx2.DbIssues
	}

	if oneTimePassword.UsedAt.Valid {
		return "", error2.TokenIssue
	}

	if time.Now().After(oneTimePassword.ExpiresAt) {
		return "", error2.TokenIssue
	}

	return oneTimePassword.Email, nil
}
