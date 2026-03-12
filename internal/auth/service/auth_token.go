package service

import (
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/framework"
	"boilerplate/internal/graph/model"
	pgx2 "boilerplate/internal/infra/pgx"
	"boilerplate/internal/infra/utils"
	userStorage "boilerplate/internal/user/storage"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"time"
)

var CannotGenerateToken = framework.NewCommonError("CannotGenerateToken", "Cannot generate an authentication token")
var CannotSaveRefreshToken = framework.NewCommonError("CannotSaveRefreshToken", "Cannot save a refresh token")
var RefreshTokenNotFound = framework.NewCommonError("RefreshTokenNotFound", "Refresh token is not found")
var RefreshTokenExpired = framework.NewCommonError("RefreshTokenExpired", "Refresh token is expired")
var IdentityNotFound = framework.NewCommonError("IdentityNotFound", "Identity is not found")

type RefreshTokenConfig interface {
	GetRefreshTokenExpiresIn() time.Duration
	GetTokenLifetime() time.Duration
}

type AuthToken struct {
	queries     *storage.Queries
	db          *pgxpool.Pool
	tokenParser *TokenParser
	userQueries *userStorage.Queries
	logger      *zap.Logger
	rtConfig    RefreshTokenConfig
}

func NewAuthToken(
	queries *storage.Queries,
	dbConn *pgxpool.Pool,
	tokenParser *TokenParser,
	userQueries *userStorage.Queries,
	logger *zap.Logger,
	rtConfig RefreshTokenConfig,
) *AuthToken {
	return &AuthToken{
		queries:     queries,
		db:          dbConn,
		tokenParser: tokenParser,
		userQueries: userQueries,
		logger:      logger,
		rtConfig:    rtConfig,
	}
}

func (r *AuthToken) GenerateNewToken(
	ctx context.Context,
	userId uuid.UUID,
	roles []string,
	tx pgx.Tx,
) (*model.AuthPayload, error) {
	sessionId := uuid.Must(uuid.NewV6())
	result, err := r.generateNewTokenPair(ctx, sessionId, userId, roles, tx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AuthToken) RevokeToken(ctx context.Context, accessToken string) error {
	claims, err := r.tokenParser.Parse(accessToken)
	if err != nil {
		return err
	}
	expired := time.Now()
	if claims.ExpiresAt != nil {
		expired = claims.ExpiresAt.UTC()
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := r.queries.WithTx(tx)

	_, err = queries.CreateRevokedToken(
		ctx, storage.CreateRevokedTokenParams{
			TokenJti: claims.ID,
			Expired:  expired,
		},
	)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	sessionId, err := uuid.FromString(claims.SessionId)
	if err != nil {
		return err
	}

	_, err = queries.RevokeRefreshTokenBySessionId(ctx, sessionId)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthToken) generateNewTokenPair(
	ctx context.Context,
	sessionId uuid.UUID,
	accountId uuid.UUID,
	roles []string,
	tx pgx.Tx,
) (*model.AuthPayload, error) {

	refreshToken, refreshTokenExpiresAt, err := r.GenerateNewRefreshToken(ctx, sessionId, accountId, tx)
	if err != nil {
		return nil, err
	}

	token, claims, err := r.tokenParser.GenerateNewToken(
		accountId.String(),
		sessionId.String(),
		roles,
		[]string{"*"},
		r.rtConfig.GetTokenLifetime(),
	)

	if err != nil {
		r.logger.Error("Log in: cannot create a token", zap.Error(err), zap.String("user_id", accountId.String()))
		return nil, CannotGenerateToken
	}

	payload := &model.AuthPayload{
		AccessToken: &model.Token{
			Value:     token,
			ExpiresAt: int(claims.ExpiresAt.Unix()),
		},
		RefreshToken: &model.Token{
			Value:     refreshToken,
			ExpiresAt: int(refreshTokenExpiresAt.Unix()),
		},
	}

	return payload, nil
}

func (r *AuthToken) GenerateNewRefreshToken(
	ctx context.Context,
	sessionId uuid.UUID,
	accountId uuid.UUID,
	tx pgx.Tx,
) (string, *time.Time, error) {
	refreshToken := utils.RandomString(64)
	refreshTokenExpiresAt := time.Now().Add(r.rtConfig.GetRefreshTokenExpiresIn())
	queries := r.queries
	if tx != nil {
		queries = r.queries.WithTx(tx)
	}

	_, err := queries.CreateRefreshToken(
		ctx, storage.CreateRefreshTokenParams{
			Hash:      r.hashRefreshToken(refreshToken),
			UserID:    accountId,
			SessionID: sessionId,
			ExpiresAt: refreshTokenExpiresAt,
		},
	)
	if err != nil {
		r.logger.Error("Log in: cannot save a refresh token", zap.Error(err), zap.String("user_id", accountId.String()))
		return "", nil, CannotSaveRefreshToken
	}

	return refreshToken, &refreshTokenExpiresAt, nil
}
func (r *AuthToken) RefreshToken(
	ctx context.Context,
	token string,
) (*model.AuthPayload, error) {
	oldRefreshToken, err := r.queries.GetRefreshTokenByHash(ctx, r.hashRefreshToken(token))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, RefreshTokenNotFound
		} else {
			return nil, pgx2.DbIssues
		}
	}

	user, err := r.userQueries.GetUser(ctx, oldRefreshToken.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, UserIsNotFound
		} else {
			return nil, pgx2.DbIssues
		}
	}

	if time.Now().Unix() > oldRefreshToken.ExpiresAt.Unix() {
		return nil, RefreshTokenExpired
	}

	return r.GenerateNewToken(ctx, user.ID, user.Roles, nil)
}

func (r *AuthToken) hashRefreshToken(token string) string {
	return utils.HashString(token)
}
