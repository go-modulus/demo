####################################################################################################
## DB COMMANDS
####################################################################################################

.PHONY: db
db: ## Run all db commands
	go install github.com/go-modulus/mtools/cmd/mtools@latest
	$(MAKE) db-sqlc-generate
	$(MAKE) db-migrate

.PHONY: db-add
db-add: ## Create new migration in the selected module
	mtools db add

.PHONY: db-migrate
db-migrate: ## Run migrations from all modules
	mtools db migrate

.PHONY: db-rollback
db-rollback: ## Rollback the last database migration over the current DB
	mtools db rollback

.PHONY: db-check-migration
db-check-migration: ## Run migrations on test environment, then rollback and migrate again
	$(MAKE) db-migrate
	$(MAKE) db-rollback
	$(MAKE) db-migrate


.PHONY: db-sqlc-generate
db-sqlc-generate: ## Generate sqlc files in all modules
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	mtools db update-sqlc-config
	mtools db generate


####################################################################################################
## END OF DB COMMANDS
####################################################################################################
