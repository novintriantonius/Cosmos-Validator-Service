# Getting Started

This guide will help you set up and run the Cosmos Validator Service.

## Requirements

- Go 1.20 or higher

## Installation

1. Clone the repository:
```sh
git clone https://github.com/novintriantonius/cosmos-validator-service.git
cd cosmos-validator-service
```

2. Build the application:
```sh
go build -o cosmos-validator-service cmd/server/main.go
```

## Running the Service

```sh
./cosmos-validator-service
```

The service will start on port 8080 by default. You can change the port by setting the `SERVER_PORT` environment variable:

```sh
SERVER_PORT=9000 ./cosmos-validator-service
```

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port on which the service listens | 8080 |

## Verifying the Service

Once the service is running, you can verify it by checking the health endpoint:

```sh
curl http://localhost:8080/health
```

You should receive a response indicating the service is healthy:

```
Service is healthy
``` 