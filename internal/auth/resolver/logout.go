package resolver

import (
	"boilerplate/internal/framework"
	"context"
	"go.uber.org/zap"
)

func (r *MutationResolver) Logout(ctx context.Context) (*bool, error) {
	authToken := framework.GetAuthToken(ctx)

	result := true
	err := r.authToken.RevokeToken(ctx, authToken)
	if err != nil {
		r.logger.Error("Error revoked token", zap.Error(err))
		result = false
	}

	if result {
		r.tokenCookie.RemoveCookie(ctx)
	}

	return &result, nil
}
