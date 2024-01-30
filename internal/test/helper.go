package test

import (
	"boilerplate/internal"
	"boilerplate/internal/framework"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

var pool *pgxpool.Pool
var isInit = false

func TestMain(m *testing.M, modules ...fx.Option) {
	setup()

	modules = append(
		modules, fx.Invoke(
			func(innerPool *pgxpool.Pool) {
				pool = innerPool
			},
		),
	)
	_ = Invoke(modules...)

	code := m.Run()
	teardown()
	os.Exit(code)
}

func initEnv() {
	// loads values from .env into the system
	os.Setenv("APP_ENV", "test")
}

func setup() {
	if isInit {
		return
	}
	initDir()
	initEnv()

	isInit = true
}

func initDir() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func teardown() {
}

func GetDb() *pgxpool.Pool {
	return pool
}

func wrapTable(table string) string {
	tableAndSchema := strings.Split(table, ".")
	if len(tableAndSchema) != 2 {
		panic(fmt.Sprintf("invalid table name: %s", table))
	}

	return fmt.Sprintf("\"%s\".\"%s\"", tableAndSchema[0], tableAndSchema[1])
}

func buildWhere(fields map[string]any) (string, []any) {
	var wheres []string
	var values []any
	var i = 0
	for key, value := range fields {
		i++
		wheres = append(wheres, fmt.Sprintf("\"%s\" = $%d", key, i))
		values = append(values, value)
	}

	return strings.Join(wheres, " and "), values
}

func CountInDb(table string, fields map[string]interface{}) int {
	wheres, values := buildWhere(fields)
	row := pool.QueryRow(
		context.Background(),
		fmt.Sprintf(
			"select count(*) as c from %s where %s",
			wrapTable(table),
			wheres,
		),
		values...,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		panic(err)
	}

	return count
}

func HasInDb(table string, fields map[string]interface{}) bool {
	return CountInDb(table, fields) > 0
}

func HasOneInDb(table string, fields map[string]interface{}) bool {
	return CountInDb(table, fields) == 1
}

func RemoveFromDb(table string, fields map[string]any) bool {
	wheres, values := buildWhere(fields)
	sql := fmt.Sprintf(
		"delete from %s where %s",
		wrapTable(table),
		wheres,
	)

	cmd, err := pool.Exec(
		context.Background(),
		sql,
		values...,
	)
	if err != nil {
		panic(err)
	}

	return cmd.RowsAffected() > 0
}

func GetServiceFromContainer[T any]() T {
	var instance T
	_ = Invoke(
		fx.Invoke(
			func(
				d1 T,
			) error {
				instance = d1
				return nil
			},
		),
	)

	return instance
}

func Invoke(options ...fx.Option) error {
	opts := append(
		internal.Modules(),
		fx.WithLogger(
			func() fxevent.Logger {
				cfg := zap.NewDevelopmentConfig()
				cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
				cfg.DisableCaller = true
				logger, _ := cfg.Build()
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)
	opts = append(opts, options...)
	app := fx.New(
		opts...,
	)

	return app.Start(context.Background())
}

func GetServiceFromContainerWithMocks[T any](mocks []any) T {
	var instance T
	_ = Invoke(
		fx.Invoke(
			func(
				d1 T,
			) error {
				instance = d1
				return nil
			},
		),
		fx.Decorate(
			mocks...,
		),
	)

	return instance
}

func IamAuthenticatedAsTestUser(ctx context.Context) (context.Context, uuid.UUID) {
	id, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	ctx = framework.SetCurrentUserId(ctx, id.String())
	return ctx, id
}

func CallPost(
	routes *framework.Routes,
	url string,
	body map[string]any,
	header http.Header,
) *httptest.ResponseRecorder {
	bodyJson, _ := json.Marshal(body)
	reader := bytes.NewReader(bodyJson)

	req, _ := http.NewRequest("POST", url, reader)
	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	var handler http.Handler
	routesInfo := routes.GetRoutesInfo()
	for _, info := range routesInfo {
		if info.Method() == "POST" && pathMatchUrl(info.Path(), url) {
			handler = info.Handler()
		}
	}
	handler.ServeHTTP(rr, req)
	return rr
}

func pathMatchUrl(path string, url string) bool {
	pathParts := strings.Split(path, "/")
	urlParts := strings.Split(url, "/")

	if len(pathParts) != len(urlParts) {
		return false
	}
	for i, part := range pathParts {
		if part != urlParts[i] {
			if !(strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")) &&
				part != "*" {
				return false
			}
		}

	}
	return true
}
