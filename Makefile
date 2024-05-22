# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=user-management-api
DOCKER_COMPOSE=docker compose

# Docker parameters
DOCKER_COMPOSE_FILE=docker-compose.yml
DOCKER_BUILD_CMD=docker build -t $(BINARY_NAME) .
DOCKER_RUN_CMD=$(DOCKER_COMPOSE) up --build

# Build, run, test, and clean commands
all: test build

# Install Swagger dependencies
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
generate-docs:
	swag init -g internal/interfaces/handler/user_handler.go

# Serve the Swagger documentation using a simple HTTP server
serve-docs:
	@echo "Serving Swagger documentation at http://localhost:8080/swagger/index.html"
	@docker-compose up --build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

run:
	$(DOCKER_RUN_CMD)

test:
	$(GOTEST) ./... -v

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
