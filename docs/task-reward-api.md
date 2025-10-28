# Task Reward API Documentation

## Overview
This API manages task rewards in the system. It provides endpoints to create, retrieve, update, and delete task rewards with different reward types.

---

## 1. Create Reward

### Endpoint
`POST /api/v1/rewards`

### Request Body
```json
{
  "reward_type": "string",
  "reward_name": "string"
}
```

### Field Validations
- **reward_type**: Required, must be one of: `"crypto_usdt_1"`, `"crypto_usdt_2"`, `"crypto_usdt_3"`
- **reward_name**: Required, string

### Success Response
**Status Code**: `201 Created`

```json
{
  "status": "success",
  "message": "Reward created successfully",
  "data": {
    "reward": 1
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

#### Invalid Reward Type
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "invalid reward type",
  "data": null,
  "errors": null
}
```

**Cause:** The `reward_type` field contains a value that is not one of the valid types

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
- Database constraint violation
- Internal server error

### Example Request

#### cURL
```bash
curl -X POST http://localhost:8080/api/v1/rewards \
  -H "Content-Type: application/json" \
  -d '{
    "reward_type": "crypto_usdt_1",
    "reward_name": "10 USDT Reward"
  }'
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/rewards', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    reward_type: 'crypto_usdt_1',
    reward_name: '10 USDT Reward'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 2. Get Reward by ID

### Endpoint
`GET /api/v1/rewards/{id}`

### Path Parameters
- **id**: Required, integer - The ID of the reward to retrieve

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Reward retrieved successfully",
  "data": {
    "reward": {
      "id": 1,
      "reward_type": "crypto_usdt_1",
      "reward_name": "10 USDT Reward"
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

**Cause:** Invalid ID parameter format

#### Reward Not Found
**Status Code**: `404 Not Found`

```json
{
  "status": "error",
  "message": "reward not found",
  "data": null,
  "errors": null
}
```

**Cause:** No reward exists with the specified ID

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
- Internal server error

### Example Request

#### cURL
```bash
curl -X GET http://localhost:8080/api/v1/rewards/1 \
  -H "Content-Type: application/json"
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/rewards/1', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  }
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 3. Get All Rewards

### Endpoint
`GET /api/v1/rewards`

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Rewards fetched successfully",
  "data": {
    "rewards": [
      {
        "id": 1,
        "reward_type": "crypto_usdt_1",
        "reward_name": "10 USDT Reward"
      },
      {
        "id": 2,
        "reward_type": "crypto_usdt_2",
        "reward_name": "20 USDT Reward"
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
- Internal server error

### Example Request

#### cURL
```bash
curl -X GET http://localhost:8080/api/v1/rewards \
  -H "Content-Type: application/json"
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/rewards', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  }
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 4. Update Reward

### Endpoint
`PUT /api/v1/rewards/{id}`

### Path Parameters
- **id**: Required, integer - The ID of the reward to update

### Request Body
```json
{
  "reward_type": "string",
  "reward_name": "string"
}
```

### Field Validations
- **reward_type**: Optional, must be one of: `"crypto_usdt_1"`, `"crypto_usdt_2"`, `"crypto_usdt_3"` if provided
- **reward_name**: Optional, string
- **Note**: At least one field must be provided

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Reward updated successfully",
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

**Cause:** Malformed JSON in request body or invalid ID parameter

#### Invalid Reward Type
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "invalid reward type",
  "data": null,
  "errors": null
}
```

**Cause:** The `reward_type` field contains a value that is not one of the valid types

#### Reward Not Found
**Status Code**: `404 Not Found`

```json
{
  "status": "error",
  "message": "reward not found",
  "data": null,
  "errors": null
}
```

**Cause:** No reward exists with the specified ID

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
- No fields provided for update
- Internal server error

### Example Request

#### cURL
```bash
curl -X PUT http://localhost:8080/api/v1/rewards/1 \
  -H "Content-Type: application/json" \
  -d '{
    "reward_type": "crypto_usdt_2",
    "reward_name": "Updated 20 USDT Reward"
  }'
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/rewards/1', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    reward_type: 'crypto_usdt_2',
    reward_name: 'Updated 20 USDT Reward'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 5. Delete Reward

### Endpoint
`DELETE /api/v1/rewards/{id}`

### Path Parameters
- **id**: Required, integer - The ID of the reward to delete

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Reward deleted successfully",
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

**Cause:** Invalid ID parameter format

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
- Reward does not exist
- Foreign key constraint violation (if reward is referenced by tasks)
- Internal server error

### Example Request

#### cURL
```bash
curl -X DELETE http://localhost:8080/api/v1/rewards/1 \
  -H "Content-Type: application/json"
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/rewards/1', {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json',
  }
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

---

## Valid Reward Types

The following reward types are currently supported:

| Reward Type      | Description                    |
|------------------|--------------------------------|
| `crypto_usdt_1`  | Crypto USDT Tier 1 Reward      |
| `crypto_usdt_2`  | Crypto USDT Tier 2 Reward      |
| `crypto_usdt_3`  | Crypto USDT Tier 3 Reward      |

---

## Field Requirements Summary

### Create Reward
| Field        | Required | Type   | Valid Values                                           |
|--------------|----------|--------|--------------------------------------------------------|
| reward_type  | ✅ Yes   | string | `crypto_usdt_1`, `crypto_usdt_2`, `crypto_usdt_3`     |
| reward_name  | ✅ Yes   | string | Any characters                                         |

### Update Reward
| Field        | Required | Type   | Valid Values                                           |
|--------------|----------|--------|--------------------------------------------------------|
| reward_type  | ❌ No    | string | `crypto_usdt_1`, `crypto_usdt_2`, `crypto_usdt_3`     |
| reward_name  | ❌ No    | string | Any characters                                         |

**Note:** For update operations, at least one field must be provided.

---

## Notes
- The reward ID is auto-generated upon creation
- Reward types are predefined and cannot be customized
- Deleting a reward may fail if it's referenced by existing tasks (foreign key constraint)
- All endpoints return consistent JSON response format
- The update endpoint supports partial updates (you can update only `reward_type` or only `reward_name`)

---

## Security Considerations
- Consider implementing authentication/authorization for these endpoints
- Use HTTPS in production to prevent man-in-the-middle attacks
- Validate all input data to prevent SQL injection (already handled by parameterized queries)
- Consider implementing rate limiting to prevent abuse
- Log all reward creation, modification, and deletion operations for audit purposes
