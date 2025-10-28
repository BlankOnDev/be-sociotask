# Rewards API Documentation

## Endpoints Overview
- `POST /rewards` - Create a new reward (record task completion by user)

---

## 1. Create Reward

### Endpoint
`POST /rewards`

### Description
Records the completion of a task by a user. This endpoint is used to track when a user completes a specific task and assigns the associated reward.

### Request Body
```json
{
  "user_id": 123,
  "task_id": 456
}
```

### Field Validations
- **user_id**: Required, must be a valid integer greater than 0
- **task_id**: Required, must be a valid integer greater than 0

### Success Response
**Status Code**: `201 Created`

```json
{
  "status": "success",
  "message": "Reward created successfully",
  "data": {
    "reward": {
      "id": 1,
      "user_id": 123,
      "task_id": 456,
      "created_at": "2024-01-15T10:30:00Z"
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

**Cause:** Malformed JSON in request body

#### Validation Failed - Missing User ID
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Validation failed",
  "data": null,
  "errors": [
    "user_id is required and cannot be zero"
  ]
}
```

**Cause:** The `user_id` field is missing, zero, or invalid

#### Validation Failed - Missing Task ID
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "Validation failed",
  "data": null,
  "errors": [
    "task_id is required and cannot be zero"
  ]
}
```

**Cause:** The `task_id` field is missing, zero, or invalid

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
- Failed to insert reward into database
- User ID does not exist in the database
- Task ID does not exist in the database
- User has already completed this task (duplicate reward)

### Example Request

#### cURL
```bash
curl -X POST http://localhost:8080/api/v1/rewards \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "task_id": 456
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
    user_id: 123,
    task_id: 456
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

#### Python (Requests)
```python
import requests

url = "http://localhost:8080/api/v1/rewards"
payload = {
    "user_id": 123,
    "task_id": 456
}
headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, json=payload, headers=headers)
print(response.json())
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

## Field Requirements Summary

### Create Reward
| Field   | Required | Type    | Constraints                    |
|---------|----------|---------|--------------------------------|
| user_id | ✅ Yes   | integer | Must be greater than 0         |
| task_id | ✅ Yes   | integer | Must be greater than 0         |

## Reward Object Structure

The reward object returned in the response contains:

| Field      | Type      | Description                                    |
|------------|-----------|------------------------------------------------|
| id         | integer   | Auto-generated unique identifier for the reward|
| user_id    | integer   | ID of the user who completed the task          |
| task_id    | integer   | ID of the completed task                       |
| created_at | timestamp | Date and time when the reward was created      |

## Business Logic

### Task Completion Flow
1. User completes a task in the application
2. Application sends a POST request to `/rewards` endpoint with `user_id` and `task_id`
3. System validates the request data
4. System records the task completion in the database
5. User receives the reward associated with the task

### Important Notes
- Each user can only complete a task once (duplicate rewards are prevented at the database level)
- Both user and task must exist in the database before creating a reward
- The reward amount is determined by the task configuration, not by this endpoint
- Rewards are recorded immediately and cannot be modified or deleted through this API

## Security Considerations
- **Authentication Required**: This endpoint should be protected with authentication to ensure only authorized users can create rewards
- **Authorization**: Verify that the authenticated user matches the `user_id` in the request to prevent users from creating rewards for others
- **Rate Limiting**: Implement rate limiting to prevent abuse and rapid-fire reward creation attempts
- **Input Validation**: All input is validated to prevent SQL injection (handled by parameterized queries)
- **HTTPS**: Use HTTPS in production to encrypt data in transit
- **Audit Logging**: All reward creation attempts should be logged for audit trails and fraud detection

## Common Use Cases

### 1. User Completes a Social Media Task
```json
POST /rewards
{
  "user_id": 101,
  "task_id": 5
}
```
Records that user 101 completed task 5 (e.g., "Share post on Twitter")

### 2. User Finishes Profile Setup
```json
POST /rewards
{
  "user_id": 202,
  "task_id": 1
}
```
Records that user 202 completed task 1 (e.g., "Complete profile information")

## Error Handling Best Practices

When integrating this API:
1. Always check the `status` field in the response
2. Handle validation errors by displaying the messages from the `errors` array
3. Implement retry logic for 500 errors with exponential backoff
4. Log all error responses for debugging purposes
5. Provide user-friendly error messages in your UI

## Database Constraints

The rewards table has the following constraints:
- Foreign key constraint on `user_id` (references users table)
- Foreign key constraint on `task_id` (references tasks table)
- Unique constraint on (`user_id`, `task_id`) combination to prevent duplicate rewards
- Cascade delete: If a task is deleted, associated rewards are also deleted

## Future Enhancements

Potential future endpoints that may be added:
- `GET /rewards` - Get all rewards
- `GET /rewards/user/:user_id` - Get all rewards for a specific user
- `GET /rewards/task/:task_id` - Get all users who completed a specific task
- `GET /rewards/:id` - Get a specific reward by ID
- `DELETE /rewards/:id` - Remove a reward (admin only)
