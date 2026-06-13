BINARY     := novobanco
CMD        := ./cmd/api
GO         := go
GOFLAGS    :=

.PHONY: build run test lint fmt imports vet tidy clean docker-up docker-down

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
