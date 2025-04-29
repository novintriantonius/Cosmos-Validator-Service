# Get Daily Delegations

Retrieves delegation data for a validator grouped by daily snapshots.

## Endpoint

```
GET /api/v1/validators/{validator_address}/delegations/daily
```

## Path Parameters

| Name | Type | Description |
|------|------|-------------|
| `validator_address` | string | The address of the validator |

## Response

### Success Response (200 OK)

```json
{
  "status": "success",
  "code": 200,
  "message": "Daily delegations retrieved successfully",
  "data": {
    "validator_address": "cosmosvaloper123...",
    "daily_delegations": {
      "2023-01-01T00:00:00Z": [
        {
          "id": 1,
          "validator_address": "cosmosvaloper123...",
          "delegator_address": "cosmos456...",
          "delegation_shares": "1000000",
          "created_at": "2023-01-01T12:15:30Z",
          "updated_at": "2023-01-01T12:15:30Z"
        },
        {
          "id": 2,
          "validator_address": "cosmosvaloper123...",
          "delegator_address": "cosmos789...",
          "delegation_shares": "2000000",
          "created_at": "2023-01-01T16:30:45Z",
          "updated_at": "2023-01-01T16:30:45Z"
        }
      ],
      "2023-01-02T00:00:00Z": [
        // More delegations...
      ]
    },
    "count": 2 // Number of days with delegation data
  }
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
curl -X GET "http://localhost:8080/api/v1/validators/cosmosvaloper123.../delegations/daily"
```

## Notes

- Delegations are grouped by day, with timestamps normalized to the start of each day (e.g., "2023-01-01T00:00:00Z").
- Each day contains an array of delegation records that were created during that day.
- The response includes a count of how many days have delegation data. 