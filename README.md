# Clean Architecture Go Project

A RESTful API built with Go following clean architecture principles, featuring a user service with MSSQL database integration.

## Project Structure

```
.
├── entity/                 # Domain Entities (Pure Logic)
├── usecase/               # UseCase Interfaces
├── internal/
│   ├── user/
│   │   ├── service.go     # UseCase Implementation
│   │   └── handler.go     # HTTP Controllers
│   └── middleware/        # Logging, Tracing Middleware
├── infra/mssql/          # Database Adapter for MSSQL
├── api/                  # Request/Response DTOs
├── pkg/                  # Shared Utilities
├── routers/              # Route Setup
├── cmd/                  # Application Entry (main.go)
├── go.mod               # Go Module Definition
├── go.sum               # Go Module Checksums
└── README.md            # Documentation
```

## Architecture Overview

- **Entity Layer**: Domain models and business logic entities
- **UseCase Layer**: Business logic interfaces and rules
- **Interface Adapters**: Controllers, handlers, and DTOs
- **Infrastructure Layer**: Database implementations, external services
- **Utilities**: Shared functions and error handling

## Multi-Service Architecture Pattern

This project demonstrates **Dependency Injection** pattern for composing multiple services (e.g., User + Address):

### How It Works

1. **Separate Repositories**: Each service has its own repository interface and MSSQL implementation
   - `UserRepository` for user data
   - `AddressRepository` for address data

2. **Service Composition**: Services are injected with dependencies from other services
   ```go
   // User handler receives both user usecase and address repository
   userHandler := user.NewHandlerWithAddress(userUsecase, addressRepo)
   ```

3. **Data Aggregation**: When fetching a user, the handler automatically includes related addresses
   ```go
   // GET /api/users/:id returns:
   {
     "id": 1,
     "name": "John Doe",
     "email": "john@example.com",
     "phone": "1234567890",
     "addresses": [
       {
         "id": 1,
         "user_id": 1,
         "street": "123 Main St",
         "city": "New York",
         "state": "NY",
         "country": "USA",
         "zip_code": "10001"
       }
     ]
   }
   ```

### Adding New Services

To add another service (e.g., "Order"):

1. Create entity: `entity/order.go`
2. Create usecase: `usecase/order.go`
3. Create service: `internal/order/service.go`
4. Create handler: `internal/order/handler.go`
5. Create repository: `infra/mssql/order_repository.go`
6. Create DTOs: `api/order_dto.go`
7. Update `main.go` to initialize and inject dependencies
8. Update `routers/router.go` to register routes

The pattern ensures:
- **Dependency Injection**: Services depend on interfaces, not implementations
- **Loose Coupling**: Services don't directly depend on each other's implementations
- **Testability**: Easy to mock dependencies for testing
- **Scalability**: New services can be added without modifying existing ones


## Prerequisites

- Go 1.21 or higher
- MSSQL Server
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd clean
```

2. Install dependencies:
```bash
go mod download
```

3. Set up MSSQL database:
```sql
CREATE DATABASE clean_db;

USE clean_db;

CREATE TABLE users (
    id INT PRIMARY KEY IDENTITY(1,1),
    name NVARCHAR(255) NOT NULL,
    email NVARCHAR(255) NOT NULL,
    phone NVARCHAR(20)
);

CREATE TABLE addresses (
    id INT PRIMARY KEY IDENTITY(1,1),
    user_id INT NOT NULL,
    street NVARCHAR(255) NOT NULL,
    city NVARCHAR(100) NOT NULL,
    state NVARCHAR(100),
    country NVARCHAR(100) NOT NULL,
    zip_code NVARCHAR(20),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_addresses_user_id ON addresses(user_id);
```

## Configuration

Set environment variables:
```bash
export DATABASE_URL="server=localhost;user id=sa;password=YourPassword;database=clean_db;port=1433"
export PORT=8080
```

## Running the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### User Management

- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID (includes associated addresses)
- `POST /api/users` - Create new user
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Address Management

- `GET /api/addresses` - List all addresses
- `GET /api/addresses/:id` - Get address by ID
- `GET /api/users/:user_id/addresses` - Get all addresses for a specific user
- `POST /api/addresses` - Create new address
- `PUT /api/addresses/:id` - Update address
- `DELETE /api/addresses/:id` - Delete address

### Health Check

- `GET /health` - Health check endpoint

### Example Usage

**Create a user:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "1234567890"
  }'
```

**Get user with addresses:**
```bash
curl http://localhost:8080/api/users/1
```

Response includes user data + associated addresses:
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "addresses": [
    {
      "id": 1,
      "user_id": 1,
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "country": "USA",
      "zip_code": "10001"
    }
  ]
}
```

**Create address for user:**
```bash
curl -X POST http://localhost:8080/api/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "street": "456 Oak Ave",
    "city": "Los Angeles",
    "state": "CA",
    "country": "USA",
    "zip_code": "90001"
  }'
```

**Get all addresses for a user:**
```bash
curl http://localhost:8080/api/users/1/addresses
```

**Update address:**
```bash
curl -X PUT http://localhost:8080/api/addresses/1 \
  -H "Content-Type: application/json" \
  -d '{
    "street": "789 Elm St",
    "city": "San Francisco",
    "state": "CA",
    "country": "USA",
    "zip_code": "94102"
  }'
```

**Delete address:**
```bash
curl -X DELETE http://localhost:8080/api/addresses/1
```

**Update user:**
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "phone": "0987654321"
  }'
```

**Delete user:**
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## Building the Application

```bash
go build -o clean cmd/main.go
./clean
```

## Testing

```bash
go test ./...
```

## Generate Mocks (Mockery)

This project uses `mockery` to generate mocks for unit tests.

Generate mocks for `IService` interfaces:

```bash
go run github.com/vektra/mockery/v2@latest --dir internal/user --name IService --output internal/user/mocks --outpkg mocks --filename user_service_mock.go --structname UserServiceMock
go run github.com/vektra/mockery/v2@latest --dir internal/address --name IService --output internal/address/mocks --outpkg mocks --filename address_service_mock.go --structname AddressServiceMock
```

Generate mocks for `IRepository` interfaces:

```bash
go run github.com/vektra/mockery/v2@latest --dir infra/mssql/user --name IRepository --output internal/user/mocks --outpkg mocks --filename user_repository_mock.go --structname UserRepositoryMock
go run github.com/vektra/mockery/v2@latest --dir infra/mssql/address --name IRepository --output internal/address/mocks --outpkg mocks --filename address_repository_mock.go --structname AddressRepositoryMock
```

After generating mocks:

```bash
go mod tidy
go test ./internal/user ./internal/address
```

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/microsoft/go-mssqldb` - MSSQL driver

## License

MIT
