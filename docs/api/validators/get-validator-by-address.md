# Get Validator by Address

Retrieves a specific validator by its address.

## Endpoint

```
GET /validators/{address}
```

## Implementation

File: `internal/handlers/validator_handler.go`
Function: `handlers.ValidatorHandler.GetByAddress`

## Description

This endpoint returns detailed information about a specific validator, identified by its unique address.

## Request

### Headers

None required.

### Parameters

**Path Parameters**:
- `address` (required): The validator address to retrieve.

## Response

### Success Response

**Code**: 200 OK

**Content Example**:
```json
{
  "name": "Binance Node",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": true
}
```

### Error Response

**Condition**: If the validator with the specified address does not exist.

**Code**: 404 Not Found

**Content Example**:
```
Validator not found
```

## Sample Call

```bash
curl -X GET http://localhost:8080/validators/cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s
```

## Notes

- The address parameter should be a valid Cosmos validator address.
- The address is case-sensitive. 