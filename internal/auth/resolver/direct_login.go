package resolver

import (
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/graph/model"
	pgx2 "boilerplate/internal/infra/pgx"
	"context"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var LoginOrPasswordIsWrong = framework.NewCommonError("LoginOrPasswordIsWrong", "Identifier or password is wrong.")
var OldPasswordIsUsed = framework.NewCommonError(
	"OldPasswordIsUsed",
	"You have used the old password. Try to recollect the new one or recover the password",
)
var UserNotFound = framework.NewCommonError("UserNotFound", "User is not found")

// DirectLogin tries to find a user with the gotten identity (phone, email or nickname).
// If the user is found, the password is checked.
// If the password is correct, the user is logged in.
// The access token is returned
// Errors:
// - LoginOrPasswordIsWrong - if the identity is not found or the password is wrong
// - OldPasswordIsUsed - if the user tries to log in with the last old password
// - CannotGenerateToken - if the access token cannot be generated
// - CannotSaveRefreshToken - if the refresh token cannot be saved into the database
// - UserNotFound - if the user is not found
func (r *MutationResolver) DirectLogin(ctx context.Context, identity string, password string) (
	*model.AuthPayload,
	error,
) {
	userIdentity, err := r.authQueries.SelectIdentity(ctx, identity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, LoginOrPasswordIsWrong
		}
		return nil, pgx2.DbIssues
	}
	passwords, err := r.authQueries.SelectUserPasswords(ctx, userIdentity.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, LoginOrPasswordIsWrong
		}
		return nil, pgx2.DbIssues
	}

	err = r.checkPassword(passwords, password)
	if err != nil {
		return nil, err
	}

	user, err := r.userQueries.GetUser(ctx, userIdentity.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, UserNotFound
		}
		return nil, pgx2.DbIssues
	}

	result, err := r.authToken.GenerateNewToken(ctx, userIdentity.UserID, user.Roles, nil)
	if err != nil {
		return nil, err
	}

	r.tokenCookie.SetCookie(ctx, result.RefreshToken.Value)
	return result, nil
}

func (r *MutationResolver) checkPassword(passwords []storage.Password, password string) error {
	if len(passwords) == 0 {
		return LoginOrPasswordIsWrong
	}
	if passwords[0].Status != storage.PasswordStatusActive {
		return LoginOrPasswordIsWrong
	}
	activePasswordError := bcrypt.CompareHashAndPassword([]byte(passwords[0].PasswordHash), []byte(password))
	if activePasswordError != nil {
		if len(passwords) > 1 {
			lastOldPasswordError := bcrypt.CompareHashAndPassword([]byte(passwords[1].PasswordHash), []byte(password))
			if lastOldPasswordError == nil {
				return OldPasswordIsUsed
			}
		}
		return LoginOrPasswordIsWrong
	}

	return nil
}
