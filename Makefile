GO = $(shell which go)

.PHONY:
	test \
	fmt \
	vet

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	set -eo pipefail
	$(GO) test ./... -cover

lint:
ifeq ($(shell which golangci-lint),)
	$(error command 'golangci-lint' not found)
else
	golangci-lint run ./....
endif
