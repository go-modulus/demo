package action

import (
	"context"
	"demo/internal/framework"
	"demo/internal/user/dao"
	"demo/internal/user/errors"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

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

func (a *UpdateAction) Register(chi chi.Router, errorHandler *framework.HttpErrorHandler) error {
	updateUser, err := framework.WrapHandler[*UpdateRequest](errorHandler, a)

	if err != nil {
		return err
	}
	chi.Put("/users", updateUser)

	return nil
}

func (a *UpdateAction) Handle(ctx context.Context, request *UpdateRequest) (*framework.ActionResponse, error) {
	user, _ := a.finder.One(ctx, request.Id)
	if user == nil {
		return nil, errors.NewUserNotFound(request.Id)
	}
	user.Name = request.Name
	err := a.saver.Update(ctx, *user)
	if err != nil {
		return nil, err
	}

	return framework.NewSuccessResponse(
		UpdateResponse{
			Id:   request.Id,
			Name: user.Name,
		},
	), nil
}
