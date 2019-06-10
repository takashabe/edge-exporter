SUBPACKAGES := $(shell go list ./...)

.DEFAULT_GOAL := help

##### Development

.PHONY: test vet lint

test: ## Run go test
	go test -v $(SUBPACKAGES)

vet: ## Run go vet
	go vet $(SUBPACKAGES)

lint: ## Run golint
	golint $(SUBPACKAGES)

##### Utilities

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
