package http

import (
	"context"
	"demo/internal/http"
	"demo/internal/user/dao"
	"demo/internal/user/dto"
	"demo/internal/user/errors"
	"demo/internal/validator"
	"encoding/json"
	"github.com/ggicci/httpin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	oHttp "net/http"
)

type UpdateRequest struct {
	httpin.JSONBody
	Ctx  context.Context `in:"ctx"`
	Id   string          `in:"path=id;required"`
	Name string
}

func (u *UpdateRequest) Validate(ctx context.Context) *validator.ValidationError {
	err := validation.ValidateStructWithContext(
		ctx,
		u,
		validation.Field(&u.Name, dto.NameRules()...),
	)

	if err != nil {
		return validator.FromOzzoError(err)
	}

	return nil
}

type UpdateResponse struct {
	Id   string
	Name string
}

type UpdateAction struct {
	finder *dao.UserFinder
	saver  *dao.UserSaver
}

func NewUpdateAction(
	finder *dao.UserFinder,
	saver *dao.UserSaver,
) http.HandlerRegistrarResult {
	return http.HandlerRegistrarResult{
		Handler: &UpdateAction{finder: finder, saver: saver},
	}
}

func (a *UpdateAction) Register(routes *http.Routes) error {
	updateUser, err := http.WrapHandler[*UpdateRequest](a)

	if err != nil {
		return err
	}
	routes.Put("/users", updateUser)

	return nil
}

func (a *UpdateAction) Handle(w oHttp.ResponseWriter, req *UpdateRequest) error {
	user, _ := a.finder.One(req.Ctx, req.Id)
	if user == nil {
		return errors.NewUserNotFound(req.Id)
	}
	user.Name = req.Name
	err := a.saver.Update(req.Ctx, *user)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(
		UpdateResponse{
			Id:   user.Id,
			Name: user.Name,
		},
	)
}
