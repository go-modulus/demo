.DEFAULT_GOAL := help
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

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
.PHONY: db-migrate
db-migrate: ## Run migrations in real database
	test -s ./bin/console ||(go build -o ./bin/console  ./cmd/console/main.go)
	./bin/console migrator migrate
	APP_ENV=test ./bin/console migrator migrate

db-rollback: ## Run migrations in real database
	test -s ./bin/console ||(go build -o ./bin/console  ./cmd/console/main.go)
	./bin/console migrator rollback
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