.DEFAULT_GOAL := help
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

start=reflex -r '(\.go$|go\.mod)' -R .idea/ -s -d none $(2) -- sh -c 'make build && $(or $(value 1), /usr/bin/demo serve)'

####################################################################################################
## MAIN COMMANDS
####################################################################################################
help: ## Commands list
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

build:
	go build -buildvcs=false -o /usr/bin/demo cmd/main.go

start:
	$(call start)

graphql-generate: ## Generate public graphql schema
	go run github.com/99designs/gqlgen generate --config gqlgen.yml

m:
	migrate -database postgres://modulus:secret@postgres:5432/demo?sslmode=disable -path internal/messenger/persistence/migration $(RUN_ARGS)

install: ## Make a binary to ./bin folder
	go build -o ./bin/server  -i /cmd/server/main.go

analyze: ## Run static analyzer
	test -s ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.45.0
	./bin/golangci-lint run -c ./.golangci.yaml ./...

test: ## Run tests
	go test ./internal/...

.PHONY: mocks
mocks:
	go run github.com/vektra/mockery/v2/

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