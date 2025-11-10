# Get Current User API Documentation

## Endpoint
`GET /users/current`

## Description
Retrieves the profile information of the currently authenticated user. This endpoint requires authentication via JWT token.

## Authentication
**Required**: Yes

**Type**: Bearer Token

**Header**:
```
Authorization: Bearer <your_jwt_token>
```

## Request Body
No request body required.

## Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "johndoe@example.com",
      "fullname": "John Doe",
      "x_id": null,
      "wallet_address": null,
      "created_at": "2025-11-10T10:00:00Z"
    }
  },
  "errors": null
}
```

### Response Fields
- **id**: User's unique identifier
- **username**: User's username
- **email**: User's email address
- **fullname**: User's full name
- **x_id**: Twitter/X account ID (nullable)
- **wallet_address**: User's crypto wallet address (nullable)
- **created_at**: Account creation timestamp

## Error Responses

### Unauthorized - Missing/Invalid Token
**Status Code**: `401 Unauthorized`

```json
{
  "status": "error",
  "message": "Unauthorized",
  "data": null,
  "errors": ["authentication required"]
}
```

**Possible Causes:**
- No Authorization header provided
- Invalid JWT token
- Expired JWT token
- Malformed token format

### Bad Request - Invalid Credentials
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Bad request",
  "data": {
    "error": "invalid credentials"
  },
  "errors": null
}
```

**Possible Causes:**
- User context not found in request
- Token valid but user data corrupted

### Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": {
    "error": "intetrnal server error"
  },
  "errors": null
}
```

**Possible Causes:**
- Database connection error
- User not found in database
- Server error while retrieving user data