.DEFAULT_GOAL := help
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

define check_console
	test -s ./bin/console ||(go build -o ./bin/console  ./cmd/console/main.go)
endef

####################################################################################################
## MAIN COMMANDS
####################################################################################################
help: ## Commands list
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'


install: ## Make a binary to ./bin folder
	go build -o ./bin/server  ./cmd/server/main.go
	go build -o ./bin/console  ./cmd/console/main.go

analyze: ## Run static analyzer
	test -s ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2
	./bin/golangci-lint run -c ./.golangci.yaml ./...

test: ## Run tests
	go test -v -failfast ./internal/...

.PHONY: mocks
mocks:
	go run github.com/vektra/mockery/v2/

####################################################################################################
## DB COMMANDS
####################################################################################################
.PHONY: migrate
migrate: ## Run migrations in both real and test databases and compiles DTOs
	$(MAKE) db-migrate
	$(MAKE) db-migrate-test
	$(MAKE) generate-db

.PHONY: check-migration
check-migration: ## Run migrations on test environment, then rollback and migrate again
	$(MAKE) db-migrate-test
	$(MAKE) db-rollback-test
	$(MAKE) db-migrate-test

.PHONY: db-add
db-add: ## Add a new migration, example: make db-add module_name migration_name
	$(check_console)
	./bin/console migrator add -m $(word 1,$(RUN_ARGS)) -n $(word 2,$(RUN_ARGS))

.PHONY: db-migrate
db-migrate: ## Run migrations in dev database
	$(check_console)
	APP_ENV=dev ./bin/console migrator migrate

.PHONY: db-migrate-test
db-migrate-test: ## Run migrations in test database
	$(check_console)
	APP_ENV=test ./bin/console migrator migrate

.PHONY: db-rollback
db-rollback: ## Rollback database migrations over the dev DB
	$(check_console)
	APP_ENV=dev ./bin/console migrator rollback

.PHONY: db-rollback-test
db-rollback-test: ## Rollback database migrations over the test DB
	$(check_console)
	APP_ENV=test ./bin/console migrator rollback

####################################################################################################
## GENERATOR COMMANDS
####################################################################################################
.PHONY: generate
generate: ## Generate public graphql schema
	go run github.com/99designs/gqlgen generate --config gqlgen.yml

.PHONY: generate-db
generate-db: ## Generate DTO and DAO for modules
	test -s ./bin/sqlc ||(cd ./bin && git clone git@github.com:debugger84/sqlc.git ./sqlc-source && cd sqlc-source && go build -o ../sqlc ./cmd/sqlc/main.go && cd .. && rm -rf ./sqlc-source)
	find . -path './internal/*/storage/sqlc.yaml' -exec ./bin/sqlc -f '{}' generate ';'


####################################################################################################
## PROFILER COMMANDS
####################################################################################################
profile-cpu:
	export PROFILER="cpu"&&./bin/dts

profile-mem:
	export PROFILER="mem"&&./bin/dts

profile-goroutine:
	export PROFILER="goroutine"&&./bin/dts

profile-block:
	export PROFILER="block"&&./bin/dts

profile-mutex:
	export PROFILER="mutex"&&./bin/dts

####################################################################################################
## VIEW PROFILER REPORTS COMMANDS
####################################################################################################
view-report-cpu:
	go tool pprof -http localhost:8111 ./profiler-reports/cpu.pprof

view-report-mem:
	go tool pprof -http localhost:8111 ./profiler-reports/mem.pprof

view-report-goroutine:
	go tool pprof -http localhost:8111 ./profiler-reports/goroutine.pprof

view-report-block:
	go tool pprof -http localhost:8111 ./profiler-reports/block.pprof

view-report-mutex:
	go tool pprof -http localhost:8111 ./profiler-reports/mutex.pprof