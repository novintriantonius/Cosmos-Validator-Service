# Testing Documentation

This document describes how to run tests for the Cosmos Validator Service.

## Quick Start Guide

To run tests, use the provided test script:

1. Make the script executable (if needed):
   ```sh chmod +x scripts/run_tests.sh```

2. Run unit tests only (quick):
   ```sh ./scripts/run_tests.sh```

3. Run all tests including E2E tests:
   ```sh ./scripts/run_tests.sh full```

## Test Structure

The project includes both unit tests and end-to-end (e2e) tests, organized in a structured directory layout:

```
tests/
├── unit/              # Unit tests
│   ├── services/      # Tests for service layer
│   ├── store/         # Tests for data store
│   └── models/        # Tests for data models
└── e2e/               # End-to-end tests
    ├── validators/    # E2E tests for validator endpoints
    └── delegations/   # E2E tests for delegation functionality
```

## Running Tests

### Using the Test Script

The project includes a test script that makes it easy to run all tests:

#### Running Unit Tests Only

For faster development cycles, run only the unit tests (no external dependencies):

```sh
./scripts/run_tests.sh
```

This command will execute all unit tests and report the results.

#### Running All Tests (Unit + E2E)

To run both unit tests and end-to-end tests:

```sh
./scripts/run_tests.sh full
```

**Note**: E2E tests may make actual API calls to external services, so they might be slower and could fail if the external services are unavailable.

### Running Tests Manually

If you prefer to run tests directly with Go commands:

#### Running All Unit Tests

```sh
go test -v ./tests/unit/...
```

#### Running Specific Test Packages

```sh
# Run only service tests
go test -v ./tests/unit/services/...

# Run only store tests
go test -v ./tests/unit/store/...
```

#### Running Specific Tests

To run a specific test or test function:

```sh
go test -v -run TestGetValidatorByAddress ./tests/unit/store/...
```

#### Running E2E Tests

```sh
# Run all E2E tests
go test -v ./tests/e2e/...

# Run only validator E2E tests
go test -v ./tests/e2e/validators/...

# Run only delegations E2E tests
go test -v ./tests/e2e/delegations/...
```

### Additional Test Options

Go's test command provides several useful flags:

- **-short**: Skip long-running tests
  ```sh
  go test -short ./tests/...
  ```

- **-race**: Enable race condition detection
  ```sh
  go test -race ./tests/...
  ```

- **-cover**: Show test coverage
  ```sh
  go test -cover ./tests/...
  ```

- **-coverprofile**: Generate a coverage profile
  ```sh
  go test -coverprofile=coverage.out ./tests/...
  go tool cover -html=coverage.out  # View coverage in browser
  ```

## Test Types

### Unit Tests

Unit tests focus on testing individual components in isolation:

- **Services Tests**: Test the service layer, with mocked dependencies
- **Store Tests**: Test the data store functionality
- **Models Tests**: Test the data models and their validation

### E2E Tests

End-to-end tests verify that components work together correctly:

- **Validators E2E Tests**: Test the validator API endpoints
- **Delegations E2E Tests**: Test the delegation retrieval from external APIs

## Writing New Tests

When adding new features, be sure to add both unit tests and e2e tests.

### Unit Test Guidelines

- Test each function/method in isolation
- Use mocks for external dependencies
- Focus on testing the logic, not the integration points

### E2E Test Guidelines

- Test complete flows, from API request to response
- For external API calls, consider using test servers with mock responses
- For real external API calls, use the `testing.Short()` flag to skip in short mode 