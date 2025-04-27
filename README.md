# Cosmos Validator Service

A simple HTTP service with a health check endpoint.

## Features

- Health check endpoint

## Requirements

- Go 1.20 or higher

## Getting Started

### Installation

1. Clone the repository:
```sh
git clone https://github.com/novintriantonius/Cosmos-Validator-Service.git
cd cosmos-validator-service
```

2. Build the application:
```sh
go build -o cosmos-validator-service cmd/server/main.go
```

### Running the Service

```sh
./cosmos-validator-service
```

The service will start on port 8080 by default. You can change the port by setting the `SERVER_PORT` environment variable:

```sh
SERVER_PORT=9000 ./cosmos-validator-service
```

## API Endpoints

### Health Check

```
GET /health
```

Response: 200 OK
```
Service is healthy
```

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port on which the service listens | 8080 | 