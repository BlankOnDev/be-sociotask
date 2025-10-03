# Login API Documentation

## Endpoint
`POST /login`

## Request Body
```json
{
  "email": "string",
  "password": "string"
}
```

### Field Validations
- **email**: Required, valid email format
- **password**: Required

## Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "errors": null
}
```

## Error Responses

### Validation Error
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Validation failed",
  "data": null,
  "errors": [
    "email is required"
  ]
}
```

**Common Validation Errors:**
- `"email is required"` - Email field is empty
- `"invalid email format"` - Email format is invalid
- `"password is required"` - Password field is empty

### Invalid Request
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

### User Not Found
**Status Code**: `404 Not Found`

```json
{
  "status": "error",
  "message": "Not found",
  "data": null,
  "errors": null
}
```

**Cause:** Email address does not exist in the database

### Invalid Credentials
**Status Code**: `401 Unauthorized`

```json
{
  "status": "error",
  "message": "Invalid credentials",
  "data": null,
  "errors": null
}
```

**Cause:** Password does not match the user's stored password

### Internal Server Error
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
- Token generation error
- Password hashing error
- Internal server error

## Example Request

### cURL
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'john.doe@example.com',
    password: 'password123'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

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

## Notes
- The `token` returned is a JWT (JSON Web Token) that should be used for authenticated requests
- The token contains the user's ID and role information
- Store the token securely (e.g., in localStorage or httpOnly cookies)
- Include the token in subsequent requests using the `Authorization` header: `Bearer <token>`
- Password validation is performed using bcrypt comparison
- Login attempts with incorrect passwords will not reveal whether the email exists

## Field Requirements Summary

| Field    | Required | Type   | Min Length | Max Length | Format             |
|----------|----------|--------|------------|-----------|--------------------|
| email    | ✅ Yes   | string | -          | -         | Valid email format |
| password | ✅ Yes   | string | -          | -         | Any characters     |

## Security Considerations
- Passwords are never stored in plain text, only bcrypt hashes are compared
- JWT tokens are signed with a secret key to prevent tampering
- Use HTTPS in production to prevent man-in-the-middle attacks
- Consider implementing rate limiting to prevent brute force attacks
- Token expiration should be configured appropriately for your security requirements
- Failed login attempts should not reveal whether the email exists in the system

## Using the Token
After successful login, use the returned token for authenticated requests:

```bash
curl -X GET http://localhost:8080/api/v1/protected-endpoint \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

```javascript
fetch('http://localhost:8080/api/v1/protected-endpoint', {
  method: 'GET',
  headers: {
    'Authorization': 'Bearer ' + token
  }
})
.then(response => response.json())
.then(data => console.log(data));
```
