.PHONY: test coverage coverage-html clean mcp-coverage mcp-coverage-html lint lint-fix build ci-setup

# Test all packages
test:
	go test -v ./...

# Generate coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Generate HTML coverage report
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html

# MCP-specific coverage
mcp-coverage:
	./scripts/test_coverage.sh

# MCP coverage with HTML reports
mcp-coverage-html: mcp-coverage
	@echo "HTML reports generated in coverage/ directory"

# Run tests with race detector
test-race:
	go test -race -v ./...

# Clean up coverage files
clean:
	rm -f coverage.out coverage.html
	rm -rf coverage/

# Run all tests and generate coverage report
test-all: clean test-race coverage coverage-html

# Run MCP tests and generate coverage report
mcp-test-all: clean mcp-coverage mcp-coverage-html

# Build the server binary
build:
	go build -o sqlite-mcp-server ./cmd/server

# Run golangci-lint
lint:
	golangci-lint run

# Run golangci-lint with auto-fix
lint-fix:
	golangci-lint run --fix

# Setup CI dependencies (install golangci-lint locally)
ci-setup:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

# Run all CI checks locally
ci-local: ci-setup lint test-race coverage

# Clean everything including build artifacts
clean-all: clean
	rm -f sqlite-mcp-server
