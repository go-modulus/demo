package graphql_test

import (
	"testing"

	"github.com/go-modulus/demo/internal/blog"
	"github.com/go-modulus/demo/internal/blog/graphql"
	"github.com/go-modulus/demo/internal/blog/storage/fixture"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/test"
	"go.uber.org/fx"
)

func createMod() *module.Module {
	return blog.NewModule().
		// add the factory to the module's dependencies
		AddProviders(fixture.NewFactory)
}

var (
	// all dependencies that you want to use in tests
	resolver *graphql.Resolver
	// add a local variable of a factory to create fixtures in tests
	fixtures *fixture.Factory
)

func TestMain(m *testing.M) {
	test.LoadEnv()
	// create a new module where tested code is placed
	mod := createMod()

	test.TestMain(
		m,
		// add all necessary dependencies to the module
		module.BuildFx(mod),
		fx.Populate(
			// populate all dependencies that you want to use in tests
			&resolver,
			// populate fixtures to work with them in tests
			&fixtures,
		),
	)
}
