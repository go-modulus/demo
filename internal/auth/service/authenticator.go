package service

import (
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	userStorage "boilerplate/internal/user/storage"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"regexp"
)

var CannotSyncUserError = framework.NewCommonError(
	"cannotSyncUser",
	"Cannot sync info about the current user. Try to reauthenticate later",
)
var WrongIdError = framework.NewCommonError(
	"wrongIdFromAuthPlatform",
	"Wrong id %s gotten from the authentication platform. It should be UUID",
)
var TokenIsRevoked = framework.NewCommonError(
	"TokenIsRevoked",
	"The user cannot be authenticated because of the refresh token is revoked.",
)
var UserIsNotFound = framework.NewCommonError("UserIsNotFound", "The user is not found.")

var authRegexp = regexp.MustCompile(`(Bearer[ ]+)([^,\n$ ]+)`)

type UserFinder interface {
	GetUser(ctx context.Context, id uuid.UUID) (userStorage.User, error)
}

type Authenticator struct {
	finder      UserFinder
	tokenParser *TokenParser
	queries     *storage.Queries
}

func NewAuthenticator(
	finder UserFinder,
	tokenParser *TokenParser,
	queries *storage.Queries,
) *Authenticator {
	return &Authenticator{finder: finder, tokenParser: tokenParser, queries: queries}
}

func (a *Authenticator) Authenticate(ctx context.Context, token string) (*framework.CurrentUser, error) {
	if token == "" {
		return nil, nil
	}
	claims, err := a.tokenParser.Parse(token)
	if err != nil {
		return nil, err
	}
	_, err = a.queries.GetRevokedTokenByJti(ctx, claims.ID)
	if err == nil {
		return nil, TokenIsRevoked
	} else if err != pgx.ErrNoRows {
		return nil, err
	}

	roles := make([]string, 0, len(claims.Permissions))
	roles = append(roles, claims.Roles...)

	permissions := make([]string, 0, len(claims.Permissions))
	permissions = append(permissions, claims.Permissions...)

	id, err := uuid.FromString(claims.RegisteredClaims.Subject)
	if err != nil {
		return nil, WrongIdError.WithTplVariables(claims.ID)
	}
	_, err = a.getUser(ctx, id)
	if err != nil {
		return nil, UserIsNotFound
	}

	return &framework.CurrentUser{
		Id:          id,
		Roles:       roles,
		Permissions: permissions,
		AccessToken: token,
	}, nil
}

func (a *Authenticator) ParseAccessToken(ctx context.Context, accessToken string) string {
	submatches := authRegexp.FindStringSubmatch(accessToken)
	authToken := ""
	if len(submatches) > 2 {
		authToken = submatches[2]
	}

	return authToken
}

func (a *Authenticator) getUser(ctx context.Context, userId uuid.UUID) (*userStorage.User, error) {
	user, err := a.finder.GetUser(ctx, userId)

	if err != nil {
		return nil, CannotSyncUserError
	}
	return &user, nil
}
