# Create Validator

Creates a new validator in the system.

## Endpoint

```
POST /validators
```

## Implementation

File: `internal/handlers/validator_handler.go`
Function: `handlers.ValidatorHandler.Create`

## Description

This endpoint allows you to create a new validator with the specified details.

## Request

### Headers

**Required**:
- `Content-Type: application/json`

### Body

**Required Fields**:
- `name` (string): The name of the validator.
- `address` (string): The unique Cosmos validator address.
- `enabledTracking` (boolean): Whether tracking is enabled for this validator.

**Example**:
```json
{
  "name": "Binance Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

## Response

### Success Response

**Code**: 201 Created

**Content Example**:
```json
{
  "name": "Binance Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

### Error Responses

**Condition**: If the request body is invalid.

**Code**: 400 Bad Request

**Content Example**:
```
Invalid request body
```

OR

**Condition**: If required fields are missing.

**Code**: 400 Bad Request

**Content Example**:
```
Name is required
```

OR

**Condition**: If a validator with the specified address already exists.

**Code**: 409 Conflict

**Content Example**:
```
Validator with this address already exists
```

## Sample Call

```bash
curl -X POST \
  http://localhost:8080/validators \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Binance Node",
    "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
    "enabledTracking": true
  }'
```

## Notes

- The validator address must be unique.
- Both `name` and `address` fields are required.
- The `enabledTracking` field defaults to `false` if not provided. 