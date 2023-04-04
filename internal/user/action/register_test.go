package action_test

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/test/expect"
	"boilerplate/internal/test/spec"
	"boilerplate/internal/user/action"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterAction_Handle(t *testing.T) {
	t.Run(
		"register successful", func(t *testing.T) {
			user, err := registerAction.Handle(
				context.Background(), &action.RegisterRequest{
					Name:  "test",
					Email: "test@test.com",
				},
			)
			defer userFixture.DeleteUser(user.Id)
			defer localAccountFixture.DeleteLocalAccount(user.Id)

			savedUser, _ := userQuery.GetUser(context.Background(), user.Id)
			count := localAccountFixture.DeleteLocalAccount(user.Id)

			spec.When(t, "try to register with valid data")
			spec.Then(
				t, "should return user with sent data",
				expect.Equal("test", user.Name),
				expect.Equal("test@test.com", user.Email),
			)
			spec.Then(t, "should not return error", expect.Nil(err))

			spec.Then(
				t, "user is saved",
				expect.Equal("test", savedUser.Name),
				expect.Equal("test@test.com", savedUser.Email),
			)
			spec.Then(t, "local account is saved", expect.Equal(int64(1), count))
		},
	)

	t.Run(
		"call API", func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/users", nil)
			routes := framework.NewRoutes()
			_ = action.InitRegisterAction(routes, errorHandler, registerAction)

			rr := httptest.NewRecorder()
			var handler http.Handler
			routesInfo := routes.GetRoutesInfo()
			for _, info := range routesInfo {
				if info.Method() == "POST" && info.Path() == "/users" {
					handler = info.Handler()
				}
			}
			handler.ServeHTTP(rr, req)
			body := rr.Body.String()
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Wrong status")
			}
			if body == "" {
				t.Errorf("Empty body")
			}
		},
	)
}
