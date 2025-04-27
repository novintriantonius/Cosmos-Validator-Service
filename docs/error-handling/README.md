# Error Handling

This document explains the error handling approach used in the Cosmos Validator Service and how to troubleshoot common issues.

## API Error Responses

The service follows a consistent error response format for API calls:

- **400 Bad Request**: Invalid input data
- **404 Not Found**: Resource doesn't exist
- **409 Conflict**: Resource already exists
- **500 Internal Server Error**: Server-side error

Error responses include a text message describing the issue.

## Common Error Types

### Validator Store Errors

The validator store returns several error types:

- `ErrValidatorNotFound`: Returned when trying to access a validator that doesn't exist
- `ErrValidatorAlreadyExists`: Returned when trying to create a validator with an address that already exists

### Cosmos API Service Errors

The Cosmos API service includes robust error handling with:

- Error wrapping to preserve context
- Detailed error messages for troubleshooting
- Consistent error formats

## Retry Mechanism

The Cosmos API service implements an exponential backoff retry mechanism:

- Configurable number of retries (default: 3)
- Configurable delay between retries (default: 500ms)
- Exponential increase in delay time with each retry
- Context cancellation support

## Debugging

If you encounter API errors, check:

1. Validation errors: Make sure request payloads match the required format.
2. Network issues: Check connectivity to external APIs.
3. Server logs: Review logs for detailed error information.

## Example Error Cases

### Missing Validator Example

Request:
```
GET /validators/nonexistent
```

Response:
```
Status: 404 Not Found
Validator not found
```

### Invalid Input Example

Request:
```
POST /validators
{
  "address": "invalid",
  "enabledTracking": true
}
```

Response:
```
Status: 400 Bad Request
Name is required
```

### Duplicate Validator Example

Request:
```
POST /validators
{
  "name": "Duplicate Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

Response (if validator already exists):
```
Status: 409 Conflict
Validator with this address already exists
``` 