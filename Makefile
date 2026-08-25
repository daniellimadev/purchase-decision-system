.PHONY: help build run test clean docker migrate

help: ## Shows this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Builds the application
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

run-api: ## Runs the API
	go run cmd/api/main.go

run-worker: ## Runs the Worker
	go run cmd/worker/main.go

migrate: ## Runs database migrations
	go run cmd/api/main.go migrate

test: ## Runs the tests
	go test -v ./...

test-coverage: ## Runs the tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Cleans compiled files
	rm -rf bin/
	rm -f coverage.out

docker-build: ## Builds the Docker image
	docker build -t purchase-decision-system .

docker-up: ## Starts the containers
	docker-compose up -d

docker-down: ## Stops the containers
	docker-compose down

docker-logs: ## Shows the container logs
	docker-compose logs -f

deps: ## Installs dependencies
	go mod download
	go mod tidy

lint: ## Runs the linter
	golangci-lint run

fmt: ## Formats the code
	go fmt ./...

.DEFAULT_GOAL := help