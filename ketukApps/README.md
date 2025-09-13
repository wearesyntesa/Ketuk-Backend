# KetukApps API Server

A clean, simple REST API server built with Go and Gin framework following clean architecture principles.

## 🏗️ Architecture

The project follows a clean architecture pattern with separated concerns:

```
ketukApps/
├── main.go                    # Application entry point
├── config/
│   └── config.go             # Configuration management
├── internal/
│   ├── handlers/             # HTTP request handlers
│   ├── services/             # Business logic
│   ├── models/               # Data structures
│   └── middleware/           # HTTP middleware
├── .env                      # Environment variables
├── go.mod                    # Go module file
├── go.sum                    # Go dependencies checksum
├── test_api.sh              # API testing script
└── README.md                # This file
```

## 🚀 Features

- **Clean Architecture**: Separation of concerns with clear layers
- **RESTful API**: Standard HTTP methods and status codes
- **Middleware Support**: CORS, Logging, Error handling, Recovery
- **Environment Configuration**: Configurable via environment variables
- **In-Memory Storage**: Simple user management (can be extended to database)
- **API Documentation**: Built-in documentation endpoint
- **Health Checks**: Monitoring endpoint for service health

## 📋 API Endpoints

### Health Check
- `GET /health` - Check API health status

### Users Management
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update existing user
- `DELETE /api/users/:id` - Delete user

### Documentation
- `GET /api/docs` - View API documentation

## 🔧 Configuration

Environment variables (can be set in `.env` file):

```env
PORT=8080           # Server port (default: 8080)
HOST=localhost      # Server host (default: localhost)
LOG_LEVEL=info      # Log level (default: info)
```

## 🏃‍♂️ Running the Server

### Prerequisites
- Go 1.21 or higher
- Git (for dependencies)

### Install Dependencies
```bash
go mod tidy
```

### Run the Server
```bash
# Using go run
go run main.go

# Or build and run
go build -o ketukApps
./ketukApps
```

The server will start on `http://localhost:8080` (or configured PORT).

## 🧪 Testing the API

### Using the Test Script
```bash
# Make sure server is running on port 8082
PORT=8082 go run main.go &

# Run the test script
./test_api.sh
```

### Manual Testing with curl

1. **Health Check:**
```bash
curl http://localhost:8080/health
```

2. **Create User:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

3. **Get All Users:**
```bash
curl http://localhost:8080/api/users
```

4. **Get User by ID:**
```bash
curl http://localhost:8080/api/users/1
```

5. **Update User:**
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated"}'
```

6. **Delete User:**
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## 📁 Project Structure Explained

### `main.go`
- Application entry point
- Server initialization and routing setup
- Dependency injection

### `config/`
- Configuration management
- Environment variable handling

### `internal/handlers/`
- HTTP request handlers
- Request/response processing
- Input validation

### `internal/services/`
- Business logic implementation
- Data processing
- Core application functionality

### `internal/models/`
- Data structures and models
- Request/response DTOs
- Domain entities

### `internal/middleware/`
- HTTP middleware functions
- CORS, logging, error handling
- Request/response processing pipeline

## 🔍 Example Responses

### Success Response
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-09-12T04:00:00Z",
    "updated_at": "2025-09-12T04:00:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "User not found",
  "error": "user not found"
}
```

## 🚀 Future Extensions

This is a foundational API server that can be extended with:

- **Database Integration** (PostgreSQL, MySQL, MongoDB)
- **Authentication & Authorization** (JWT, OAuth)
- **Caching** (Redis)
- **Validation** (Enhanced input validation)
- **Swagger Documentation** (API specs)
- **Testing** (Unit tests, integration tests)
- **Docker** (Containerization)
- **Monitoring** (Metrics, logging)
- **Rate Limiting** (API throttling)
- **Message Queues** (RabbitMQ, Kafka)

## 📝 Notes

- This implementation uses in-memory storage for simplicity
- User IDs are auto-incremented integers
- Email uniqueness is enforced
- All endpoints return JSON responses
- CORS is enabled for all origins (configure for production)
- The server uses Gin's release mode (change to debug for development)
