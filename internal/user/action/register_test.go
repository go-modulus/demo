package action_test

import (
	"boilerplate/internal/framework"
	"boilerplate/internal/test"
	"boilerplate/internal/test/spec"
	"boilerplate/internal/user/action"
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
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

			t.Log("When try to register with valid data")
			assert.Equal(t,
				"test",
				user.Name,
				"should return user with sent data")
			assert.NoError(t, err, "should not return error")

			assert.Equal(t, "test", savedUser.Name,
				"username is saved")
			assert.Equal(t, "test@test.com", savedUser.Email,
				"user's email is saved")

			assert.Equal(t, int64(1), count,
				"local account is saved")
		},
	)

	t.Run(
		"call API", func(t *testing.T) {
			email := "test@test.com"
			routes := framework.NewRoutes()
			_ = action.InitRegisterAction(routes, errorHandler, registerAction)
			rr := test.CallPost(
				routes, "/users", map[string]interface{}{
					"name":  "test",
					"email": email,
				}, nil,
			)

			defer userFixture.DeleteUserByEmail(email)
			defer localAccountFixture.DeleteLocalAccountByEmail(email)

			t.Log(t, "When try to register with valid data")
			assert.Equal(t, http.StatusCreated, rr.Code,
				"should return status 201")
			spec.ThenJsonContains(
				t,
				"should return user with sent data",
				map[string]interface{}{
					"name":  "test",
					"email": email,
				},
				rr.Body.Bytes(),
			)
		},
	)
}
