# Update Validator

Updates an existing validator's information.

## Endpoint

```
PUT /validators/{address}
```

## Implementation

File: `internal/handlers/validator_handler.go`
Function: `handlers.ValidatorHandler.Update`

## Description

This endpoint allows you to update the information of an existing validator, identified by its address.

## Request

### Headers

**Required**:
- `Content-Type: application/json`

### Parameters

**Path Parameters**:
- `address` (required): The validator address to update.

### Body

**Optional Fields** (at least one required):
- `name` (string): The new name for the validator.
- `enabledTracking` (boolean): Whether tracking should be enabled for this validator.

**Example**:
```json
{
  "name": "Updated Node Name",
  "enabledTracking": false
}
```

## Response

### Success Response

**Code**: 200 OK

**Content Example**:
```json
{
  "name": "Updated Node Name",
  "address": "cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s",
  "enabledTracking": false
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

**Condition**: If the validator with the specified address does not exist.

**Code**: 404 Not Found

**Content Example**:
```
Validator not found
```

## Sample Call

```bash
curl -X PUT \
  http://localhost:8080/validators/cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Updated Node Name",
    "enabledTracking": false
  }'
```

## Notes

- The address in the path parameter identifies the validator to update.
- The address of a validator cannot be changed.
- You can update either the name, the enabledTracking flag, or both.
- Omitted fields will retain their current values. 