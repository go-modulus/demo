package auth

import (
	"boilerplate/internal/auth/provider/local"
	"errors"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	logger2 "github.com/go-pkgz/auth/logger"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"strings"
	"time"
)

type Auth struct {
	service      *auth.Service
	sessionStore *local.Session
	guard        *GuardMiddleware
}

func NewAuth(
	logger logger2.L,
	sessionStore *local.Session,
	guard *GuardMiddleware,
) *Auth {
	// define options
	options := auth.Opts{
		SecretReader: token.SecretFunc(
			func(id string) (string, error) { // secret key for JWT
				return "secret", nil
			},
		),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "my-test-app",
		URL:            "http://localhost:8181",
		AvatarStore:    avatar.NewLocalFS("/tmp"),
		Logger:         logger,
		Validator: token.ValidatorFunc(
			func(_ string, claims token.Claims) bool {
				// allow only dev_* names
				return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
			},
		),
	}

	// create auth service with providers
	service := auth.NewService(options)
	service.AddProvider(
		"google",
		"228495550698.apps.googleusercontent.com",
		"6CWB8q4p7kcWb6ECx1nzz6Ib",
	)
	service.AddDirectProvider(
		"local", provider.CredCheckerFunc(
			func(user, password string) (ok bool, err error) {
				if user == "test" && password == "test" {
					return true, nil
				}
				return false, errors.New("wrong password")
			},
		),
	)

	//m := service.Middleware()

	return &Auth{
		service:      service,
		sessionStore: sessionStore,
		guard:        guard,
	}
}

func (a *Auth) AuthGuard() *GuardMiddleware {
	return a.guard
}

func (a *Auth) RegisterAccount() error {
	return nil
}
