.PHONY: run build test migrate-up migrate-down swagger

run:
	air

build:
	go build -o server ./cmd/server

test:
	go test ./... -v -coverprofile=coverage.out -covermode=atomic

swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# CR-03: Use environment variables for DB connection string in Makefile.
# Set CRUISE_DATABASE_* env vars or source .env before running these targets.
MIGRATE_DB_URL ?= postgresql://$(CRUISE_DATABASE_USER):$(CRUISE_DATABASE_PASSWORD)@$(CRUISE_DATABASE_HOST):$(CRUISE_DATABASE_PORT)/$(CRUISE_DATABASE_DBNAME)?sslmode=$(CRUISE_DATABASE_SSLMODE)

migrate-up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path migrations \
		-database "$(MIGRATE_DB_URL)" \
		-verbose up

migrate-down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path migrations \
		-database "$(MIGRATE_DB_URL)" \
		-verbose down
