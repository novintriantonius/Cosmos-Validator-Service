# Getting Started

This guide will help you set up and run the Cosmos Validator Service.

## Requirements

- Go 1.20 or higher (for direct installation)
- Docker and Docker Compose (for containerized installation)

## Installation

### Option 1: Direct Installation
1. Clone the repository:
```sh
git clone https://github.com/novintriantonius/cosmos-validator-service.git
cd cosmos-validator-service
```

2. Build the application:
```sh
go build -o cosmos-validator-service cmd/server/main.go
```

### Option 2: Docker Installation
1. Clone the repository:
```sh
git clone https://github.com/novintriantonius/cosmos-validator-service.git
cd cosmos-validator-service
```

2. Start the service using Docker Compose:
```sh
docker-compose up -d
```

## Running the Service

### Direct Installation
```sh
./cosmos-validator-service
```

The service will start on port 8080 by default. You can change the port by setting the `SERVER_PORT` environment variable:

```sh
SERVER_PORT=9000 ./cosmos-validator-service
```

### Docker Installation
The service will automatically start when you run `docker-compose up -d`. To stop the service:

```sh
docker-compose down
```

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port on which the service listens | 8080 |
| DB_HOST | PostgreSQL host | postgres |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL username | cosmos |
| DB_PASSWORD | PostgreSQL password | cosmos123 |
| DB_NAME | PostgreSQL database name | cosmos_validator |

## Verifying the Service

Once the service is running, you can verify it by checking the health endpoint:

```sh
curl http://localhost:8080/health
```

You should receive a response indicating the service is healthy:

```
Service is healthy
```

## Docker Setup Details

The Docker setup includes:
- The main application service
- PostgreSQL database
- Automatic health checks
- Volume persistence for database data

The database is automatically initialized and configured when using Docker Compose. The data is persisted in a Docker volume named `postgres_data`. 