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

### Complete API Documentation

For detailed API documentation with examples, request/response formats, and all available endpoints, see **[API.md](./API.md)**.

### Quick Reference

#### Health Check
- `GET /health` - Check API health status

#### Users Management
- `GET /api/users/v1` - Get all users
- `GET /api/users/v1/:id` - Get user by ID
- `POST /api/users/v1` - Create new user
- `PUT /api/users/v1/:id` - Update existing user
- `DELETE /api/users/v1/:id` - Delete user

#### Tickets Management
- `GET /api/tickets/v1` - Get all tickets
- `GET /api/tickets/v1/:id` - Get ticket by ID
- `GET /api/tickets/user/:user_id` - Get tickets by user
- `GET /api/tickets/status/:status` - Get tickets by status
- `GET /api/tickets/category/:category` - Get tickets by category
- `GET /api/tickets/pending` - Get pending tickets
- `GET /api/tickets/search?q=query` - Search tickets
- `POST /api/tickets/v1` - Create new ticket
- `PUT /api/tickets/v1/:id` - Update ticket
- `PATCH /api/tickets/v1/:id/status` - Update ticket status
- `POST /api/tickets/v1/bulk-status` - Bulk update status
- `DELETE /api/tickets/v1/:id` - Delete ticket
- `GET /api/tickets/statistics` - Get ticket statistics

#### Schedule Management

- `GET /api/schedules/tickets` - Get all schedule tickets
- `GET /api/schedules/tickets/:id` - Get schedule ticket by ID
- `GET /api/schedules/tickets/user/:user_id` - Get schedule tickets by user
- `GET /api/schedules/tickets/category/:category` - Get schedule tickets by category
- `POST /api/schedules/tickets` - Create new schedule ticket
- `PUT /api/schedules/tickets/:id` - Update schedule ticket
- `DELETE /api/schedules/tickets/:id` - Delete schedule ticket

#### Schedule Reguler

- `GET /api/schedules/reguler` - Get all regular schedules
- `GET /api/schedules/reguler/:id` - Get regular schedule by ID
- `GET /api/schedules/reguler/user/:user_id` - Get regular schedules by user
- `POST /api/schedules/reguler` - Create new regular schedule
- `PUT /api/schedules/reguler/:id` - Update regular schedule
- `DELETE /api/schedules/reguler/:id` - Delete regular schedule

#### Unblocking (Semester Management)

- `GET /api/unblocking` - Get all unblocking records
- `GET /api/unblocking/:id` - Get unblocking record by ID
- `GET /api/unblocking/user/:user_id` - Get unblocking records by user
- `GET /api/unblocking/semester/:tahun/:semester` - Get unblocking by semester
- `POST /api/unblocking` - Create new unblocking record
- `PUT /api/unblocking/:id` - Update unblocking record
- `DELETE /api/unblocking/:id` - Delete unblocking record

#### Documentation
- `GET /api/docs` - View API documentation (HTML)

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

## �️ Database

The application uses PostgreSQL with GORM as the ORM. Database migrations are managed using golang-migrate.

### Database Schema

#### Users Table
- User authentication via Google OAuth2
- Roles: `user`, `admin`
- Unique constraints on email and google_sub

#### Tickets Table
- Booking requests for rooms and equipment
- Categories: `Kelas`, `Praktikum`, `Skripsi`, `Lainnya`
- Status: `pending`, `accepted`, `rejected`
- Includes contact information and scheduling details

#### Schedule Tables
- `schedule_ticket` - Schedules created from accepted tickets
- `schedule_reguler` - Regular recurring schedules
- `unblocking` - Semester unblocking periods

#### Items Tables
- `items_category` - Equipment/room categories
- `items` - Individual equipment/room items

See [sample_data.sql](../sample_data.sql) for example data.

## 🔧 Message Queue

The application uses RabbitMQ for asynchronous processing:
- Schedule worker processes approved tickets
- Automatically creates schedule entries
- Handles bulk operations efficiently

## 🚀 Future Extensions

Potential improvements and features:

- **Enhanced Authentication** (JWT tokens, refresh tokens)
- **Caching** (Redis for frequently accessed data)
- **Advanced Search** (Full-text search, filters)
- **Swagger Documentation** (Interactive API specs)
- **Testing** (Unit tests, integration tests)
- **Monitoring** (Prometheus, Grafana)
- **Rate Limiting** (API throttling)
- **WebSocket Support** (Real-time notifications)
- **Email Notifications** (Booking confirmations)
- **Calendar Integration** (Export to Google Calendar)

## 📝 Notes

- Database: PostgreSQL with GORM ORM
- User IDs are auto-incremented integers
- Email and google_sub uniqueness is enforced
- All endpoints return JSON responses
- CORS is enabled for all origins (configure for production)
- The server uses Gin's release mode (change to debug for development)
- RabbitMQ is used for asynchronous task processing

## 🔍 Example Curl Commands

```bash
# Get all users
curl http://localhost:8080/api/users/v1

# Create a ticket
curl -X POST http://localhost:8080/api/tickets/v1 \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "request_data": {
      "name": "Peminjaman Lab",
      "desc": "Lab untuk praktikum",
      "category": "Praktikum",
      "requestDate": "2025-10-20T08:00:00+07:00",
      "email": "user@example.com",
      "phone": "081234567890",
      "pic": "John Doe"
    }
  }'

# Update ticket status
curl -X PATCH http://localhost:8080/api/tickets/v1/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "accepted"}'

# Search tickets
curl "http://localhost:8080/api/tickets/search?q=lab"
```

## 📚 Additional Resources

- **[API Documentation](./API.md)** - Complete API reference
- **[Sample Data](../sample_data.sql)** - Database sample data
- **[Migrations](../migrations/)** - Database migration files
- **[Makefile Reference](../README.md)** - Development commands
