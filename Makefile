COVERAGE_FILE ?= coverage.out

.PHONY: build lint lint-fix test coverage coverage-html

build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...

coverage:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
