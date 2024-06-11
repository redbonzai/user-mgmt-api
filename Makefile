# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) ginkgo
GOGET=$(GOCMD) get
BINARY_NAME=user-management-api
DOCKER_COMPOSE=docker compose

# Docker parameters
DOCKER_COMPOSE_FILE=docker-compose.yml
DOCKER_BUILD_CMD=docker build -t $(BINARY_NAME) .
DOCKER_RUN_CMD=$(DOCKER_COMPOSE) up --build

# Migration parameters
MIGRATE_CMD=migrate
MIGRATE_DIR=internal/db/migrations
DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DB)?sslmode=disable



# Build, run, test, and clean commands
all: ginkgo build

# Install Swagger dependencies
install-swag:
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/swaggo/echo-swagger

swag-init:
	swag init --dir cmd/server --output docs/

# Serve the Swagger documentation using a simple HTTP server
serve-docs:
	@echo "Serving Swagger documentation at http://localhost:8080/swagger/index.html"
	@docker-compose up --build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

run:
	$(DOCKER_RUN_CMD)

repository-mocks:
	mockgen -source=internal/interfaces/repository/user_repository.go -destination=internal/interfaces/repository/mocks/mock_user_repository.go -package=mocks

service-mocks:
	mockgen -source=internal/services/user_service.go -destination=internal/services/mocks/mock_service.go -package=mocks



# Run tests
ginkgo:
	ginkgo -r -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Commands for managing dependencies
deps:
	$(GOGET) -u github.com/labstack/echo/v4
	$(GOGET) -u github.com/Masterminds/squirrel
	$(GOGET) -u github.com/onsi/ginkgo/ginkgo
	$(GOGET) -u github.com/onsi/gomega
	$(GOGET) -u github.com/swaggo/echo-swagger
	$(GOGET) -u go.uber.org/zap
	$(GOGET) -u github.com/golang-migrate/migrate/v4/cmd/migrate

# Commands for generating Swagger documentation
swagger-init:
	swag init -g cmd/server/main.go

# Docker build and run commands
docker-build:
	$(DOCKER_BUILD_CMD)

docker-run:
	$(DOCKER_RUN_CMD)

lint:
	golangci-lint run

# Migration commands
migration-create:
	@read -p "Enter migration name: " name; \
	migrate create -dir $(MIGRATE_DIR) -ext sql $$name

migration-up:
	@export $(cat .env | xargs) && \
	$(MIGRATE_CMD) -path $(MIGRATE_DIR) -database "postgres://root:admin@localhost:5432/userapi?sslmode=disable" up

migration-down:
	@export $(cat .env | xargs) && \
	$(MIGRATE_CMD) -path $(MIGRATE_DIR) -database "postgres://root:admin@localhost:5432/userapi?sslmode=disable" down 1

