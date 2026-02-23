.PHONY: run build test migrate-up swag docker-up docker-down

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


run:
	go run cmd/app/main.go

# Run unit tests
test:
	go test -v ./internal/usecase/...

lint: 
	golangci-lint run

swag:
	swag init -g cmd/app/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

# Helper to run everything needed for a fresh start
init: swag docker-up