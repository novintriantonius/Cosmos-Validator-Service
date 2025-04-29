# Get Hourly Delegations

Retrieves delegation data for a validator grouped by hourly snapshots.

## Endpoint

```
GET /api/v1/validators/{validator_address}/delegations/hourly
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
  "message": "Hourly delegations retrieved successfully",
  "data": {
    "validator_address": "cosmosvaloper123...",
    "hourly_delegations": {
      "2023-01-01T12:00:00Z": [
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
          "created_at": "2023-01-01T12:30:45Z",
          "updated_at": "2023-01-01T12:30:45Z"
        }
      ],
      "2023-01-01T13:00:00Z": [
        // More delegations...
      ]
    },
    "count": 2 // Number of hours with delegation data
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
curl -X GET "http://localhost:8080/api/v1/validators/cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s/delegations/hourly"
```

## Notes

- Delegations are grouped by hour, with timestamps normalized to the start of each hour (e.g., "2023-01-01T12:00:00Z").
- Each hour contains an array of delegation records that were created during that hour.
- The response includes a count of how many hours have delegation data. 