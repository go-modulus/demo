package storage_test

import (
	graphInit "boilerplate/internal/graph"
	"boilerplate/internal/pgx"
	"boilerplate/internal/user"
	"boilerplate/internal/user/dao"
	"boilerplate/internal/user/storage"
	"context"
	application "github.com/debugger84/modulus-application"
	db "github.com/debugger84/modulus-db-pg-gorm"
	graphql "github.com/debugger84/modulus-graphql"
	logger "github.com/debugger84/modulus-logger-zap"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func BenchmarkGorm(t *testing.B) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	os.Chdir(basepath + "/../../../")
	// loads values from .env into the system
	if err := godotenv.Load(basepath + "/../../../.env"); err != nil {
		log.Print("No .env file found ")
	}
	loggerConfig := logger.NewModuleConfig(nil)

	userConfig := user.NewModuleConfig()
	dbConfig := db.NewModuleConfig()
	graphQlConfig := graphql.NewModuleConfig()
	graphQlInitConfig := graphInit.NewModuleConfig()

	pgxConfig := pgx.NewModuleConfig()

	app := application.New(
		[]interface{}{
			loggerConfig,
			dbConfig,
			userConfig,
			pgxConfig,
			graphQlConfig,
			graphQlInitConfig,
		},
	)
	err := app.Run()
	if err != nil {
		log.Print("Cannot run app")
		return
	}
	containerInstance := app.Container()

	var db *dao.UserFinder
	err = containerInstance.Invoke(
		func(dbInstance *dao.UserFinder) {
			db = dbInstance
		},
	)
	if err != nil {
		log.Print("Cannot get GORM DB")
		return
	}

	for i := 0; i < t.N; i++ {
		query := db.CreateQuery(context.Background())
		query.NewerFirst()
		users := db.ListByQuery(query, 10)
		if len(users) == 0 {
			log.Print("Cannot get users")
			return
		}
	}
}

func BenchmarkSqlc(t *testing.B) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	os.Chdir(basepath + "/../../../")
	// loads values from .env into the system
	if err := godotenv.Load(basepath + "/../../../.env"); err != nil {
		log.Print("No .env file found ")
	}
	loggerConfig := logger.NewModuleConfig(nil)

	userConfig := user.NewModuleConfig()
	dbConfig := db.NewModuleConfig()
	graphQlConfig := graphql.NewModuleConfig()
	graphQlInitConfig := graphInit.NewModuleConfig()

	pgxConfig := pgx.NewModuleConfig()

	app := application.New(
		[]interface{}{
			loggerConfig,
			dbConfig,
			userConfig,
			pgxConfig,
			graphQlConfig,
			graphQlInitConfig,
		},
	)
	err := app.Run()
	if err != nil {
		log.Print("Cannot run app")
		return
	}
	containerInstance := app.Container()

	var db *storage.Queries
	err = containerInstance.Invoke(
		func(dbInstance *storage.Queries) {
			db = dbInstance
		},
	)
	if err != nil {
		log.Print("Cannot get GORM DB")
		return
	}

	for i := 0; i < t.N; i++ {
		users, _ := db.GetNewerUsers(context.Background(), 10)
		if len(users) == 0 {
			log.Print("Cannot get users")
			return
		}
	}
}
