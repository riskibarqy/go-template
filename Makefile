# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

MIGRATE=migrate -path databases/migrations -database "$(DB_CONNECTION_STRING)"

.PHONY: migrate-up migrate-down migrate-new migrate-force migrate-version

# Create a new migration file with timestamp prefix
migrate-new:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir databases/migrations -format unix $$name

# Apply all migrations
migrate-up:
	$(MIGRATE) up

# Rollback the last migration
migrate-down:
	$(MIGRATE) down 1

# Force a specific migration version
migrate-force:
	@read -p "Enter version: " version; \
	$(MIGRATE) force $$version

# Check current migration version
migrate-version:
	$(MIGRATE) version
