package resolver

import (
	"boilerplate/internal/graph/model"
	"context"
)

// Refresh access token
// Errors:
// - NotAuthenticated - if the refreshToken is nil and empty cookies
// - RefreshTokenNotFound - if the refreshToken is not found
// - UserIsNotFound - if the refreshToken's user is not found
// - RefreshTokenExpired - if the refreshToken is expired
func (r *MutationResolver) RefreshToken(ctx context.Context, refreshToken *string) (*model.AuthPayload, error) {
	token, err := r.tokenCookie.GetTokenFromCookie(ctx)
	if err != nil {
		if refreshToken == nil {
			return nil, err
		} else {
			token = *refreshToken
		}
	}

	return r.authToken.RefreshToken(ctx, token)
}
