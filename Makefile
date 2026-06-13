BINARY     := novobanco
CMD        := ./cmd/api
GO         := go
GOFLAGS    :=

.PHONY: build run test lint fmt imports vet tidy clean docker-up docker-db docker-down docker-build

## build: compile the binary
build:
	$(GO) build $(GOFLAGS) -o $(BINARY) $(CMD)

## run: build and run locally (requires DATABASE_URL set)
run: build
	./$(BINARY)

## test: run all tests with race detector
test:
	$(GO) test ./... -race -count=1 -timeout 60s

## fmt: format code with gofmt
fmt:
	gofmt -w ./..

## imports: format code and fix imports with goimports
imports:
	goimports -w ./..

## vet: run go vet
vet:
	$(GO) vet ./...

## lint: vet + imports check
lint: vet imports

## tidy: tidy and verify go.mod
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## cover: run tests and open html coverage report
cover:
	$(GO) test ./... -race -count=1 -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out

## clean: remove binary and coverage artifacts
clean:
	rm -f $(BINARY) coverage.out

## docker-up: start all services via docker compose
docker-up:
	docker compose up -d

## docker-db: start only postgres (for local go run)
docker-db:
	docker compose up -d db

## docker-down: stop and remove containers
docker-down:
	docker compose down

## docker-build: build the Docker image
docker-build:
	docker build -t $(BINARY):local .

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## //'
