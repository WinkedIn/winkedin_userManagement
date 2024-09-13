# Set the Go version to use
GO_VERSION := 1.20

# Set the linter version
GOLANGCI_LINT_VERSION := v1.53.3

# Install golangci-lint if not already installed
install-lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VERSION); \
	}

# Run go mod tidy and go mod vendor
prepare:
	@echo "Running go mod tidy and go mod vendor..."
	@go mod tidy
	@go mod vendor

# Run the linter on the codebase
lint: install-lint prepare
	@echo "Running golangci-lint..."
	@./bin/golangci-lint run ./...

# Clean up binaries or build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/

# Default target when running "make"
default: lint

