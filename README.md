# User Management HTTP Server

A production-ready RESTful API for managing user data, built with Go, Gorilla/Mux, and BoltDB. Implements clean architecture with separation of concerns, proper error handling, and context-aware operations.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [API Reference](#api-reference)
- [Getting Started](#getting-started)
- [Design Decisions](#design-decisions)
- [Testing Strategy](#testing-strategy)
- [Future Improvements](#future-improvements)

## Overview

This project demonstrates a modular HTTP server implementing full CRUD operations on a user resource. Key highlights:

- Interface-driven design — Loose coupling between HTTP layer and persistence via database.Database contract
- Proper error handling — Error wrapping with %w, no silent failures, graceful degradation on data corruption
- Context propagation — Timeout and cancellation support throughout the stack
- Resource cleanup — Deferred database closes, graceful shutdown patterns

Built as a showcase of modern Go practices including idiomatic error patterns, proper HTTP semantics, and clean separation of concerns.

## Architecture

### Layered Design

Response:

Response:

Response:

main.go (Application Entry Point)

    Initializes BoltDB
    Configures router & server
    Handles graceful shutdown | +---> server_operations/server.go --+ | - HandleIndex | | - HandleCreateUser |---> database/bolt/bolt.go | - HandleUser (CRUD) | - Create(), Get(), Update(), Delete() +-----------------------------------+ ^ | database/database.go - Database interface - User struct


### Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| main.go | Bootstrap, dependency injection, server lifecycle |
| server.go | HTTP handlers, request/response translation, validation |
| bolt.go | Persistence implementation, transaction management |
| database.go | Shared contracts, domain models (User struct) |

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Router | github.com/gorilla/mux | Request routing and path variables |
| Database | github.com/boltdb/bolt | Embedded key-value storage |
| Language | Go 1.20+ | Type safety, concurrency, performance |
| Serialization | encoding/json | Request/response payloads |
| Context | context package | Timeout/cancellation propagation |

## Features

| Feature | Implementation Detail |
|---------|----------------------|
| RESTful API | Standard HTTP methods (GET, POST, PUT, PATCH, DELETE) |
| JSON Request/Response | Content-Type enforcement, automatic marshaling |
| Input Validation | Name required, existence checks, type validation |
| Error Wrapping | %w for traceable error chains |
| Graceful Degradation | Corrupt data logged but does not crash application |
| Context Timeouts | 30s database operation timeout, 10s HTTP timeout |
| Deferred Cleanup | Database closes reliably on exit |
| Idempotent Operations | CREATE and DELETE handle duplicates safely |

## API Reference

### Endpoints

| Method | Endpoint | Description | Success | Error Codes |
|--------|----------|-------------|---------|-------------|
| GET | / | Serve dashboard HTML | 200 OK | N/A |
| GET | /user/{name} | Retrieve user by name | 200 OK + JSON | 404 Not Found |
| POST | /user/create | Create new user | 201 Created + JSON | 400 Bad Request, 415 Unsupported Media |
| PUT | /user/create | Alias for POST | 201 Created + JSON | Same as POST |
| PATCH | /user/{name} | Partially update user | 200 OK + JSON | 400, 404, 415 |
| DELETE | /user/{name} | Remove user | 204 No Content | 500 Internal Server |

### Request/Response Examples

**Create User**

Request:

bash curl -X POST http://localhost:9090/user/create
-H "Content-Type: application/json"
-d '{"name": "alice", "email": "alice@example.com", "age": 30}'

Response:

HTTP/1.1 201 Created Content-Type: application/json

{"name":"alice","email":"alice@example.com","age":30}

**Get User**

Request:

bash curl http://localhost:9090/user/alice

Response:

HTTP/1.1 200 OK Content-Type: application/json

{"name":"alice","email":"alice@example.com","age":30}

**Update User (Partial)**

Request:

bash curl -X PATCH http://localhost:9090/user/alice
-H "Content-Type: application/json"
-d '{"age": 31}'

Response:

HTTP/1.1 200 OK Content-Type: application/json

{"name":"alice","email":"alice@example.com","age":31}

**Delete User**

Request:

bash curl -X DELETE http://localhost:9090/user/alice

Response:

HTTP/1.1 204 No Content

### Error Responses

| Status | Body | Meaning |
|--------|------|---------|
| 400 Bad Request | "user already exists: alice" | Duplicate or invalid input |
| 404 Not Found | (empty) | User does not exist |
| 405 Method Not Allowed | "method not allowed" | Invalid HTTP verb |
| 415 Unsupported Media | (empty) | Content-Type not application/json |
| 500 Internal Server | (empty) | Unexpected error (check logs) |

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Access to module dependencies (go mod tidy)

### Installation

1. Clone repository
   ```bash
   git clone <repository-url>
   cd Http-Server
   ```

2. Download dependencies
   ```bash
   go mod tidy
   ```

3. Build binary
   ```bash
   go build -o server ./cmd
   ```

4. Run server
   ```bash
   ./server
   ```

Or run directly:

bash go run .

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | :9090 | Bind address for HTTP server |
| DATA_DIR | ./data | Directory for BoltDB file storage |

Example:

bash PORT=:8080 DATA_DIR=/tmp/server go run .

### Project Structure

Http-Server/ ├── main.go # Application entry point ├── server.go # HTTP handlers (server_operations/) ├── database.go # Domain model + interface (database/) ├── bolt.go # BoltDB implementation (bolt/) ├── index.go # Dashboard HTML template (pages/) └── go.mod # Module definition

## Design Decisions

### Why Interface-Driven Architecture?

The Database interface in database.go decouples HTTP handlers from the persistence layer:

go type Database interface { Create(ctx context.Context, user User) error Get(ctx context.Context, name string) *User Update(ctx context.Context, user User) (*User, error) Delete(ctx context.Context, name string) error }

This enables:
1. Testability — Mock implementations for unit tests
2. Flexibility — Swap BoltDB for PostgreSQL without changing handlers
3. Single Responsibility — Each layer has one clear purpose

### Why BoltDB?

- Zero dependency — No separate database server required
- ACID compliance — Transactional guarantees via db.Update()
- Simple API — Ideal for demonstrating Go database patterns
- Embedded — Single binary deployment, perfect for demos

### Why No External Validation Library?

Kept dependencies minimal. Validation logic lives in handlers:
- Check empty name before DB operations
- Verify Content-Type before parsing
- Fail fast with specific HTTP status codes

### Error Handling Strategy

1. Wrap errors with context — Use fmt.Errorf("%w") for traceability
2. Log warnings, do not crash — Corrupt data returns nil instead of panicking
3. Return appropriate status codes — Map internal errors to HTTP semantics

## Testing Strategy

### Unit Tests Needed

| File | Test Cases |
|------|------------|
| server.go | Invalid JSON handling, missing fields, method validation |
| bolt.go | Create duplicate rejection, Get non-existent, Update partial merge |
| database.go | N/A — pure interface definition |

### Integration Tests

| Scenario | Expected Result |
|----------|-----------------|
| POST /user/create with existing name | 400 Bad Request |
| PATCH /user/{name} with empty name field | 400 Bad Request |
| DELETE /user/{name} then GET same name | 404 Not Found |
| GET / without authentication | 200 OK (public endpoint) |

### Test Coverage Goals

- Handlers: >= 80% (valid + invalid inputs)
- Persistence: >= 90% (transaction paths, error conditions)

## Future Improvements

| Area | Enhancement | Priority |
|------|-------------|----------|
| Authentication | JWT or API key middleware | High |
| Pagination | List() endpoint with offset/limit | High |
| Database Migration | Replace BoltDB with PostgreSQL for scalability | Medium |
| Health Checks | /health endpoint exposing DB stats | Medium |
| Logging | Structured logging with log/slog | Low |
| OpenAPI Spec | Generate Swagger documentation | Low |
| Rate Limiting | Middleware to prevent abuse | Low |

## Lessons Learned

1. Interface-first design pays dividends when adding new features or storage backends
2. Error wrapping (%w) makes debugging significantly easier in production
3. Graceful shutdown prevents data loss and improves UX during deployments
4. Minimal dependencies reduce build complexity and attack surface
5. Context propagation is essential for timeouts across network calls

## License

MIT License — see LICENSE file for details.

Built with ❤️ using Go
