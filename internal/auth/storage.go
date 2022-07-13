package auth

import (
	"context"
	uuid2 "github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/authboss/v3"
	aboauth "github.com/volatiletech/authboss/v3/oauth2"
	"gorm.io/gorm"
)

type GormStorage struct {
	db  *gorm.DB
	cfg *ModuleConfig
}

func NewGormStorage(db *gorm.DB, cfg *ModuleConfig) *GormStorage {
	return &GormStorage{db: db, cfg: cfg}
}

// Save the user
func (s GormStorage) Save(ctx context.Context, user authboss.User) error {
	u := user.(*User)
	result := s.db.WithContext(ctx).Table(s.cfg.UserTable).Create(u)

	return result.Error
}

func (s GormStorage) loadOauth2User(ctx context.Context, provider string, uid string) (authboss.User, error) {
	var userObj User
	result := s.db.
		WithContext(ctx).
		Table(s.cfg.UserTable).
		Where("oauth2_provider = ?", provider).
		Where("oauth2_uid = ?", uid).
		Limit(1).
		Find(&userObj)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, authboss.ErrUserNotFound
	}

	return &userObj, nil
}

func (s GormStorage) loadDirectUser(ctx context.Context, id string) (authboss.User, error) {
	var userObj User
	result := s.db.
		WithContext(ctx).
		Table(s.cfg.UserTable).
		Where("id = ?", id).
		Limit(1).
		Find(&userObj)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, authboss.ErrUserNotFound
	}

	return &userObj, nil
}

// Load the user
func (s GormStorage) Load(ctx context.Context, key string) (user authboss.User, err error) {
	// Check to see if our key is actually an oauth2 pid
	provider, uid, err := authboss.ParseOAuth2PID(key)
	if err == nil {
		return s.loadOauth2User(ctx, provider, uid)
	}

	return s.loadDirectUser(ctx, key)
}

// New user creation
func (s GormStorage) New(ctx context.Context) authboss.User {
	uuid, _ := uuid2.NewV6()
	return &User{
		ID: uuid.String(),
	}
}

// Create the user
func (s GormStorage) Create(ctx context.Context, user authboss.User) error {
	return s.Save(ctx, user)
}

// LoadByConfirmSelector looks a user up by confirmation token
func (s GormStorage) LoadByConfirmSelector(ctx context.Context, selector string) (
	authboss.ConfirmableUser,
	error,
) {
	var userObj User
	result := s.db.
		WithContext(ctx).
		Table(s.cfg.UserTable).
		Where("confirm_selector = ?", selector).
		Limit(1).
		Find(&userObj)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, authboss.ErrUserNotFound
	}

	return &userObj, nil
}

// LoadByRecoverSelector looks a user up by confirmation selector
func (s GormStorage) LoadByRecoverSelector(ctx context.Context, selector string) (
	authboss.RecoverableUser,
	error,
) {
	var userObj User
	result := s.db.
		WithContext(ctx).
		Table(s.cfg.UserTable).
		Where("recover_selector = ?", selector).
		Limit(1).
		Find(&userObj)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, authboss.ErrUserNotFound
	}

	return nil, authboss.ErrUserNotFound
}

// AddRememberToken to a user
func (s GormStorage) AddRememberToken(ctx context.Context, pid, token string) error {
	id, err := uuid2.NewV6()
	if err != nil {
		return err
	}
	tokenObj := Token{
		ID:        id.String(),
		Token:     token,
		AccountId: pid,
	}
	res := s.db.WithContext(ctx).Table(s.cfg.TokenTable).Create(tokenObj)
	return res.Error
}

// DelRememberTokens removes all tokens for the given pid
func (s GormStorage) DelRememberTokens(ctx context.Context, pid string) error {
	res := s.db.WithContext(ctx).
		Table(s.cfg.TokenTable).
		Where("account_id = ?", pid).
		Delete(&Token{})
	return res.Error
}

// UseRememberToken finds the pid-token pair and deletes it.
// If the token could not be found return ErrTokenNotFound
func (s GormStorage) UseRememberToken(ctx context.Context, pid, token string) error {
	res := s.db.WithContext(ctx).
		Table(s.cfg.TokenTable).
		Where("account_id = ?", pid).
		Where("token = ?", token).
		Delete(&Token{})
	if res.RowsAffected == 0 || res.Error != nil {
		return authboss.ErrTokenNotFound
	}

	return nil
}

// NewFromOAuth2 creates an oauth2 user (but not in the database, just a blank one to be saved later)
func (s GormStorage) NewFromOAuth2(ctx context.Context, provider string, details map[string]string) (
	authboss.OAuth2User,
	error,
) {
	switch provider {
	case "google":
		email := details[aboauth.OAuth2Email]
		var user User
		result := s.db.
			WithContext(ctx).
			Table(s.cfg.UserTable).
			Where("email = ?", email).
			Limit(1).
			Find(&user)
		if result.Error != nil {
			return nil, result.Error
		}
		user.Email = details[aboauth.OAuth2Email]
		user.OAuth2UID = details[aboauth.OAuth2UID]
		user.Confirmed = true

		return &user, nil
	}

	return nil, errors.Errorf("unknown provider %s", provider)
}

// SaveOAuth2 user
func (s GormStorage) SaveOAuth2(ctx context.Context, user authboss.OAuth2User) error {

	return s.Save(ctx, user)
}
