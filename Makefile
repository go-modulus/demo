default: help

###
## Add these lines to your .zshrc to have autocompletion for make commands
## zstyle ':completion:*:make:*:targets' call-command true
## zstyle ':completion:*:*:make:*' tag-order 'targets'
##
####################################################################################################
## MAIN COMMANDS
####################################################################################################
.PHONY: help
help: ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: install
install: ## Make a binary to ./bin folder
	go build -o ./bin/console  ./cmd/console/main.go

.PHONY: analyze
analyze: ## Run static analyzer
	test -s ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.10.1
	./bin/golangci-lint run -c ./.golangci.yaml ./...

.PHONY: test
test: ## Run tests
	go install github.com/rakyll/gotest@latest
	gotest -v -failfast  ./internal/...

.PHONY: cover
cover: ## Run tests with coverage
	go test -v -failfast -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: mocks
mocks:
	go install github.com/vektra/mockery/v3@latest
	mockery --config .mockery.yaml

####################################################################################################
## END OF MAIN COMMANDS
####################################################################################################

####################################################################################################
## MODULE COMMANDS
####################################################################################################
.PHONY: module-install
module-install: ## install the modules from the modules manifest file https://github.com/go-modulus/modulus/blob/main/modules.json
	go install github.com/go-modulus/mtools@latest
	mtools module install

.PHONY: module-create
module-create: ## create a new module in the project
	go install github.com/go-modulus/mtools@latest
	mtools module create

.PHONY: module-add-cli
module-add-cli: ## add a new cli command to the module
	go install github.com/go-modulus/mtools@latest
	mtools module add-cli

.PHONY: module-add-json-api
module-add-json-api: ## add a new json api route to process in the module
	go install github.com/go-modulus/mtools@latest
	mtools module add-json-api

.PHONY: module-init-storage
module-init-storage: ## inits the storage feature (SQLc, migrations, queries) for an existing local module
	go install github.com/go-modulus/mtools@latest
	mtools module init-storage


####################################################################################################
## END OF MODULE COMMANDS
####################################################################################################


MAKEFILE_FOLDER := ./mk

exist := $(wildcard $(MAKEFILE_FOLDER))
ifneq ($(strip $(exist)),)
  include $(MAKEFILE_FOLDER)/*.mk
endif
