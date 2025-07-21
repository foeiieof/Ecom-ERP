.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Go commands
build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf tmp/

# Development with hot reload
dev:
	air

# Dependency management
deps:
	go mod download
	go mod tidy

# Docker commands
docker-build:
	docker build -t ecommerce-api .

docker-run:
	docker run -p 8080:8080 --env-file .env ecommerce-api

# Docker Compose commands
docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f

# Linting and formatting
fmt:
	go fmt ./...

lint:
	golangci-lint run

# Security check
security:
	gosec ./...

# Generate documentation
docs:
	swag init -g cmd/server/main.go

# Database migrations (if using a database)
migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  dev                - Run with hot reload (requires air)"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-compose-up  - Start with docker-compose"
	@echo "  docker-compose-down- Stop docker-compose"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  security           - Run security check"
	@echo "  help               - Show this help"
