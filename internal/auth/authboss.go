package auth

import (
	"github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"
	"os"
)

type Auth struct {
	ab *authboss.Authboss
}

func NewAuth(
	logger authboss.Logger,
	storage authboss.ServerStorer,
) *Auth {
	ab := authboss.New()

	sessKey := []byte("encription_key_32_byte_length!!!")
	cookieKey := []byte("cookieKey")

	ab.Config.Storage.Server = storage
	ab.Config.Storage.SessionState = abclientstate.NewSessionStorer("authId", sessKey, nil)
	ab.Config.Storage.CookieState = abclientstate.NewCookieStorer(cookieKey, nil)

	ab.Config.Paths.Mount = "/authboss"
	ab.Config.Paths.RootURL = "https://www.example.com/"

	ab.Core.Mailer = defaults.NewLogMailer(os.Stdout)
	ab.Core.Logger = logger

	//// This is using the renderer from: github.com/volatiletech/authboss
	//ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	//// Probably want a MailRenderer here too.
	//
	//// This instantiates and uses every default implementation
	//// in the Config.Core area that exist in the defaults package.
	//// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, true, false)

	if err := ab.Init(); err != nil {
		panic(err)
	}
	return &Auth{
		ab: ab,
	}
}

func (a *Auth) RegisterAccount() error {
	a.ab.Modules.ReStorage.Server.Save()
	return nil
}
