# Cosmos Validator Service

A simple HTTP service with health check and validator management endpoints.

## Features

- Health check endpoint
- Validator CRUD operations
- Delegations retrieval from Cosmos API
- Robust error handling with retry mechanism

## Documentation

Comprehensive documentation is available in the [docs](docs/) directory:

- [Getting Started](docs/getting-started/): Installation and setup
- [API Documentation](docs/api/): Detailed API reference
- [Testing](docs/testing/): How to run and write tests
- [Error Handling](docs/error-handling/): Error handling approach and troubleshooting

## Quick Start

```sh
# Build the service
go build -o cosmos-validator-service cmd/server/main.go

# Run the service
./cosmos-validator-service

# Check the health endpoint
curl http://localhost:8080/health
```

## License

[MIT License](LICENSE) 