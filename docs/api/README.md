# API Documentation

This document provides an overview of the API endpoints provided by the Cosmos Validator Service.

## API Categories

The API is divided into the following categories:

| Category | Description | Documentation |
|----------|-------------|---------------|
| Validators | Endpoints for managing validators | [Validators API](validators/README.md) |
| Health | Endpoints for checking service health | [See below](#health-check) |

## Directory Structure

The API implementation is organized in the following directory structure:

```
internal/
├── handlers/           # API handlers and routing
│   ├── router.go       # Sets up all routes
│   └── validator_handler.go
├── models/             # Data models
│   ├── validator.go
│   └── delegation.go
├── services/           # Business logic and external API calls
│   └── cosmos_service.go
└── store/              # Data storage
    └── validator_store.go
```

## Base URL

All endpoints are relative to the base URL of the service, which defaults to:

```
http://localhost:8080
```

## Health Check

### Check Service Health

```
GET /health
```

**Response**: 200 OK
```
Service is healthy
```

## Authentication

The API currently does not require authentication.

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests:

- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **204 No Content**: Request successful, no content returned
- **400 Bad Request**: Invalid input or parameters
- **404 Not Found**: Requested resource not found
- **409 Conflict**: Resource already exists
- **500 Internal Server Error**: Server-side error
- **502 Bad Gateway**: Error communicating with external API

For more details on error handling, see the [Error Handling documentation](../error-handling/README.md). 