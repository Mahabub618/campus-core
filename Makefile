# Variables
APP_NAME := campus-core
BUILD_DIR := bin
MAIN_FILE := cmd/server/main.go
MIGRATION_DIR := internal/database/migrations

# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Default Database connection values if not in .env
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= campus_core
DB_SSLMODE ?= disable

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: all build run test clean fmt vet deps docker-build docker-run migrate-up migrate-down migrate-create migrate-force version help

all: build

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/server $(MAIN_FILE)

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_FILE)

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Vet code
	@echo "Vetting code..."
	go vet ./...

deps: ## Install dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Installing tools..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create: ## Create a new migration file. Usage: make migrate-create NAME=init_schema
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is undefined. Usage: make migrate-create NAME=init_schema"; exit 1; fi
	@echo "Creating migration $(NAME)..."
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

migrate-up: ## Run migrations up
	@echo "Running migrations up..."
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down: ## Run migrations down
	@echo "Running migrations down..."
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down

migrate-force: ## Force migration version. Usage: make migrate-force VERSION=1
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION is undefined. Usage: make migrate-force VERSION=1"; exit 1; fi
	@echo "Forcing migration version $(VERSION)..."
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(VERSION)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME)

help: ## Show help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
