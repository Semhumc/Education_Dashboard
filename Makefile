# Simple Makefile for a Go project

include .env
export

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"


DB_URL=postgres://$(DB_POSTGRES_USER):$(DB_POSTGRES_PASSWORD)@localhost:$(DB_POSTGRES_PORT)/$(DB_POSTGRES_APP_NAME)?sslmode=disable

## Migration komutları
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-force:
	migrate -path migrations -database "$(DB_URL)" force 1

migrate-drop:
	migrate -path migrations -database "$(DB_URL)" drop -f

## Sqlc generate
sqlc:
	sqlc generate

## Hepsi bir arada: migration + sqlc
dev:
	make migrate-up
	make sqlc

migration:
	migrate create -ext sql -dir migrations -seq $(name)


.PHONY: all build run test clean watch docker-run docker-down itest
