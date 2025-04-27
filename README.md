# Cosmos Validator Service

A simple HTTP service with health check and validator management endpoints.

## Features

- Health check endpoint
- Validator CRUD operations

## Requirements

- Go 1.20 or higher

## Getting Started

### Installation

1. Clone the repository:
```sh
git clone https://github.com/novintriantonius/cosmos-validator-service.git
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

### Validator Endpoints

#### Get All Validators

```
GET /validators
```

Response: 200 OK
```json
{
  "data": [
    {
      "name": "Binance Node",
      "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
      "enabledTracking": true
    }
  ],
  "count": 1
}
```

#### Get Validator by Address

```
GET /validators/{address}
```

Response: 200 OK
```json
{
  "name": "Binance Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

#### Create Validator

```
POST /validators
```

Request body:
```json
{
  "name": "Binance Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

Response: 201 Created
```json
{
  "name": "Binance Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

#### Update Validator

```
PUT /validators/{address}
```

Request body:
```json
{
  "name": "Updated Node Name",
  "enabledTracking": false
}
```

Response: 200 OK
```json
{
  "name": "Updated Node Name",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": false
}
```

#### Delete Validator

```
DELETE /validators/{address}
```

Response: 204 No Content

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port on which the service listens | 8080 | 