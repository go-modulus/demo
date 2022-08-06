package auth

import (
	"boilerplate/internal/auth/context"
	"boilerplate/internal/auth/provider/local"
	"boilerplate/internal/framework"
	"errors"
	"net/http"
)

type GuardMiddleware struct {
	sessionStore *local.Session
	errorHandler *framework.HttpErrorHandler
	logger       framework.Logger
	allowGuest   bool
}

func NewAuthGuardMiddleware(
	sessionStore *local.Session,
	errorHandler *framework.HttpErrorHandler,
	logger framework.Logger,
) *GuardMiddleware {
	return &GuardMiddleware{
		sessionStore: sessionStore,
		errorHandler: errorHandler,
		logger:       logger,
		allowGuest:   false,
	}
}

func (m GuardMiddleware) WithErrorHandler(errorHandler *framework.HttpErrorHandler) *GuardMiddleware {
	mNew := m
	mNew.errorHandler = errorHandler
	return &mNew
}

func (m GuardMiddleware) WithAllowGuest(allowGuest bool) *GuardMiddleware {
	mNew := m
	mNew.allowGuest = allowGuest
	return &mNew
}

func (m GuardMiddleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userId, err := m.sessionStore.Get(req)
		ctx := req.Context()
		if err != nil {
			m.errorHandler.Handle(err, w, req)
			return
		}
		if userId == "" && !m.allowGuest {
			m.errorHandler.Handle(errors.New("unauthenticated"), w, req)
			return
		}
		ctx = context.SetCurrentUserId(ctx, userId)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	}
}
