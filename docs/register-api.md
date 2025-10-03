# Register API Documentation

## Endpoint
`POST /register`

## Request Body
```json
{
  "fullname": "string",
  "username": "string",
  "email": "string",
  "password": "string"
}
```

### Field Validations
- **fullname**: Required, max 255 characters
- **username**: Required, max 50 characters
- **email**: Required, valid email format
- **password**: Required, minimum 8 characters

## Success Response
**Status Code**: `201 Created`

```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "user_id": 123
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
    "fullname is required"
  ]
}
```

**Common Validation Errors:**
- `"fullname is required"` - Fullname field is empty
- `"fullname must be less than 255 characters"` - Fullname is too long
- `"username is required"` - Username field is empty
- `"username must be less than 50 character"` - Username is too long
- `"email is required"` - Email field is empty
- `"invalid email format"` - Email format is invalid
- `"password is required"` - Password field is empty
- `"password must be at least 8 characters long"` - Password is too short

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

### Registration Failed
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Registration failed",
  "data": null,
  "errors": null
}
```

**Possible Causes:**
- Email already exists in database (duplicate email)
- Username already exists in database (duplicate username)
- Database connection error
- Internal server error

## Example Request

### cURL
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "fullname": "John Doe",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/auth/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    fullname: 'John Doe',
    username: 'johndoe',
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
- The `bio` field is **no longer part** of the registration process
  - Bio field will be `NULL` or empty string in the database after registration
  - Users can update their `bio` later through a profile update endpoint
- The `fullname` field is **required** during registration
- The response only returns the `user_id` for security purposes (sensitive user information is not exposed)
- Email and username must be unique (duplicates will result in registration failure)
- Password is hashed using bcrypt before storage (never stored in plain text)

## Field Requirements Summary

| Field      | Required | Type   | Min Length | Max Length | Format                    |
|------------|----------|--------|------------|------------|---------------------------|
| fullname   | ✅ Yes   | string | 1          | 255        | Any characters            |
| username   | ✅ Yes   | string | 1          | 50         | Any characters            |
| email      | ✅ Yes   | string | -          | -          | Valid email format        |
| password   | ✅ Yes   | string | 8          | -          | Any characters            |

## Security Considerations
- Passwords are hashed using bcrypt with cost factor 12
- User ID is returned as a numeric value, not a sensitive token
- For authentication, use the `/api/v1/auth/login` endpoint to obtain a JWT token
- All endpoints should be served over HTTPS in production
