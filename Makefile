SHELL := /bin/bash

.PHONY: all
all: noop

.PHONY: noop
noop: ;

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: format
format:
	@for f in $$(find . -name '*.go' -print); do \
		gofmt -s -w $$f ; \
	done

.PHONY: help
help:
	$(info available make targets)
	$(info ----------------------)
	@grep "^[a-z0-9-][a-z0-9.-]*:" Makefile  | cut -d':' -f1 | sort
