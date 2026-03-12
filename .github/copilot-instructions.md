# Clean Architecture Go Project - Copilot Instructions

## Project Overview

This is a clean architecture Go project implementing a RESTful API with user service management and MSSQL database integration.

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: MSSQL
- **Architecture Pattern**: Clean Architecture with Dependency Injection

## Project Structure Conventions

- **entity/**: Domain models with pure business logic
- **usecase/**: Interface definitions for business logic contracts
- **internal/user/**: Service implementations and HTTP handlers
- **internal/middleware/**: Cross-cutting concerns (logging, CORS, etc.)
- **infra/mssql/**: Database adapter implementations
- **api/**: Data Transfer Objects (DTOs) for API communication
- **pkg/**: Shared utilities and error handling
- **routers/**: Route registration and middleware setup
- **cmd/**: Application entry points

## Key Design Principles

1. **Dependency Injection**: Services depend on interfaces, not concrete implementations
2. **Separation of Concerns**: Each layer has a specific responsibility
3. **Testability**: Interfaces enable easy mocking and testing
4. **Database Independence**: Repository pattern abstracts database details

## Running the Project

1. Set up MSSQL database with the users table
2. Configure DATABASE_URL environment variable
3. Run: `go run cmd/main.go`
4. Server starts on port 8080 (configurable via PORT env var)

## Adding New Features

When adding new domain entities/features:

1. Create entity in `entity/` directory
2. Define usecase interface in `usecase/` directory
3. Implement service in `internal/{feature}/service.go`
4. Create handler in `internal/{feature}/handler.go`
5. Implement repository in `infra/mssql/{feature}_repository.go`
6. Create DTOs in `api/{feature}_dto.go`
7. Register routes in `routers/router.go`

## Code Style

- Follow Go conventions (camelCase for variables, PascalCase for exported symbols)
- Use interfaces for abstraction and dependency injection
- Keep functions small and focused
- Add error handling for all I/O operations
- Use named return values with meaningful names
