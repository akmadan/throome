.PHONY: all build clean test run install lint fmt help dev

# Variables
BINARY_NAME=throome
CLI_BINARY_NAME=throome-cli
VERSION?=0.1.0
BUILD_DIR=bin
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Colors for output
GREEN=\033[0;32m
NC=\033[0m # No Color

all: clean build

## help: Display this help message
help:
	@echo "Throome - Gateway Management System"
	@echo ""
	@echo "Available targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'

## build: Build all binaries
build: build-gateway build-cli
	@echo "${GREEN}✓ Build complete${NC}"

## build-gateway: Build the gateway server
build-gateway:
	@echo "Building gateway server..."
	@mkdir -p ${BUILD_DIR}
	@${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/throome

## build-cli: Build the CLI tool
build-cli:
	@echo "Building CLI tool..."
	@mkdir -p ${BUILD_DIR}
	@${GO} build ${GOFLAGS} ${LDFLAGS} -o ${BUILD_DIR}/${CLI_BINARY_NAME} ./cmd/throome-cli

## install: Install binaries to GOPATH/bin
install:
	@echo "Installing binaries..."
	@${GO} install ${LDFLAGS} ./cmd/throome
	@${GO} install ${LDFLAGS} ./cmd/throome-cli
	@echo "${GREEN}✓ Installation complete${NC}"

## run: Run the gateway server
run: build-gateway
	@echo "Starting Throome Gateway..."
	@./${BUILD_DIR}/${BINARY_NAME}

## dev: Run in development mode with auto-reload (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

## test: Run all tests (unit + integration)
test:
	@echo "Running all tests..."
	@INTEGRATION_TESTS=true ${GO} test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "${GREEN}✓ All tests complete${NC}"

## test-unit: Run unit tests only (fast, no external dependencies)
test-unit:
	@echo "Running unit tests..."
	@${GO} test -v -short ./...
	@echo "${GREEN}✓ Unit tests complete${NC}"

## test-integration: Run integration tests (requires Docker services)
test-integration:
	@echo "Running integration tests..."
	@echo "Make sure Docker services are running: cd test && docker-compose up -d"
	@INTEGRATION_TESTS=true ${GO} test -v -run Integration ./test/integration/...
	@echo "${GREEN}✓ Integration tests complete${NC}"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@${GO} test -v -coverprofile=coverage.out ./...
	@${GO} tool cover -html=coverage.out -o coverage.html
	@${GO} tool cover -func=coverage.out | grep "total:" | awk '{print "Total Coverage: " $$3}'
	@echo "${GREEN}✓ Coverage report generated: coverage.html${NC}"

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	@${GO} test -race ./...
	@echo "${GREEN}✓ Race tests complete${NC}"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@${GO} test -v ./...

## test-watch: Watch and run tests on file changes (requires entr)
test-watch:
	@which entr > /dev/null || (echo "Installing entr..." && brew install entr)
	@find . -name "*.go" | entr -c make test-unit

## lint: Run linters
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "Running linters..."
	@golangci-lint run ./...
	@echo "${GREEN}✓ Linting complete${NC}"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@${GO} fmt ./...
	@echo "${GREEN}✓ Formatting complete${NC}"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	@${GO} vet ./...
	@echo "${GREEN}✓ Vet complete${NC}"

## tidy: Tidy go modules
tidy:
	@echo "Tidying modules..."
	@${GO} mod tidy
	@echo "${GREEN}✓ Modules tidied${NC}"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@rm -f coverage.txt coverage.out coverage.html
	@rm -rf tmp/
	@echo "${GREEN}✓ Clean complete${NC}"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@${GO} mod download
	@echo "${GREEN}✓ Dependencies downloaded${NC}"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t throome:${VERSION} -f deployments/docker/Dockerfile .
	@echo "${GREEN}✓ Docker image built${NC}"

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 9000:9000 -v $(PWD)/clusters:/app/clusters throome:${VERSION}

## docker-test: Run tests in Docker
docker-test:
	@echo "Running tests in Docker..."
	@cd test && docker-compose up -d
	@sleep 5
	@INTEGRATION_TESTS=true ${GO} test -v ./...
	@cd test && docker-compose down
	@echo "${GREEN}✓ Docker tests complete${NC}"

## test-setup: Start Docker services for testing
test-setup:
	@echo "Starting test services..."
	@cd test && docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "${GREEN}✓ Test services ready${NC}"
	@echo "  Redis:      localhost:6379"
	@echo "  PostgreSQL: localhost:5432 (user: test, pass: test, db: throome_test)"
	@echo "  Kafka:      localhost:9092"

## test-teardown: Stop Docker services
test-teardown:
	@echo "Stopping test services..."
	@cd test && docker-compose down -v
	@echo "${GREEN}✓ Test services stopped${NC}"

## create-cluster: Create a new cluster (usage: make create-cluster NAME=mycluster)
create-cluster:
	@./${BUILD_DIR}/${CLI_BINARY_NAME} create-cluster --name $(NAME)

## proto: Generate protobuf files (for future use)
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto
	@echo "${GREEN}✓ Protobuf generation complete${NC}"

## benchmark: Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	@${GO} test -bench=. -benchmem ./...

## tools: Install development tools
tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "${GREEN}✓ Tools installed${NC}"

## init-project: Initialize project structure
init-project:
	@echo "Initializing project structure..."
	@mkdir -p clusters configs scripts examples deployments/docker deployments/kubernetes
	@touch clusters/.gitkeep
	@echo "${GREEN}✓ Project initialized${NC}"

.DEFAULT_GOAL := help

