# CI/CD Pipeline Documentation

This document describes the automated testing and security scanning pipelines for GoCharge.

## Overview

GoCharge uses GitHub Actions for continuous integration and deployment. The pipelines ensure code quality, security, and reliability across multiple Go versions.

## Workflows

### 1. Main Testing Workflow (`.github/workflows/go.yml`)

Runs comprehensive tests and code quality checks on all pushes and pull requests.

#### Jobs

**Test Job**
- Runs on multiple Go versions: 1.21, 1.22, 1.25
- Steps:
  - Build verification
  - Unit tests with race detector (`-race` flag)
  - Coverage reports (uploaded to Codecov)
- Platforms: Ubuntu Latest

**Lint Job**
- Code quality analysis using `golangci-lint`
- Checks for:
  - Code style violations
  - Unused variables and functions
  - Potential bugs
  - Performance issues
- Configuration: `.golangci.yml`

**Security Job (gosec)**
- Security vulnerability scanning
- Reports findings in SARIF format
- Uploads to GitHub Security Dashboard
- Runs all gosec security checks

**Static Analysis Job**
- Runs `staticcheck` for deep static analysis
- Detects unreachable code, unused variables, and other issues

### 2. Security Scanning Workflow (`.github/workflows/gosec.yml`)

Dedicated workflow for security scanning with gosec.

**Purpose:**
- Consistent security checks across all branches
- Automatic SARIF report generation
- Integration with GitHub Security Dashboard

**Key Features:**
- Runs on push to `main` branch
- Runs on all pull requests to `main`
- `-no-fail` flag ensures workflow succeeds even if issues found
- SARIF output for GitHub integration

## Configuration Files

### `.golangci.yml`

Comprehensive linter configuration with the following checks:

**Enabled Linters:**
- `errcheck` - Checks for unchecked errors
- `gocyclo` - Cyclomatic complexity (max 10)
- `gocritic` - Critical programming errors
- `gofmt` - Code formatting
- `goimports` - Import ordering and unused imports
- `govet` - Go vet analyzer
- `ineffassign` - Inefficient assignments
- `misspell` - Spelling mistakes
- `unused` - Unused variables and functions
- `gosec` - Security issues
- `stylecheck` - Code style violations
- `revive` - Rule-based linter

**Test File Exclusions:**
- Excludes gosec, errcheck from test files

### `.codecov.yml`

Coverage report configuration:

- **Precision:** 2 decimal places
- **Coverage Range:** 70-100%
- **Behavior:** Default (requires code changes in pull requests)
- **Flags:** Unit tests with carryforward enabled

### Trigger Events

All workflows trigger on:

- **Push Events:**
  - `main` branch
  - `feat/*` branches
  - `bugfix/*` branches

- **Pull Request Events:**
  - To `main` branch

## Running Locally

### Run All Tests
```bash
go test -v ./...
```

### Run Tests with Race Detector
```bash
go test -v -race ./...
```

### Run Tests with Coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Linter
```bash
golangci-lint run
```

### Run Security Scanner
```bash
gosec ./...
```

### Run Static Analysis
```bash
staticcheck ./...
```

## GitHub Security Dashboard

Security findings from gosec are automatically uploaded to GitHub's Security Dashboard.

**To view:**
1. Go to repository Settings → Security → Code scanning
2. Or navigate to Security tab → Code scanning alerts

## Test Coverage

- **Current Target:** 70-100% coverage
- **Tracked via:** Codecov integration
- **Reports:** Automatically generated on every test run

## Pull Request Checks

All PR checks must pass before merging:

- ✅ Test (all Go versions)
- ✅ Lint (golangci-lint)
- ✅ Security (gosec)
- ✅ Static Analysis (staticcheck)
- ✅ Coverage (70%+ minimum)

## Known Limitations

- All workflows require GitHub Actions to be enabled
- Some security checks are informational (gosec uses `-no-fail`)
- Static analysis requires additional tools installation

## Future Improvements

Potential enhancements to CI/CD:

1. **Performance Benchmarks** - Track performance across versions
2. **Dependency Scanning** - Automated dependency updates
3. **Docker Image Building** - Build and push container images
4. **Deployment** - Automated deployment to staging/production
5. **API Documentation** - Generate OpenAPI specs
6. **SBOM Generation** - Software bill of materials for supply chain safety

## Troubleshooting

### Tests Fail on My Machine

- Ensure you're using a compatible Go version (1.21+)
- Run: `go mod tidy` to sync dependencies
- Run: `go test -v -race ./...` to check for race conditions

### Linter Errors

- Run locally: `golangci-lint run`
- Most issues can be auto-fixed: `golangci-lint run --fix`

### Security Warnings

- Review gosec findings in GitHub Security Dashboard
- Some warnings may be false positives - use `// #nosec` to suppress if necessary
- Always evaluate security warnings carefully

## Resources

- [golangci-lint Documentation](https://golangci-lint.run/)
- [gosec Documentation](https://github.com/securego/gosec)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Codecov Documentation](https://docs.codecov.com/)
