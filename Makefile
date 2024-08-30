SHELL := /bin/bash

.PHONY: noop
noop:
	$(info noop)

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
