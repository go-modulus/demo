package service

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/infra/cache"
	"crypto/rsa"
	"errors"
	"github.com/gofrs/uuid"
	"time"
)

var IncorrectToken = framework.NewCommonError("incorrectToken", "An incorrect authentication token")
var TokenExpired = framework.NewCommonError("tokenIsExpired", "The authentication token is expired")

const refreshTokenFastTtl = 60

type TokenParserConfig interface {
	GetPrivateKey() string
	GetPublicKey() string
	GetCacheEnabled() bool
	GetUsersWithFastTokenRefreshing() []string
}

type TokenParser struct {
	config     TokenParserConfig
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	tokenCache *cache.Cache[string, *AuthClaims]
}

type GenerateTokenRequest struct {
	TokenJti     string
	RefreshToken string
	UserId       string
	Grants       []string
}

type AuthClaims struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	SessionId   string   `json:"sessionId"`
	jwt.RegisteredClaims
}

func newTokenCache(config TokenParserConfig) *cache.Cache[string, *AuthClaims] {
	return cache.NewCache[string, *AuthClaims](
		&cache.Config{
			MaxCachedItems: 100,
			CacheEnabled:   config.GetCacheEnabled(),
			LifeTime:       time.Minute * 5,
		},
	)
}

func NewTokenParser(
	config TokenParserConfig,
) (*TokenParser, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.GetPrivateKey()))
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.GetPublicKey()))
	if err != nil {
		return nil, err
	}

	return &TokenParser{
		config:     config,
		publicKey:  publicKey,
		privateKey: privateKey,
		tokenCache: newTokenCache(config),
	}, nil
}

// GenerateNewToken generates a new token for the user.
// The token is signed with the private key.
func (f *TokenParser) GenerateNewToken(
	userId string,
	sessionId string,
	roles []string,
	permissions []string,
	lifetime time.Duration,
) (tokenString string, claims AuthClaims, err error) {
	signingMethod := jwt.SigningMethodRS256

	id := uuid.Must(uuid.NewV6()).String()

	expAt := time.Now().Add(lifetime)
	if f.IsRefreshingTokenFast(userId) {
		expAt = time.Now().Add(time.Second * time.Duration(refreshTokenFastTtl))
	}
	claims = AuthClaims{
		Roles:       roles,
		Permissions: permissions,
		SessionId:   sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "AUTH",
			Subject:   userId,
			Audience:  jwt.ClaimStrings{"00000000-0000-0000-0000-100000000001"},
			ExpiresAt: &jwt.NumericDate{Time: expAt},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        id,
		},
	}
	token := jwt.NewWithClaims(
		signingMethod, claims,
	)

	tokenString, err = token.SignedString(f.privateKey)

	return tokenString, claims, err
}

func (f *TokenParser) IsRefreshingTokenFast(user string) bool {
	for _, authUser := range f.config.GetUsersWithFastTokenRefreshing() {
		if authUser == user {
			return true
		}
	}
	return false
}

// Parse parses the token string and returns the claims.
// If the token is invalid (expired, malformed, etc.) then the error is returned.
// Errors:
// - IncorrectToken - any error except the expired token
// - TokenExpired - the token is expired
func (f *TokenParser) Parse(tokenString string) (*AuthClaims, error) {
	claimsCache, hasItem := f.tokenCache.Get(tokenString)
	if hasItem && claimsCache != nil {
		isAlive := claimsCache.ExpiresAt.After(time.Now())
		if isAlive {
			return claimsCache, nil
		}
	}

	tokenData, err := jwt.ParseWithClaims(
		tokenString,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return f.publicKey, nil
		},
		jwt.WithAudience("00000000-0000-0000-0000-100000000001"),
		jwt.WithIssuer("AUTH"),
	)

	if err != nil || tokenData == nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, TokenExpired
		}
		return nil, IncorrectToken
	}

	if claims, ok := tokenData.Claims.(*AuthClaims); ok && tokenData.Valid {
		_ = f.tokenCache.Set(
			tokenString,
			claims,
		)
		return claims, nil
	}

	return nil, IncorrectToken
}
