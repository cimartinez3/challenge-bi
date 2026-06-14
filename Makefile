BINARY     := novobanco
CMD        := ./cmd/api
GO         := go
GOFLAGS    :=

.PHONY: build run test test-service lint fmt imports vet tidy clean docker-up docker-db docker-down docker-build

build:
	$(GO) build $(GOFLAGS) -o $(BINARY) $(CMD)

fmt:
	gofmt -w ./..

imports:
	goimports -w ./..

validate:
	CGO_ENABLED=0 $(GO) test -v -count=1 ./internal/service/...

coverage:
	CGO_ENABLED=0 $(GO) test -count=1 -coverprofile=coverage.out ./internal/service/...
	$(GO) tool cover -html=coverage.out
