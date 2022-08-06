package local

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Session struct {
	sessionStore sessions.Store
	config       *ProviderConfig
}

func NewSession(sessionStore sessions.Store, config *ProviderConfig) *Session {
	return &Session{sessionStore: sessionStore, config: config}
}

func (s *Session) Get(
	request *http.Request,
) (string, error) {
	session, err := s.sessionStore.Get(request, s.config.SessionIdCookieName)
	if err != nil {
		return "", err
	}
	if userId, ok := session.Values["userId"]; ok {
		return userId.(string), nil
	}

	return "", nil
}

func (s *Session) Save(
	writer http.ResponseWriter,
	request *http.Request,
	userId string,
) error {
	session, _ := s.sessionStore.Get(request, s.config.SessionIdCookieName)
	session.Values["userId"] = userId
	err := session.Save(request, writer)
	if err != nil {
		return err
	}
	return nil
}
