#!make
include .env
MIGRATIONS_PATH = ./migrations
DB_ADDR = "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable&search_path=$(DB_SCHEMA)"

# Install dependencies
dependencies:
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/air-verse/air@latest


# Create migration file
migration:
	@if [ -z "$(NAME)" ]; then \
		echo "Please provide a migration name using: make migration NAME=your_migration_name"; \
		exit 1; \
	fi
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(NAME)

# Migrate up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

# Migrate down
migrate-down:
	@if [ -z "$(VERSION)" ]; then \
		echo "Please provide a migration version using: make migration VERSION=migration_version"; \
		exit 1; \
	fi
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(VERSION)

# Build the application
build:
	@echo "Building..."	
	@go build -o bin/main cmd/api/main.go

# Generate api mocks
generate-mocks:
		mockgen -package=mocks -source=internal/domain/services/api.go -destination=mocks/services_mock.go
		mockgen -package=mocks -source=internal/domain/repositories/api.go -destination=mocks/repositories_mock.go

# Generate api docs
generate-docs:
	swag init -g cmd/api/main.go -o ./docs

# Run the application
run:
	@go run cmd/api/main.go
	
# Run docker compose
up:
	docker compose up --build --force-recreate

# Shutdown docker compose
down:
	docker-compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/app

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

