package service

import (
	authErrors "boilerplate/internal/auth/error"
	"boilerplate/internal/framework"
	"context"
	"net/http"
	"time"
)

var sessionCookie = "auth-rt"

type TokenCookieConfig interface {
	GetRefreshTokenCookieName() string
	GetSecureCookies() bool
	GetCookiesDomain() string
	GetRefreshTokenExpiresIn() time.Duration
}
type TokenCookie struct {
	cfg TokenCookieConfig
}

func NewTokenCookie(cfg TokenCookieConfig) *TokenCookie {
	if cfg.GetRefreshTokenCookieName() != "" {
		sessionCookie = cfg.GetRefreshTokenCookieName()
	}
	return &TokenCookie{
		cfg: cfg,
	}
}

func (r *TokenCookie) GetTokenFromCookie(ctx context.Context) (string, error) {
	req := framework.GetHttpRequest(ctx)
	cookie, err := req.Cookie(sessionCookie)
	if err != nil || cookie == nil {
		return "", authErrors.NotAuthenticated
	}

	return cookie.Value, nil
}

func (r *TokenCookie) SetCookie(
	ctx context.Context,
	token string,
) {
	w := framework.GetHttpResponseWriter(ctx)

	cookies := &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now().Add(r.cfg.GetRefreshTokenExpiresIn()),
		Domain:   r.cfg.GetCookiesDomain(),
	}
	if r.cfg.GetSecureCookies() {
		cookies.Secure = true
		cookies.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookies)
}

func (r *TokenCookie) RemoveCookie(ctx context.Context) {
	w := framework.GetHttpResponseWriter(ctx)

	cookies := &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now().Add(-24 * time.Hour),
		Domain:   r.cfg.GetCookiesDomain(),
	}
	if r.cfg.GetSecureCookies() {
		cookies.Secure = true
		cookies.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookies)
}
