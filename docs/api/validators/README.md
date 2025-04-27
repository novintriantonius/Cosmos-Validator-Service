# Validator API

This section documents all the available API endpoints for managing validators.

## Directory Structure

The validators API is implemented in the following directory structure:

```
internal/
├── handlers/
│   ├── router.go            # Sets up all routes
│   └── validator_handler.go # Handles validator endpoints
├── models/
│   └── validator.go         # Validator data model
└── store/
    └── validator_store.go   # In-memory storage for validators
```

## Available Endpoints

| Method | Endpoint | Description | Documentation |
|--------|----------|-------------|---------------|
| GET | `/validators` | Get all validators | [Get All Validators](get-all-validators.md) |
| GET | `/validators/{address}` | Get a specific validator | [Get Validator by Address](get-validator-by-address.md) |
| POST | `/validators` | Create a new validator | [Create Validator](create-validator.md) |
| PUT | `/validators/{address}` | Update an existing validator | [Update Validator](update-validator.md) |
| DELETE | `/validators/{address}` | Delete a validator | [Delete Validator](delete-validator.md) |

## Data Model

The validator model contains the following fields:

```json
{
  "name": "Validator Name",            // Required: The name of the validator
  "address": "cosmosvaloper...",       // Required: The unique Cosmos validator address
  "enabledTracking": true              // Whether this validator is being tracked
}
```

## Common Errors

- **400 Bad Request**: The request body is invalid or required fields are missing
- **404 Not Found**: The requested validator does not exist
- **409 Conflict**: A validator with the same address already exists
- **500 Internal Server Error**: An unexpected error occurred on the server 