# Delegation API

This section documents all the available API endpoints for retrieving delegation data.

## Available Endpoints

| Method | Endpoint | Description | Documentation |
|--------|----------|-------------|---------------|
| GET | `/api/v1/validators/{validator_address}/delegations/hourly` | Get hourly delegation snapshots | [Get Hourly Delegations](hourly-delegations.md) |
| GET | `/api/v1/validators/{validator_address}/delegations/daily` | Get daily delegation snapshots | [Get Daily Delegations](daily-delegations.md) |
| GET | `/api/v1/validators/{validator_address}/delegator/{delegator_address}/history` | Get delegation history for a specific delegator | [Get Delegator History](delegator-history.md) |

## Delegation Data Model

The delegation data model contains the following fields:

```json
{
  "id": 123,                                  // The internal ID of the delegation record
  "validator_address": "cosmosvaloper...",    // The validator's address
  "delegator_address": "cosmos...",           // The delegator's address
  "delegation_shares": "1000000",             // The delegation amount in shares
  "created_at": "2023-01-01T12:00:00Z",       // When the delegation was created
  "updated_at": "2023-01-01T12:00:00Z"        // When the delegation was last updated
}
```

## Common Response Format

All endpoints return responses in the following format:

```json
{
  "status": "success",                     // "success" or "error"
  "code": 200,                             // HTTP status code
  "message": "Operation successful",       // Human-readable message
  "data": {                                // Response data (only on success)
    // Endpoint-specific data
  },
  "errors": [                              // Error details (only on error)
    "Error message 1",
    "Error message 2"
  ]
}
```

## Common Errors

- **400 Bad Request**: The request parameters are invalid
- **404 Not Found**: The requested validator or delegator does not exist
- **500 Internal Server Error**: An unexpected error occurred on the server 