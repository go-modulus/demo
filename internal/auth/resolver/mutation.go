package resolver

import (
	"boilerplate/internal/auth/service"
	authStorage "boilerplate/internal/auth/storage"
	"boilerplate/internal/infra/sendgrid"
	service2 "boilerplate/internal/marketing/service"
	userService "boilerplate/internal/user/service"
	userStorage "boilerplate/internal/user/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type MutationResolver struct {
	db               *pgxpool.Pool
	userQueries      *userStorage.Queries
	authQueries      *authStorage.Queries
	authToken        *service.AuthToken
	logger           *zap.Logger
	sendOtpConfig    SendOneTimePasswordConfig
	tokenCookie      *service.TokenCookie
	oneTimePassword  *service.OneTimePassword
	sender           *sendgrid.Sender
	userRegistration *userService.UserRegistration
	eventSender      service2.MarketplaceEventSender
}

func NewMutationResolver(
	db *pgxpool.Pool,
	userQueries *userStorage.Queries,
	authQueries *authStorage.Queries,
	authToken *service.AuthToken,
	logger *zap.Logger,
	sendOtpConfig SendOneTimePasswordConfig,
	tokenCookie *service.TokenCookie,
	oneTimePassword *service.OneTimePassword,
	sender *sendgrid.Sender,
	userRegistration *userService.UserRegistration,
	eventSender service2.MarketplaceEventSender,
) *MutationResolver {
	return &MutationResolver{
		db:               db,
		userQueries:      userQueries,
		authQueries:      authQueries,
		authToken:        authToken,
		logger:           logger,
		sendOtpConfig:    sendOtpConfig,
		tokenCookie:      tokenCookie,
		oneTimePassword:  oneTimePassword,
		sender:           sender,
		userRegistration: userRegistration,
		eventSender:      eventSender,
	}
}
