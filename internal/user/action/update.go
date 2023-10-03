package action

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/user/dao"
	"context"
	"github.com/ggicci/httpin"
	"go.uber.org/zap"
)

var CannotUpdateUser = framework.NewCommonError("CannotUpdateUser", "Cannot update user %s")

type UpdateRequest struct {
	httpin.JSONBody
	Id   string `in:"path=id;required"`
	Name string `json:"name"`
}
type UpdateResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UpdateAction struct {
	finder *dao.UserFinder
	saver  *dao.UserSaver
	logger *zap.Logger
}

func NewUpdateAction(
	finder *dao.UserFinder,
	saver *dao.UserSaver,
	logger *zap.Logger,
) *UpdateAction {
	return &UpdateAction{finder: finder, saver: saver, logger: logger}
}

func InitUpdateAction(
	routes *framework.Routes,
	errorHandler *framework.HttpErrorHandler,
	action *UpdateAction,
) error {
	updateUser, err := framework.WrapHandler[*UpdateRequest, UpdateResponse](errorHandler, action, 200)

	if err != nil {
		return err
	}
	routes.Put("/api/users", updateUser)

	return nil
}

func (a *UpdateAction) Handle(ctx context.Context, request *UpdateRequest) (UpdateResponse, error) {
	user, _ := a.finder.One(ctx, request.Id)
	if user == nil {
		return UpdateResponse{}, UserNotFound.WithTplVariables(request.Id)
	}
	user.Name = request.Name
	err := a.saver.Update(ctx, *user)
	if err != nil {
		a.logger.Error("error during user update", zap.Error(err))
		return UpdateResponse{}, CannotUpdateUser.WithTplVariables(request.Id)
	}
	r := UpdateResponse{
		Id:   request.Id,
		Name: user.Name,
	}
	return r, nil
}
