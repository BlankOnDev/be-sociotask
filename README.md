# SocioTask - Task Marketplace Backend

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12.4-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A robust RESTful API backend for a task marketplace platform where users can post tasks for others to complete and earn rewards. Built with Go, PostgreSQL, and modern web technologies.

## 📋 Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [License](#license)

## 🎯 Overview

SocioTask is a task marketplace backend system that enables users to create, manage, and complete tasks in exchange for rewards. The platform supports multiple authentication methods including traditional email/password, Google OAuth, and Twitter OAuth, making it accessible to a wide range of users.

The system handles:
- **Task Management**: Users can create tasks with descriptions, deadlines, participant limits, and reward amounts
- **Reward System**: Automated reward tracking and distribution for completed tasks
- **User Authentication**: Multiple OAuth providers and JWT-based authentication
- **Task Actions**: Categorized task types (type_1, type_2, type_3) for different task categories
- **User Profiles**: User management with profiles, wallet addresses, and task history

## ✨ Key Features

### User Management
- ✅ User registration and authentication
- ✅ Google OAuth 2.0 integration
- ✅ Twitter OAuth integration
- ✅ JWT-based session management
- ✅ User profiles with bio and wallet addresses
- ✅ Android-specific Google login support

### Task System
- ✅ Create tasks with detailed specifications
- ✅ Set task rewards (USDT-based)
- ✅ Define participant limits
- ✅ Set deadlines for task completion
- ✅ Task categorization through action types
- ✅ Task status tracking
- ✅ Update and delete tasks
- ✅ Filter tasks by various criteria

### Reward System
- ✅ Automated reward tracking
- ✅ Task completion verification
- ✅ Reward distribution management
- ✅ User reward history

### Security
- ✅ Password hashing using bcrypt
- ✅ JWT token authentication
- ✅ Protected API endpoints
- ✅ CORS configuration
- ✅ Request validation and sanitization

## 🛠 Tech Stack

### Backend Framework
- **Language**: Go 1.25.0
- **Router**: Chi v5 - Lightweight and fast HTTP router
- **Authentication**: JWT (golang-jwt/jwt v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)

### Database
- **Database**: PostgreSQL 12.4
- **Driver**: pgx/v5 - High-performance PostgreSQL driver
- **Migration**: Goose v3 - Database migration tool

### OAuth Integration
- **Google OAuth**: Google API Client
- **Twitter OAuth**: OAuth2 implementation

### Testing
- **Framework**: Testify - Testing toolkit with assertions
- **Strategy**: Unit tests for store layer
- **Database**: Docker-based PostgreSQL for testing

### DevOps
- **Containerization**: Docker & Docker Compose
- **Server**: Custom HTTP server with timeouts

## 🏗 Architecture

### Project Structure
```
be-sociotask/
├── internal/
│   ├── api/              # HTTP handlers
│   │   ├── task_handler.go
│   │   ├── user_handler.go
│   │   ├── auth_handler.go
│   │   ├── action_handler.go
│   │   ├── reward_handler.go
│   │   └── rewards_handler.go
│   ├── app/              # Application initialization
│   │   └── app.go
│   ├── auth/             # Authentication logic
│   │   ├── jwt.go
│   │   └── google/
│   ├── middleware/       # HTTP middleware
│   │   └── middleware.go
│   ├── pages/            # HTML pages
│   │   └── html.go
│   ├── routes/           # Route definitions
│   │   └── routes.go
│   ├── store/            # Data access layer
│   │   ├── database.go
│   │   ├── task_store.go
│   │   ├── user_store.go
│   │   ├── rewards_store.go
│   │   ├── task_action_store.go
│   │   └── *_test.go
│   └── utils/            # Utility functions
│       └── utils.go
├── migrations/           # Database migrations
│   ├── 00001_users.sql
│   ├── 00002_tasks.sql
│   ├── 00003_*.sql
│   └── ...
├── docs/                 # API documentation
│   ├── task-api.md
│   ├── user-current.md
│   ├── rewards-api.md
│   ├── task-action-api.md
│   ├── task-reward-api.md
│   ├── login-api.md
│   ├── register-api.md
│   ├── oauth-google-api.md
│   └── oauth-twitter-api.md
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
└── docker-compose.yml    # Docker configuration
```

### Database Schema

**Users Table**
- User credentials and profile information
- OAuth identities (X/Twitter ID)
- Wallet addresses for rewards
- Full name and bio

**Tasks Table**
- Task details (title, description)
- Creator reference (user_id)
- Reward information (USDT amount)
- Deadline and participant limits
- Task status tracking
- Associated action type

**Task Actions Table**
- Action type classification (type_1, type_2, type_3)
- Action name and description
- Used to categorize tasks

**Task Rewards Table**
- Reward definitions
- Reference data for task rewards

**Rewards Table**
- Records of task completions
- Links users to completed tasks
- Timestamp tracking

### API Design Pattern

The project follows a clean architecture pattern:

1. **Handlers** (`internal/api/`) - HTTP request/response handling
2. **Store** (`internal/store/`) - Database operations and business logic
3. **Middleware** (`internal/middleware/`) - Authentication and authorization
4. **Routes** (`internal/routes/`) - API route definitions

## 🚀 Getting Started

### Prerequisites

- **Go**: Version 1.25.0 or higher
- **PostgreSQL**: Version 12.4 or higher (or use Docker)
- **Docker & Docker Compose**: (Optional, for containerized database)
- **Git**: For cloning the repository

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/BlankOnDev/be-sociotask.git
   cd be-sociotask
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL**

   Using Docker (recommended for development):
   ```bash
   docker-compose up -d
   ```

   This will start a PostgreSQL instance on port 5433.

   Or install PostgreSQL manually and create a database:
   ```sql
   CREATE DATABASE socialtask;
   ```

### Configuration

1. **Create environment file**
   ```bash
   cp .env.example .env
   ```

2. **Configure environment variables**

   Edit `.env` with your settings:

   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=socialtask
   DB_PORT=5432

   # JWT Secret
   JWT_SECRET=your_super_secret_jwt_key

   # Twitter OAuth
   TWITTER_CLIENT_ID=your_twitter_client_id
   TWITTER_CLIENT_SECRET=your_twitter_client_secret
   TWITTER_REDIRECT_URL=http://localhost:8080/login/twitter/callback

   # Google OAuth
   Google_Client_ID_Web=your_google_client_id
   Google_Client_Secret_Web=your_google_client_secret
   Google_Redirect_URL=http://localhost:8080/login/google/callback
   ```

3. **Load environment variables** (Linux/Mac)
   ```bash
   set -a
   source .env
   set +a
   ```

   Or for Windows PowerShell:
   ```powershell
   Get-Content .env | ForEach-Object {
       if ($_ -match '^([^=]+)=(.*)$') {
           [System.Environment]::SetEnvironmentVariable($matches[1], $matches[2])
       }
   }
   ```

### Database Setup

1. **Run migrations**

   The application uses Goose for database migrations. Migrations are automatically applied when the application starts.

   To run migrations manually:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   goose -dir migrations postgres "host=localhost user=postgres password=postgres dbname=socialtask port=5432 sslmode=disable" up
   ```

2. **Verify migrations**
   ```bash
   goose -dir migrations postgres "host=localhost user=postgres password=postgres dbname=socialtask port=5432 sslmode=disable" status
   ```

### Running the Application

1. **Start the server**
   ```bash
   go run main.go
   ```

   Or specify a custom port:
   ```bash
   go run main.go -port=3000
   ```

2. **Verify the server is running**
   ```bash
   curl http://localhost:8080/health
   ```

   You should receive a health check response.

3. **Access the application**
   - API Base URL: `http://localhost:8080`
   - Home Page: `http://localhost:8080/`
   - Health Check: `http://localhost:8080/health`

## 📚 API Documentation

The API provides comprehensive endpoints for managing tasks, users, rewards, and authentication.

### Base URL
```
http://localhost:8080
```

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Main Endpoints

#### Authentication
- `POST /register` - Register a new user
- `POST /login` - Login with email/password
- `GET /login/google` - Initiate Google OAuth
- `GET /login/google/callback` - Google OAuth callback
- `POST /login/google/android` - Android Google login
- `GET /login/twitter` - Initiate Twitter OAuth
- `GET /login/twitter/callback` - Twitter OAuth callback

#### Users
- `GET /users/current` - Get current authenticated user (Protected)
- `GET /users/{id}/tasks` - Get tasks created by a specific user

#### Tasks
- `GET /tasks` - Get all tasks
- `GET /tasks/{id}` - Get task by ID
- `POST /tasks` - Create a new task (Protected)
- `PUT /tasks/{id}` - Update a task (Protected)
- `DELETE /tasks/{id}` - Delete a task (Protected)

#### Task Actions
- `GET /actions` - Get all task actions
- `GET /actions/{id}` - Get task action by ID

#### Task Rewards
- `GET /reward` - Get all rewards
- `GET /reward/{id}` - Get reward by ID

#### Rewards (Completion Tracking)
- `POST /rewards` - Record task completion (Protected)

### Detailed API Documentation

For detailed request/response examples, validation rules, and error codes, see:
- [Task API Documentation](docs/task-api.md)
- [User API Documentation](docs/user-current.md)
- [Rewards API Documentation](docs/rewards-api.md)
- [Task Actions API Documentation](docs/task-action-api.md)
- [Task Rewards API Documentation](docs/task-reward-api.md)
- [Login API Documentation](docs/login-api.md)
- [Register API Documentation](docs/register-api.md)
- [Google OAuth Documentation](docs/oauth-google-api.md)
- [Twitter OAuth Documentation](docs/oauth-twitter-api.md)

## 🧪 Testing

The project includes comprehensive unit tests for the store layer.

### Running Tests

1. **Start the test database**
   ```bash
   docker-compose up -d
   ```

2. **Run all tests**
   ```bash
   go test ./...
   ```

3. **Run tests with coverage**
   ```bash
   go test -cover ./...
   ```

4. **Run tests with verbose output**
   ```bash
   go test -v ./...
   ```

5. **Run specific test**
   ```bash
   go test -v ./internal/store -run TestCreateTask
   ```

### Test Structure

Tests are located in `internal/store/*_test.go` files:
- `task_store_test.go` - Task CRUD operations
- `user_store_test.go` - User management operations
- `rewards_store_test.go` - Reward tracking operations
- `task_reward_store_test.go` - Task reward operations

Each test:
- Sets up a clean test database
- Runs migrations
- Executes test cases
- Cleans up after completion

## 🚢 Deployment

### Production Considerations

1. **Environment Variables**
   - Use secure, randomly generated JWT secrets
   - Configure production database credentials
   - Set up production OAuth redirect URLs

2. **Database**
   - Use a managed PostgreSQL service (AWS RDS, Google Cloud SQL, etc.)
   - Enable SSL connections (`sslmode=require`)
   - Set up automated backups
   - Configure connection pooling

3. **Server Configuration**
   - Use a reverse proxy (Nginx, Caddy)
   - Enable HTTPS/TLS
   - Configure rate limiting
   - Set up logging and monitoring

4. **Build for Production**
   ```bash
   go build -o sociotask-server main.go
   ```

5. **Run in Production**
   ```bash
   ./sociotask-server -port=8080
   ```

### Docker Deployment

Create a `Dockerfile`:
```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./server"]
```

Build and run:
```bash
docker build -t sociotask-backend .
docker run -p 8080:8080 --env-file .env sociotask-backend
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Harun Darat**

- GitHub: [@BlankOnDev](https://github.com/BlankOnDev)
