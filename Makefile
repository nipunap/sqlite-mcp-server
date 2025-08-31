.PHONY: test coverage coverage-html clean mcp-coverage mcp-coverage-html

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
