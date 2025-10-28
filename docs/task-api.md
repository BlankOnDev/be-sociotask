# Task API Documentation

## Endpoints Overview
- [Create Task](#create-task) - `POST /tasks`
- [Get Task by ID](#get-task-by-id) - `GET /tasks/{id}`
- [Get All Tasks](#get-all-tasks) - `GET /tasks`
- [Edit Task](#edit-task) - `PUT /tasks/{id}`
- [Delete Task](#delete-task) - `DELETE /tasks/{id}`

---

## Create Task

### Endpoint
`POST /tasks`

### Authentication
**Required**: Yes (JWT Token)

### Request Headers
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Request Body
```json
{
  "title": "string",
  "description": "string",
  "reward_task": 1,
  "reward_usdt": 100.50,
  "due_date": "2024-12-31T23:59:59Z",
  "max_participant": "50",
  "task_image": "https://example.com/image.jpg",
  "action_id": 1
}
```

### Field Descriptions
- **title**: Task title
- **description**: Detailed task description
- **reward_task**: Reward ID reference
- **reward_usdt**: Reward amount in USDT
- **due_date**: Task deadline (ISO 8601 format)
- **max_participant**: Maximum number of participants (string)
- **task_image**: URL to task image
- **action_id**: Task action type ID

### Success Response
**Status Code**: `201 Created`

```json
{
  "status": "success",
  "message": "Task created successfully",
  "data": {
    "task": {
      "id": 1,
      "title": "Complete Social Media Task",
      "description": "Follow and share our content",
      "user_id": 123,
      "reward_task": 1,
      "reward_usdt": 100.50,
      "due_date": "2024-12-31T23:59:59Z",
      "max_participant": "50",
      "created_at": "2024-10-28T10:00:00Z",
      "task_image": "https://example.com/image.jpg",
      "action_id": 1,
      "updated_at": "2024-10-28T10:00:00Z"
    }
  },
  "errors": null
}
```

### Error Responses

#### Unauthorized
**Status Code**: `401 Unauthorized`

```json
{
  "status": "error",
  "message": "Unauthorized",
  "data": null,
  "errors": null
}
```

**Cause:** Missing or invalid JWT token

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
- Failed to create task in database

### Example Request

#### cURL
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete Social Media Task",
    "description": "Follow and share our content",
    "reward_task": 1,
    "reward_usdt": 100.50,
    "due_date": "2024-12-31T23:59:59Z",
    "max_participant": "50",
    "task_image": "https://example.com/image.jpg",
    "action_id": 1
  }'
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/tasks', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_JWT_TOKEN',
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'Complete Social Media Task',
    description: 'Follow and share our content',
    reward_task: 1,
    reward_usdt: 100.50,
    due_date: '2024-12-31T23:59:59Z',
    max_participant: '50',
    task_image: 'https://example.com/image.jpg',
    action_id: 1
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## Get Task by ID

### Endpoint
`GET /tasks/{id}`

### Authentication
**Required**: No

### Path Parameters
- **id**: Task ID (integer)

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Task retrieved successfully",
  "data": {
    "task": {
      "id": 1,
      "title": "Complete Social Media Task",
      "description": "Follow and share our content",
      "user_id": 123,
      "reward_task": 1,
      "reward_usdt": 100.50,
      "due_date": "2024-12-31T23:59:59Z",
      "max_participant": "50",
      "created_at": "2024-10-28T10:00:00Z",
      "task_image": "https://example.com/image.jpg",
      "action_id": 1,
      "updated_at": "2024-10-28T10:00:00Z"
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

**Cause:** Invalid task ID parameter

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
- Task not found
- Database connection error

### Example Request

#### cURL
```bash
curl -X GET http://localhost:8080/api/v1/tasks/1
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/tasks/1')
  .then(response => response.json())
  .then(data => console.log(data));
```

---

## Get All Tasks

### Endpoint
`GET /tasks`

### Authentication
**Required**: No

### Query Parameters
- **page**: Page number (default: 1)

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Tasks fetched successfully",
  "data": {
    "tasks": [
      {
        "id": 1,
        "title": "Complete Social Media Task",
        "description": "Follow and share our content",
        "user_id": 123,
        "reward_task": 1,
        "reward_usdt": 100.50,
        "due_date": "2024-12-31T23:59:59Z",
        "max_participant": "50",
        "created_at": "2024-10-28T10:00:00Z",
        "task_image": "https://example.com/image.jpg",
        "action_id": 1,
        "updated_at": "2024-10-28T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 100
    }
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
- Failed to fetch tasks

### Example Request

#### cURL
```bash
# Get first page
curl -X GET http://localhost:8080/api/v1/tasks

# Get specific page
curl -X GET http://localhost:8080/api/v1/tasks?page=2
```

#### JavaScript (Fetch)
```javascript
// Get first page
fetch('http://localhost:8080/api/v1/tasks')
  .then(response => response.json())
  .then(data => console.log(data));

// Get specific page
fetch('http://localhost:8080/api/v1/tasks?page=2')
  .then(response => response.json())
  .then(data => console.log(data));
```

---

## Edit Task

### Endpoint
`PUT /tasks/{id}`

### Authentication
**Required**: No (but should be implemented)

### Path Parameters
- **id**: Task ID (integer)

### Request Body
```json
{
  "title": "string",
  "description": "string",
  "reward_task": 1,
  "reward_usdt": 100.50,
  "due_date": "2024-12-31T23:59:59Z",
  "max_participant": "50",
  "task_image": "https://example.com/image.jpg",
  "action_id": 1
}
```

### Field Notes
- All fields are **optional**
- Only provided fields will be updated
- Empty or zero values will be ignored
- `updated_at` is automatically set to current timestamp

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Task updated successfully",
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

**Possible Causes:**
- Invalid task ID parameter
- Malformed JSON in request body

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
- Task not found
- No fields provided to update
- Database connection error

### Example Request

#### cURL
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Task Title",
    "description": "Updated description",
    "reward_usdt": 150.00
  }'
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/tasks/1', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    title: 'Updated Task Title',
    description: 'Updated description',
    reward_usdt: 150.00
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## Delete Task

### Endpoint
`DELETE /tasks/{id}`

### Authentication
**Required**: No (but should be implemented)

### Path Parameters
- **id**: Task ID (integer)

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "Task deleted successfully",
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

**Cause:** Invalid task ID parameter

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
- Task not found
- Database connection error

### Example Request

#### cURL
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```

#### JavaScript (Fetch)
```javascript
fetch('http://localhost:8080/api/v1/tasks/1', {
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

## Task Model

| Field           | Type      | Description                              |
|-----------------|-----------|------------------------------------------|
| id              | integer   | Unique task identifier                   |
| title           | string    | Task title                               |
| description     | string    | Detailed task description                |
| user_id         | integer   | ID of user who created the task          |
| reward_task     | integer   | Reward ID reference                      |
| reward_usdt     | float     | Reward amount in USDT                    |
| due_date        | timestamp | Task deadline                            |
| max_participant | string    | Maximum number of participants           |
| created_at      | timestamp | Task creation timestamp                  |
| task_image      | string    | URL to task image                        |
| action_id       | integer   | Task action type ID                      |
| updated_at      | timestamp | Last update timestamp                    |

## Notes
- The `user_id` is automatically set from the authenticated user's JWT token in the Create Task endpoint
- Pagination is implemented with a default limit of 10 tasks per page
- The `page` query parameter starts from 1 (not 0)
- Total count in metadata represents the total number of tasks, not total pages
- Edit Task uses partial updates - only send fields you want to change
- Delete operation is permanent and cannot be undone

## Security Considerations
- Create Task endpoint requires JWT authentication
- Edit and Delete endpoints should implement authentication and authorization checks
- Verify that users can only edit/delete their own tasks
- All endpoints should be served over HTTPS in production
- Validate all input data before processing
- Implement rate limiting to prevent abuse

## Pagination Details
- **Default limit**: 10 tasks per page
- **Page numbering**: Starts from 1
- **Offset calculation**: `(page - 1) * limit`
- **Total**: Total number of tasks in database

### Pagination Example
```
Page 1: offset = 0, limit = 10  (tasks 1-10)
Page 2: offset = 10, limit = 10 (tasks 11-20)
Page 3: offset = 20, limit = 10 (tasks 21-30)
```
