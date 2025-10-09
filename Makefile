.PHONY: build all install help

BINARY=mikros

all: build

build: ## Builds the application locally
	@go build -o $(BINARY) ./cmd/mikros

clean: ## Removes the application binary
	@rm -rf $(BINARY)

install: ## Installs the application
	@go install ./cmd/memed

help: ## Shows all available options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
