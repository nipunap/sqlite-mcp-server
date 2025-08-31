# GitHub Actions CI/CD

This directory contains GitHub Actions workflows for the sqlite-mcp-server project.

## Workflows

### `ci.yml` - Continuous Integration

This workflow runs on:
- Every push to `main` and `develop` branches
- Every pull request targeting `main` and `develop` branches
- Pull request events (opened, synchronize, reopened)

#### Jobs

1. **Test** - Runs tests across multiple Go versions
   - Go versions: 1.21, 1.22, 1.23
   - Runs unit tests with race detection
   - Generates coverage reports
   - Uploads coverage to Codecov (only for Go 1.23)

2. **Lint** - Code quality checks
   - Runs golangci-lint with comprehensive linting rules
   - Checks code formatting, style, and potential issues
   - Uses custom configuration from `.golangci.yml`

3. **Integration Test** - End-to-end testing
   - Builds the server binary
   - Sets up test databases
   - Runs quick integration tests
   - Depends on test and lint jobs passing

4. **Security** - Security scanning
   - Runs Gosec security scanner
   - Uploads results to GitHub Security tab
   - Scans for common security vulnerabilities

5. **Status Check** - PR status summary
   - Summarizes all job results for PR requirements
   - Posts status comments on pull requests
   - Updates existing status comments instead of creating duplicates
   - Required for PR merges

### `release.yml` - Automatic Tagging and Releases

This workflow automatically creates tags and releases when code is merged to the `main` branch.

#### Triggers
- **Automatic**: Every push to `main` branch (after successful PR merge)
- **Manual**: Workflow dispatch with version type selection

#### Jobs

1. **Check Changes** - Analyzes commits for release-worthy changes
   - Examines commit messages since last tag
   - Determines appropriate version bump (major/minor/patch)
   - Skips release if no significant changes

2. **Create Tag** - Generates new version tag and GitHub release
   - Calculates semantic version based on commit analysis
   - Supports manual version override via workflow dispatch
   - Runs final tests before tagging
   - Generates automated changelog
   - Creates annotated Git tag

3. **Build Release Assets** - Cross-platform binary compilation
   - Linux AMD64/ARM64
   - macOS AMD64/ARM64 (Intel/Apple Silicon)
   - Windows AMD64
   - Generates SHA256 checksums

4. **Create GitHub Release** - Publishes release with assets
   - Uploads all platform binaries
   - Includes automated changelog
   - Links to full commit comparison

5. **Notify** - Reports release status
   - Success/failure notifications
   - Links to new release

#### Version Bump Logic
- **Major** (`1.x.x`): Commits with `BREAKING`, `major`, `feat!`, `fix!`
- **Minor** (`x.1.x`): Commits with `feat`, `feature`
- **Patch** (`x.x.1`): All other changes (fixes, docs, etc.)

## Manual Release

You can manually trigger a release from the GitHub Actions tab:

1. Go to **Actions** → **Release and Tagging**
2. Click **Run workflow**
3. Choose:
   - **Version bump type**: `patch`, `minor`, or `major`
   - **Custom version**: Override with specific version (e.g., `v2.1.0`)
4. Click **Run workflow**

This is useful for:
- Creating releases outside of the normal merge cycle
- Fixing version numbering issues
- Creating custom version numbers

## Local Development

You can run the same checks locally using the Makefile:

```bash
# Install golangci-lint and run all CI checks
make ci-local

# Run just the linter
make lint

# Run linter with auto-fix
make lint-fix

# Run tests with race detection
make test-race

# Generate coverage report
make coverage
```

## Configuration Files

- `.golangci.yml` - golangci-lint configuration
  - Enables comprehensive set of linters
  - Customized rules for the project
  - Excludes certain checks for test files

## Pull Request Features

The CI workflow includes several PR-specific enhancements:

### Automated Comments
- **Test Results**: Comments with build and test status
- **Status Summary**: Comprehensive status check with all job results
- **Smart Updates**: Updates existing comments instead of creating duplicates

### Status Checks
- All jobs must pass for PR merge approval
- Clear visual indicators for each job status
- Links to detailed action logs

### Security Integration
- SARIF upload to GitHub Security tab
- Security findings visible in PR conversations
- Automated security scanning on every PR

## Coverage Reports

- Coverage reports are generated for each test run
- Codecov integration provides detailed coverage tracking
- HTML coverage reports are generated locally with `make coverage-html`
- Coverage changes are tracked and reported on PRs

## Badge Status

Add these badges to your README.md:

```markdown
[![CI](https://github.com/nipunap/sqlite-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/nipunap/sqlite-mcp-server/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/nipunap/sqlite-mcp-server/branch/main/graph/badge.svg)](https://codecov.io/gh/nipunap/sqlite-mcp-server)
```
