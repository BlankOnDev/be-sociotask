# Task Action API Documentation

## Endpoints Overview
- `POST /actions` - Create a new task action
- `GET /actions` - Get all task actions
- `GET /actions/:id` - Get a specific task action by ID
- `PUT /actions/:id` - Update a task action
- `DELETE /actions/:id` - Delete a task action

---

## 1. Create Task Action

### Endpoint
`POST /actions`

### Request Body
```json
{
  "type": "type_1" | "type_2" | "type_3",
  "name": "string",
  "description": "string"
}
```

### Field Validations
- **type**: Required, must be one of: `"type_1"`, `"type_2"`, `"type_3"`
- **name**: Required
- **description**: Required

### Success Response
**Status Code**: `201 Created`

```json
{
  "status": "success",
  "message": "Action created successfully",
  "data": {
    "action": 123
  },
  "errors": null
}
```

### Error Responses

#### Invalid Request
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Invalid request",
  "data": null,
  "errors": null
}
```

**Cause:** Malformed JSON in request body

#### Invalid Task Type
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Invalid action type",
  "data": null,
  "errors": null
}
```

**Cause:** The `type` field contains a value other than `"type_1"`, `"type_2"`, or `"type_3"`

#### Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null,
  "errors": null
}
```

**Possible Causes:**
- Database connection error
- Failed to insert action into database

### Example Request

#### cURL
```bash
curl -X POST http://localhost:8080/api/v1/actions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "type_1",
    "name": "Complete Profile",
    "description": "Fill out your profile information completely"
  }'
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/actions', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    type: 'type_1',
    name: 'Complete Profile',
    description: 'Fill out your profile information completely'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 2. Get All Task Actions

### Endpoint
`GET /actions`

### Request Parameters
None

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Tasks fetched successfully",
  "data": {
    "action": [
      {
        "id": 1,
        "type": "type_1",
        "name": "Complete Profile",
        "description": "Fill out your profile information completely"
      },
      {
        "id": 2,
        "type": "type_2",
        "name": "Share Post",
        "description": "Share a post on social media"
      }
    ]
  },
  "errors": null
}
```

### Error Responses

#### Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null,
  "errors": null
}
```

**Possible Causes:**
- Database connection error
- Failed to retrieve actions from database

### Example Request

#### cURL
```bash
curl -X GET http://localhost:8080/api/v1/actions
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/actions', {
  method: 'GET'
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 3. Get Task Action by ID

### Endpoint
`GET /actions/:id`

### URL Parameters
- **id**: Integer, the ID of the task action to retrieve

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Action retrieved successfully",
  "data": {
    "action": {
      "id": 1,
      "type": "type_1",
      "name": "Complete Profile",
      "description": "Fill out your profile information completely"
    }
  },
  "errors": null
}
```

### Error Responses

#### Invalid Request
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Invalid request",
  "data": null,
  "errors": null
}
```

**Cause:** Invalid ID parameter (not a valid integer)

#### Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null,
  "errors": null
}
```

**Possible Causes:**
- Database connection error
- Action with specified ID does not exist
- Failed to retrieve action from database

### Example Request

#### cURL
```bash
curl -X GET http://localhost:8080/api/v1/actions/1
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/actions/1', {
  method: 'GET'
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 4. Update Task Action

### Endpoint
`PUT /actions/:id`

### URL Parameters
- **id**: Integer, the ID of the task action to update

### Request Body
```json
{
  "type": "type_1" | "type_2" | "type_3",
  "name": "string",
  "description": "string"
}
```

### Field Validations
- **type**: Optional, if provided must be one of: `"type_1"`, `"type_2"`, `"type_3"`
- **name**: Optional
- **description**: Optional
- At least one field must be provided for update

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Action updated successfully",
  "data": null,
  "errors": null
}
```

### Error Responses

#### Invalid Request
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Invalid request",
  "data": null,
  "errors": null
}
```

**Causes:**
- Invalid ID parameter (not a valid integer)
- Malformed JSON in request body

#### Invalid Task Type
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "invalid task type",
  "data": null,
  "errors": null
}
```

**Cause:** The `type` field contains a value other than `"type_1"`, `"type_2"`, or `"type_3"`

#### Action Not Found
**Status Code**: `404 Not Found`

```json
{
  "status": "error",
  "message": "action not found",
  "data": null,
  "errors": null
}
```

**Cause:** Action with the specified ID does not exist in the database

#### Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null,
  "errors": null
}
```

**Possible Causes:**
- Database connection error
- Failed to update action in database
- No fields provided for update

### Example Request

#### cURL
```bash
curl -X PUT http://localhost:8080/api/v1/actions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Complete Your Profile",
    "description": "Fill out all required profile information"
  }'
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/actions/1', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'Complete Your Profile',
    description: 'Fill out all required profile information'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 5. Delete Task Action

### Endpoint
`DELETE /actions/:id`

### URL Parameters
- **id**: Integer, the ID of the task action to delete

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Action deleted successfully",
  "data": null,
  "errors": null
}
```

### Error Responses

#### Invalid Request
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Invalid request",
  "data": null,
  "errors": null
}
```

**Cause:** Invalid ID parameter (not a valid integer)

#### Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null,
  "errors": null
}
```

**Possible Causes:**
- Database connection error
- Failed to delete action from database
- Action with specified ID does not exist

### Example Request

#### cURL
```bash
curl -X DELETE http://localhost:8080/api/v1/actions/1
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/actions/1', {
  method: 'DELETE'
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## Response Format
All responses follow a consistent format:

```json
{
  "status": "success" | "error",
  "message": "string",
  "data": {
    // Response data (only on success)
  } | null,
  "errors": [
    // Array of error messages (only on validation errors)
  ] | null
}
```

## Task Action Types

The system supports three types of task actions:

| Type    | Value     | Description                           |
|---------|-----------|---------------------------------------|
| Type 1  | `type_1`  | First category of task actions        |
| Type 2  | `type_2`  | Second category of task actions       |
| Type 3  | `type_3`  | Third category of task actions        |

## Field Requirements Summary

### Create Action
| Field       | Required | Type   | Constraints                           |
|-------------|----------|--------|---------------------------------------|
| type        | ✅ Yes   | string | Must be: type_1, type_2, or type_3    |
| name        | ✅ Yes   | string | Any characters                        |
| description | ✅ Yes   | string | Any characters                        |

### Update Action
| Field       | Required | Type   | Constraints                           |
|-------------|----------|--------|---------------------------------------|
| type        | ❌ No    | string | Must be: type_1, type_2, or type_3    |
| name        | ❌ No    | string | Any characters                        |
| description | ❌ No    | string | Any characters                        |

*Note: At least one field must be provided when updating*

## Notes
- All endpoints return consistent JSON response format
- The `id` field is auto-generated by the database and returned after creation
- Update operations support partial updates (only specified fields are updated)
- The update endpoint uses dynamic query building to only update provided fields
- Empty strings are not considered valid values for optional fields in updates
- Task action types are predefined and cannot be customized

## Security Considerations
- Consider implementing authentication/authorization for these endpoints
- Validate all input data to prevent SQL injection (currently handled by parameterized queries)
- Use HTTPS in production to encrypt data in transit
- Consider implementing rate limiting to prevent abuse
- Add proper logging for audit trails
