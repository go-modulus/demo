package action_test

import (
	"boilerplate/internal/blog/action"
	"boilerplate/internal/blog/storage/fixture"
	"boilerplate/internal/framework"
	"boilerplate/internal/test"
	"go.uber.org/fx"
	"testing"
)

var getPostsAction *action.GetPostsAction
var postFixture *fixture.PostFixture
var errorHandler *framework.HttpErrorHandler

func TestMain(m *testing.M) {
	test.TestMain(
		m, fx.Populate(
			&getPostsAction,
			&postFixture,
			&errorHandler,
		),
	)
}
