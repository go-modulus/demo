package user

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func registerRouters(
	chi chi.Router,
	errorHandler *framework.HttpErrorHandler,
	registerAction *action.RegisterAction,
	getUserAction *action.GetUserAction,
	getUsersAction *action.GetUsersAction,
	updateAction *action.UpdateAction,
) error {
	err := registerAction.Register(chi, errorHandler)
	if err != nil {
		return err
	}
	err = getUserAction.Register(chi, errorHandler)
	if err != nil {
		return err
	}
	err = getUsersAction.Register(chi, errorHandler)
	if err != nil {
		return err
	}
	err = updateAction.Register(chi, errorHandler)
	if err != nil {
		return err
	}

	return nil
}

func UserPlugin() fx.Option {
	return fx.Module(
		"user",
		fx.Provide(
			action.NewRegisterAction,
			action.NewGetUserAction,
			action.NewGetUsersAction,
			action.NewUpdateAction,

			dao.NewUserFinder,
			dao.NewUserSaver,

			service.NewRegistration,
		),
		fx.Invoke(registerRouters),
	)
}
