package action_test

import (
	fixture2 "boilerplate/internal/auth/storage/fixture"
	"boilerplate/internal/framework"
	"boilerplate/internal/test"
	"boilerplate/internal/user/action"
	"boilerplate/internal/user/storage"
	"boilerplate/internal/user/storage/fixture"
	"go.uber.org/fx"
	"testing"
)

var registerAction *action.RegisterAction
var userFixture *fixture.UserFixture
var userQuery *storage.Queries
var localAccountFixture *fixture2.LocalAccountFixture

var errorHandler *framework.HttpErrorHandler

func TestMain(m *testing.M) {
	test.TestMain(
		m, fx.Invoke(
			func(
				r *action.RegisterAction,
				f *fixture.UserFixture,
				q *storage.Queries,
				l *fixture2.LocalAccountFixture,
				eh *framework.HttpErrorHandler,
			) {
				registerAction = r
				userFixture = f
				userQuery = q
				localAccountFixture = l
				errorHandler = eh
			},
		),
	)
}
