# Inventhier

Production-grade inventory management system demonstrating enterprise Go development practices: hexagonal architecture, multi-database integration, comprehensive testing, and clean code patterns.

## What This Project Shows

- **Systems Design**: Multi-protocol integration (gRPC vs REST), database selection by use case, scalability patterns
- **Production Go**: Proper context usage, structured logging, transaction management, error handling
- **Testing Culture**: Multi-layer testing (50+ tests), mock generation, performance benchmarks
- **Clean Architecture**: Port-driven design, dependency injection, separation of concerns

## Architecture

**Hexagonal (Ports & Adapters)**: Business logic is isolated in `core/` domain, completely independent of infrastructure. All external dependencies are interfaces.

```
HTTP Request → [Router] → [Middleware] → [Handler]
                                            ↓
                                     [Service] (Core Logic)
                                            ↓
                    ┌─────────────┬─────────┴──────────┐
                    ↓             ↓                     ↓
            [gRPC Client]   [REST Client]      [Repositories]
                    ↓             ↓                     ↓
            [User Service] [Payment API]   [SQL/Redis/MongoDB]
```

Key architectural files:
- **`internal/bootstrap/app.go`** - DI container (wires 6+ dependencies without frameworks)
- **`internal/core/services/v1/product/`** - Business logic (no framework coupling)
- **`internal/adapters/`** - Infrastructure adapters (swappable implementations)

## Code Quality Highlights

### 1. Dependency Injection Without Frameworks
Pure Go DI container in `bootstrap/app.go`. Explicitly wires repositories, clients, and services with compile-time type safety. No reflection, no service locators.

### 2. Advanced SQL Patterns
Pessimistic locking with `FOR UPDATE` prevents race conditions on concurrent updates:
```go
tx, err := db.BeginTx(ctx, nil)
defer tx.Rollback() // Safety rollback
queryLock := `SELECT stock FROM products WHERE sku = ? FOR UPDATE`
// ... execute within transaction
tx.Commit()
```

### 3. Multi-Protocol Integration
- **gRPC** (`internal/adapters/client/grpc/v1/user/`) - Type-safe, high-performance inter-service communication
- **REST** (`internal/adapters/client/rest/v1/payment/`) - Flexible external API integration
- Both implement port interfaces; service layer doesn't care which protocol

### 4. Request Tracing
Auto-generates request IDs, propagates via context, includes in all logs and HTTP headers. Enables tracing a request through multiple services.

### 5. Structured Logging
Uses **zerolog** everywhere (no `fmt.Printf` in production code). JSON output integrates with ELK/Datadog/CloudWatch for production debugging.

### 6. Multi-Layer Testing
- **Service tests** (8+): Business logic with mocked dependencies
- **Repository tests** (20+): Database operations with `go-sqlmock` (no real DB needed)
- **Client tests**: In-process gRPC mocking with `bufconn`
- **50+ total test cases**, coverage reporting, performance benchmarks

## Tech Stack

- **Language**: Go 1.25.0+
- **Router**: Gorilla Mux
- **Databases**: MySQL (SQL), MongoDB (audit logs), Redis (cache)
- **RPC**: gRPC, REST APIs
- **Testing**: mockery (auto mock generation), testify, go-sqlmock, bufconn
- **Logging**: zerolog
- **Quality**: golangci-lint, coverage reporting

## Quick Start

```bash
# Clone & setup
git clone https://github.com/bagus-aulia/inventhier
cd inventhier
go mod download

# Start services
docker-compose up -d

# Configure environment
cp .env.example .env

# Run tests
make test

# Start server
go run cmd/api/main.go
```

## Design Decisions

| Decision | Why | Trade-off |
|----------|-----|-----------|
| **Hexagonal Architecture** | Keep business logic independent of frameworks and databases | Slightly more code structure upfront |
| **Pessimistic Locking** | Serialize concurrent updates, prevent lost updates | Higher latency but guaranteed consistency |
| **Multi-Database** | Right tool for each job (ACID for transactions, flexible schema for logs) | Operational complexity |
| **gRPC for internal, REST for external** | gRPC is 2-3x faster for services; REST is universal for third-party APIs | Learning curve for gRPC |
| **Request ID Tracing** | Trace requests across services and logs for production debugging | Small overhead |

## Key Files for Code Review

| File | What It Shows |
|------|---------------|
| `internal/bootstrap/app.go` | Sophisticated DI without external frameworks |
| `internal/core/services/v1/product/product.go` | Clean business logic (no framework coupling) |
| `internal/adapters/repository/sql/product/product.go` | Production SQL patterns (transactions, locking, audit trails) |
| `internal/adapters/client/grpc/v1/user/user.go` | Adapter pattern for external service integration |
| `internal/adapters/middleware/request_id.go` | Proper Go context usage and request tracing |
| `internal/core/services/v1/product/product_test.go` | Multi-layer testing with mocks |

## Testing & Quality

```bash
make test              # Run all tests with coverage
make test-short        # Quick tests (skip long-running)
make test-coverage     # HTML coverage report
make bench             # Performance benchmarks
make lint              # Run linter
make quality           # Full quality check (lint + test + coverage)
```

**Coverage**: 50+ test cases covering service, repository, client, and middleware layers. Tests validate both success and error paths.

## External Package Integration

Integrates `github.com/bagus-aulia/go-tools` for:
- Structured logging helpers with zerolog
- Reusable utility patterns
- Custom middleware components

## API Endpoints

### Stock In
```
POST /api/v1/stock-in
{
  "staff_uuid": "...",
  "product_sku": "SKU-001",
  "product_qty": 10
}
```

### Get Product
```
GET /api/v1/products/{sku}
```
- Cached in Redis
- Falls back to SQL on cache miss
- Auto-invalidates cache on stock updates

## Production Patterns Demonstrated

- **Connection pooling** - Database connection management
- **Audit trails** - `updated_by`, `updated_at` fields, MongoDB logging
- **Concurrency control** - Pessimistic locking, distributed cache locking
- **Error handling** - Consistent error types, proper HTTP status codes
- **Context propagation** - Timeouts, cancellation, request IDs
- **Graceful degradation** - Cache fallback to database
- **Observability** - Structured logging, request tracing, performance metrics

## Interview Questions You Might Ask

- "Why pessimistic locking instead of optimistic?"
- "How would you handle distributed transactions across MySQL and MongoDB?"
- "What's the difference between your gRPC and REST clients, and when would you use each?"
- "How do you ensure request tracing works across async operations?"
- "What monitoring/observability gaps exist in this project?"

## Project Structure

```
internal/
├── core/               # Business domain (no framework coupling)
│   ├── services/       # Business logic
│   ├── ports/          # Interfaces for adapters
│   ├── dto/            # Data transfer objects
│   ├── constants/      # Error codes, business constants
│   └── helpers/        # Shared utilities
├── adapters/           # Infrastructure layer
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # CORS, logging, recovery, request ID
│   ├── client/         # gRPC and REST clients
│   ├── repository/     # SQL, Redis, MongoDB implementations
│   └── router/         # Route definitions
└── bootstrap/          # Dependency injection
cmd/api/main.go        # Application entry point
```
