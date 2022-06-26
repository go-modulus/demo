package http

import (
	"context"
	"demo/internal/http"
	"demo/internal/user/storage"
	"encoding/json"
	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
	oHttp "net/http"
)

type GetUserRequest struct {
	Ctx context.Context `in:"ctx"`
	Id  uuid.UUID       `in:"path=id;required"`
}

type UserResponse struct {
	Id   string
	Name string
}

type GetUserAction struct {
	finder *storage.Queries
}

func NewGetUserAction(finder *storage.Queries) http.HandlerRegistrarResult {
	return http.HandlerRegistrarResult{Handler: &GetUserAction{finder: finder}}
}

func (a *GetUserAction) Register(routes *http.Routes) error {
	getUser, err := http.WrapHandler[*GetUserRequest](a)

	if err != nil {
		return err
	}
	routes.Get("/users/{id}", getUser)

	return nil
}

func (a *GetUserAction) Handle(w oHttp.ResponseWriter, req *GetUserRequest) error {
	user, err := a.finder.GetUser(req.Ctx, guid.UUID(req.Id))
	if err != nil {
		return err
	}
	//if user == nil {
	//	return errors.NewUserNotFound(req.Id)
	//}

	var res UserResponse
	res.Id = req.Id.String()
	res.Name = user.Name

	return json.NewEncoder(w).Encode(res)
}
