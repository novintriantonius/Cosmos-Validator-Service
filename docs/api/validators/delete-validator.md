# Delete Validator

Deletes a validator from the system.

## Endpoint

```
DELETE /validators/{address}
```

## Implementation

File: `internal/handlers/validator_handler.go`
Function: `handlers.ValidatorHandler.Delete`

## Description

This endpoint allows you to delete a validator from the system, identified by its address.

## Request

### Headers

None required.

### Parameters

**Path Parameters**:
- `address` (required): The validator address to delete.

### Body

No body required.

## Response

### Success Response

**Code**: 204 No Content

No content is returned upon successful deletion.

### Error Response

**Condition**: If the validator with the specified address does not exist.

**Code**: 404 Not Found

**Content Example**:
```
Validator not found
```

## Sample Call

```bash
curl -X DELETE http://localhost:8080/api/v1/validators/cosmosvaloper18ruzecmqj9pv8ac0gvkgryuc7u004te9rh7w5s
```

## Notes

- The deletion is permanent and cannot be undone.
- If the validator is successfully deleted, no content is returned in the response body.
- If you attempt to get a deleted validator, you will receive a 404 Not Found response. 