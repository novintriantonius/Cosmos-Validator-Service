# Get All Validators

Retrieves a list of all validators in the system.

## Endpoint

```
GET /validators
```

## Implementation

File: `internal/handlers/validator_handler.go`
Function: `handlers.ValidatorHandler.GetAll`

## Description

This endpoint returns a list of all validators currently registered in the system, along with a count of the total number of validators.

## Request

### Headers

None required.

### Parameters

None required.

## Response

### Success Response

**Code**: 200 OK

**Content Example**:
```json
{
  "data": [
    {
      "name": "Binance Node",
      "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
      "enabledTracking": true
    },
    {
      "name": "Another Validator",
      "address": "cosmosvaloper1example2address3here4validator5check",
      "enabledTracking": false
    }
  ],
  "count": 2
}
```

### Error Response

**Condition**: If an internal server error occurs.

**Code**: 500 Internal Server Error

## Sample Call

```bash
curl -X GET http://localhost:8080/api/v1/validators
```

## Notes

- The response includes a `count` field that indicates the total number of validators.
- Validators are returned in no specific order. 