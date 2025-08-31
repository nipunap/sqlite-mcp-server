# Local Testing for GitHub Actions

This document explains how to test your GitHub Actions workflows locally before pushing to GitHub.

## 🧪 Methods to Test GitHub Actions Locally

### 1. **Direct Command Testing** (Recommended)

Run the exact same commands that GitHub Actions uses:

```bash
# Test linting (same as CI)
$(go env GOPATH)/bin/golangci-lint run --timeout=5m

# Test building
go build -v ./...

# Test with race detection and coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Run security scanner (if you have it)
# gosec ./...
```

### 2. **Using Makefile** (Convenient)

Use the local CI command:

```bash
# Install dependencies and run all CI checks
make ci-local

# Individual commands
make lint          # Run linting
make test-race     # Run tests with race detection
make coverage      # Generate coverage report
make build         # Build the project
```

### 3. **Using `act` Tool** (Full GitHub Actions Simulation)

Install and use `act` to run GitHub Actions locally with Docker:

```bash
# Install act (requires Docker)
brew install act

# Run specific jobs
act -j lint                    # Run just the lint job
act -j test                    # Run just the test job
act -j integration-test        # Run integration tests
act                           # Run all jobs

# Run with specific event
act push                      # Simulate push event
act pull_request             # Simulate PR event
```

### 4. **VS Code Extensions**

- **GitHub Actions**: Syntax highlighting and validation
- **YAML**: Better YAML editing experience

## 🔧 Local Setup

### Install Required Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install goimports (for formatting)
go install golang.org/x/tools/cmd/goimports@latest

# Install act (optional, requires Docker)
brew install act
```

### Verify Installation

```bash
# Check golangci-lint
$(go env GOPATH)/bin/golangci-lint --version

# Check Go tools
go version
```

## 🚀 Quick Test Script

Create a `test-ci.sh` script:

```bash
#!/bin/bash
set -e

echo "🔍 Running local CI tests..."

echo "📦 Building..."
go build -v ./...

echo "🧪 Running tests..."
go test -v -race ./...

echo "🔍 Running linter..."
$(go env GOPATH)/bin/golangci-lint run --timeout=5m

echo "✅ All local CI tests passed!"
```

Make it executable and run:

```bash
chmod +x test-ci.sh
./test-ci.sh
```

## 📋 Common Issues and Solutions

### Linting Issues

**Problem**: `golangci-lint` reports many issues
**Solution**:
1. Run `$(go env GOPATH)/bin/goimports -w .` to fix imports
2. Check `.golangci.yml` configuration
3. Fix critical issues, exclude overly strict rules

**Problem**: Import formatting issues
**Solution**:
```bash
$(go env GOPATH)/bin/goimports -w .
```

### Test Issues

**Problem**: Tests hang or timeout
**Solution**:
1. Add `t.Parallel()` to tests
2. Use in-memory databases for testing
3. Add proper cleanup functions
4. Set reasonable timeouts

### Build Issues

**Problem**: Build fails locally but works in CI
**Solution**:
1. Check Go version compatibility
2. Run `go mod tidy`
3. Ensure all dependencies are available

## 🎯 Best Practices

1. **Test Early**: Run local tests before every commit
2. **Use Makefile**: Standardize common commands
3. **Parallel Tests**: Use `t.Parallel()` for faster execution
4. **Clean Setup**: Use proper test setup and cleanup
5. **Reasonable Timeouts**: Don't let tests hang indefinitely

## 📊 Performance Tips

- **Parallel Execution**: Tests run in parallel by default
- **In-Memory Databases**: Use `:memory:` for SQLite tests
- **Minimal Setup**: Only create what you need for tests
- **Proper Cleanup**: Clean up resources to avoid conflicts

## 🔗 Related Files

- `.golangci.yml` - Linting configuration
- `Makefile` - Build and test commands
- `.github/workflows/ci.yml` - GitHub Actions workflow
- `.github/workflows/release.yml` - Release workflow

This approach ensures your code passes CI before you push, saving time and avoiding failed builds!
