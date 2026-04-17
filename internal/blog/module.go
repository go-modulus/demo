package blog

import (
	"embed"

	"github.com/go-modulus/demo/internal/blog/graphql"
	"github.com/go-modulus/demo/internal/blog/storage"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/pgx"
	"github.com/go-modulus/pgx/migrator"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed storage/migration/*.sql
var migrationFS embed.FS

type ModuleConfig struct {
	// Add your module configuration here
	// e.g. Var1 string `env:"MYMODULE_VAR1, default=test"`
}

func NewModule() *module.Module {
	return module.NewModule("blog").
		// Add all dependencies of a module here
		AddDependencies(
			pgx.NewModule(),
		).
		// Add all your services here. DO NOT DELETE AddProviders call. It is used for code generation
		AddProviders(
			func(db *pgxpool.Pool) storage.DBTX {
				return db
			},
			func(db storage.DBTX) *storage.Queries {
				return storage.New(db)
			},
			migrator.ProvideMigrationFS(migrationFS),
			graphql.NewResolver,
		).
		// Add all your CLI commands here
		AddCliCommands().
		// Add all your configs here
		InitConfig(ModuleConfig{})
}
