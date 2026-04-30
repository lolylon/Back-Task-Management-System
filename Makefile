DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= 1234
DB_NAME ?= taskmanager

DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

MIGRATE_PATH := migrations
MIGRATE := migrate

GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m 

.PHONY: help
help: 
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	@echo "$(GREEN)Applying migrations...$(NC)"
	@$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Migrations applied successfully!$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" down 1
	@echo "$(GREEN)Migration rolled back successfully!$(NC)"

.PHONY: migrate-down-all
migrate-down-all: ## Rollback all migrations
	@echo "$(YELLOW)Rolling back all migrations...$(NC)"
	@$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" down
	@echo "$(GREEN)All migrations rolled back successfully!$(NC)"

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: VERSION parameter is required. Usage: make migrate-force VERSION=1$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Forcing migration version to $(VERSION)...$(NC)"
	@$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" force $(VERSION)
	@echo "$(GREEN)Migration version forced to $(VERSION)!$(NC)"

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@echo "$(GREEN)Current migration version:$(NC)"
	@$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" version

.PHONY: migrate-create
migrate-create: ## Create new migration (usage: make migrate-create NAME=create_table_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME parameter is required. Usage: make migrate-create NAME=create_table_name$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)Creating migration files for $(NAME)...$(NC)"
	@$(MIGRATE) create -ext sql -dir $(MIGRATE_PATH) $(NAME)
	@echo "$(GREEN)Migration files created successfully!$(NC)"

.PHONY: migrate-status
migrate-status: ## Show migration status
	@echo "$(GREEN)Migration status:$(NC)"
	@$(MIGRATE) -path $(MIGRATE_PATH) -database "$(DATABASE_URL)" version

.PHONY: db-reset
db-reset: ## Reset database (drop all and re-migrate)
	@echo "$(RED)Warning: This will drop all tables!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [ $$REPLY = "y" ]; then \
		$(MAKE) migrate-down-all; \
		$(MAKE) migrate-up; \
	else \
		echo "Aborted."; \
	fi

.PHONY: install-migrate
install-migrate: ## Install golang-migrate tool
	@echo "$(GREEN)Installing golang-migrate...$(NC)"
	@if command -v brew > /dev/null 2>&1; then \
		brew install golang-migrate; \
	else \
		echo "$(RED)Homebrew not found. Please install golang-migrate manually from https://github.com/golang-migrate/migrate/releases$(NC)"; \
	fi
	@echo "$(GREEN)golang-migrate installed successfully!$(NC)"

.PHONY: run
run: ## Run the application
	@echo "$(GREEN)Starting application...$(NC)"
	@go run .

.PHONY: build
build: ## Build the application
	@echo "$(GREEN)Building application...$(NC)"
	@go build -o bin/taskmanager .

.PHONY: test
test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	@go test ./...

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
