package auth

import (
	"boilerplate/internal/auth/resolver"
	"boilerplate/internal/auth/service"
	"boilerplate/internal/auth/storage"
	"boilerplate/internal/auth/storage/fixture"
	framework2 "boilerplate/internal/framework"
	userStorage "boilerplate/internal/user/storage"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	RefreshTokenName                  string        `mapstructure:"AUTH_RT_NAME"`
	SecureCookies                     bool          `mapstructure:"AUTH_SECURE_COOKIES"`
	CookiesDomain                     string        `mapstructure:"AUTH_COOKIES_DOMAIN"`
	RefreshTokenExpiresIn             time.Duration `mapstructure:"AUTH_RT_EXPIRES_IN"`
	PrivateKey                        string        `mapstructure:"AUTH_PRIVATE_KEY"`
	PublicKey                         string        `mapstructure:"AUTH_PUBLIC_KEY"`
	AccessTokenLifetime               time.Duration `mapstructure:"AUTH_AT_EXPIRES_IN"`
	CacheEnabled                      bool          `mapstructure:"AUTH_CACHE_ENABLED"`
	OneTimePasswordTtl                time.Duration `mapstructure:"ONE_TIME_PASSWORD_TTL"`
	OneTimePasswordAfterPurchasingTtl time.Duration `mapstructure:"ONE_TIME_PASSWORD_AFTER_PURCHASING_TTL"`
	OneTimePasswordResendTimeout      time.Duration `mapstructure:"ONE_TIME_PASSWORD_RESEND_TIMEOUT"`
	VerificationCodeTtl               time.Duration `mapstructure:"VERIFICATION_CODE_TTL"`
	VerificationCodeForDanaUserTtl    time.Duration `mapstructure:"VERIFICATION_CODE_FOR_DANA_USER_TTL"`
	FrontendHost                      string        `mapstructure:"FRONTEND_HOST"`
	UsersWithFastTokenRefreshing      string        `mapstructure:"USERS_WITH_FAST_TOKEN_REFRESHING"`
}

func (m *ModuleConfig) GetPrivateKey() string {
	return m.PrivateKey
}

func (m *ModuleConfig) GetPublicKey() string {
	return m.PublicKey
}

func (m *ModuleConfig) GetRefreshTokenCookieName() string {
	return m.RefreshTokenName
}

func (m *ModuleConfig) GetSecureCookies() bool {
	return m.SecureCookies
}

func (m *ModuleConfig) GetCookiesDomain() string {
	return m.CookiesDomain
}

func (m *ModuleConfig) GetRefreshTokenExpiresIn() time.Duration {
	return m.RefreshTokenExpiresIn
}

func (m *ModuleConfig) GetTokenLifetime() time.Duration {
	return m.AccessTokenLifetime
}

func (m *ModuleConfig) GetCacheEnabled() bool {
	return m.CacheEnabled
}

func (m *ModuleConfig) GetOneTimePasswordTtl() time.Duration {
	return m.OneTimePasswordTtl
}

func (m *ModuleConfig) GetOneTimePasswordAfterPurchasingTtl() time.Duration {
	return m.OneTimePasswordAfterPurchasingTtl
}

func (m *ModuleConfig) GetOneTimePasswordResendTimeout() time.Duration {
	return m.OneTimePasswordResendTimeout
}

func (m *ModuleConfig) GetVerificationCodeTtl() time.Duration {
	return m.VerificationCodeTtl
}

func (m *ModuleConfig) GetVerificationCodeForDanaUserTtl() time.Duration {
	return m.VerificationCodeForDanaUserTtl
}

func (m *ModuleConfig) GetUsersWithFastTokenRefreshing() []string {
	return strings.Split(
		strings.ReplaceAll(
			m.UsersWithFastTokenRefreshing,
			" ",
			"",
		), ",",
	)
}

func (m *ModuleConfig) GetFrontendHost() string {
	return m.FrontendHost
}

func ProvidedServices() []interface{} {
	return []interface{}{
		service.NewAuthenticator,
		resolver.NewQueryResolver,
		service.NewTokenCookie,
		service.NewTokenParser,
		service.NewAccountSaver,
		service.NewAuthToken,
		service.NewVerificationCode,
		service.NewOneTimePassword,

		resolver.NewMutationResolver,

		fixture.NewRefreshTokenFixture,
		fixture.NewAuthFixture,
		fixture.NewRevokedTokenFixture,
		fixture.NewOneTimePasswordFixture,
		fixture.NewVerificationCodeFixture,

		func(db *pgxpool.Pool) storage.DBTX {
			return db
		},
		func(db storage.DBTX) *storage.Queries {
			return storage.New(db)
		},

		func(s *userStorage.Queries) service.UserFinder { return s },
		func(s *service.OneTimePassword) service.OtpGenerator { return s },
		func(config *ModuleConfig) service.TokenCookieConfig { return config },
		func(config *ModuleConfig) service.TokenParserConfig { return config },
		func(config *ModuleConfig) service.RefreshTokenConfig { return config },
		func(config *ModuleConfig) service.VerificationCodeConfig { return config },
		func(config *ModuleConfig) service.OtpConfig { return config },
		func(config *ModuleConfig) resolver.SendOneTimePasswordConfig { return config },
	}
}

func NewModule(config ModuleConfig) fx.Option {
	return fx.Options(
		fx.Decorate(
			func(a *service.Authenticator) framework2.Authenticator {
				return a
			},
		),
		fx.Module(
			"auth",
			fx.Provide(
				append(
					ProvidedServices(),
					func(viper *viper.Viper) (*ModuleConfig, error) {
						err := viper.Unmarshal(&config)
						if err != nil {
							return nil, err
						}
						return &config, nil
					},
				)...,
			),
		),
	)
}
