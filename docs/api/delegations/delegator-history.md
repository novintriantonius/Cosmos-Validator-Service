# Get Delegator History

Retrieves historical delegation data for a specific delegator with a specific validator.

## Endpoint

```
GET /api/v1/validators/{validator_address}/delegator/{delegator_address}/history
```

## Path Parameters

| Name | Type | Description |
|------|------|-------------|
| `validator_address` | string | The address of the validator |
| `delegator_address` | string | The address of the delegator |

## Response

### Success Response (200 OK)

```json
{
  "status": "success",
  "code": 200,
  "message": "Delegator history retrieved successfully",
  "data": {
    "validator_address": "cosmosvaloper123...",
    "delegator_address": "cosmos456...",
    "history": [
      {
        "id": 5,
        "validator_address": "cosmosvaloper123...",
        "delegator_address": "cosmos456...",
        "delegation_shares": "1500000",
        "created_at": "2023-01-03T14:20:15Z",
        "updated_at": "2023-01-03T14:20:15Z"
      },
      {
        "id": 3,
        "validator_address": "cosmosvaloper123...",
        "delegator_address": "cosmos456...",
        "delegation_shares": "1200000",
        "created_at": "2023-01-02T10:45:00Z",
        "updated_at": "2023-01-02T10:45:00Z"
      },
      {
        "id": 1,
        "validator_address": "cosmosvaloper123...",
        "delegator_address": "cosmos456...",
        "delegation_shares": "1000000",
        "created_at": "2023-01-01T12:15:30Z",
        "updated_at": "2023-01-01T12:15:30Z"
      }
    ],
    "count": 3 // Number of delegation records for this delegator
  }
}
```

### Error Response (404 Not Found)

```json
{
  "status": "error",
  "code": 404,
  "message": "No delegations found for this delegator",
  "errors": [
    "No delegation history found for delegator cosmos456... with validator cosmosvaloper123..."
  ]
}
```

### Error Response (500 Internal Server Error)

```json
{
  "status": "error",
  "code": 500,
  "message": "Failed to retrieve delegations",
  "errors": [
    "Error message describing what went wrong"
  ]
}
```

## Sample Call

```bash
curl -X GET "http://localhost:8080/api/v1/validators/cosmosvaloper123.../delegator/cosmos456.../history"
```

## Notes

- Returns the complete history of delegation records for a specific delegator with a specific validator.
- Delegation records are sorted by creation time (most recent first).
- Each record represents a change in delegation amount.
- Changes in delegation amount can be tracked by comparing the `delegation_shares` field across different records. 