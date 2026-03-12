package service_test

import (
	"boilerplate/internal/infra/test"
	"boilerplate/internal/user/storage"
	"go.uber.org/fx"
	"testing"
)

var userQuery *storage.Queries

func TestMain(m *testing.M) {

	test.TestMain(
		m, fx.Invoke(
			func(
				d2 *storage.Queries,
			) error {
				userQuery = d2
				return nil
			},
		),
	)
}
