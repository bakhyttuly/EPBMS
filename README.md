# EPBMS — Event Performer & Booking Management System

A production-grade REST API built with Go, Gin, GORM, and PostgreSQL.

---

## Architecture

The project follows **Clean Architecture** principles with strict layer separation:

```
backend/
├── cmd/api/           # Entry point — wires all dependencies (DI root)
├── config/            # Database initialisation (env-based, no hardcoded secrets)
├── internal/
│   ├── domain/        # Entities, interfaces (repository + service), errors, DTOs
│   ├── repository/    # GORM implementations of domain.Repository interfaces
│   ├── service/       # Business logic (auth, performers, bookings)
│   │   └── mocks/     # In-memory mocks for unit testing
│   ├── handler/       # Gin HTTP handlers (thin layer, delegates to services)
│   └── middleware/    # JWT auth, RBAC role guard, request logger, rate limiter
├── pkg/
│   ├── utils/         # JWT generation & parsing
│   ├── logger/        # Structured slog logger
│   └── response/      # Standardised JSON response envelope
└── routes/            # Route registration with role-based middleware
```

**Dependency flow:** `handler → service → repository → database`  
Each layer depends only on the **interface** defined in `domain/`, never on a concrete implementation. This makes every layer independently unit-testable.

---

## Role-Based Access Control

| Role          | Permissions |
|---------------|-------------|
| **ADMIN**     | Full access. Sees all bookings. Approves/rejects/completes bookings. Manages performers. |
| **PERFORMER** | Sees only their own **confirmed** bookings. Can update their own profile. |
| **CLIENT**    | Browses performers. Creates booking requests (status: `pending`). Sees only their own bookings. |

---

## Booking Flow

```
CLIENT  →  POST /api/v1/bookings          →  status: pending
ADMIN   →  PUT  /api/v1/admin/bookings/:id/status
              body: { "status": "confirmed" }  →  conflict check → status: confirmed
              body: { "status": "rejected"  }  →  status: rejected
PERFORMER  →  GET /api/v1/bookings        →  sees only confirmed bookings assigned to them
```

---

## Conflict Detection

When an ADMIN **confirms** a booking, the system performs a database-level overlap check:

```sql
WHERE performer_id = ?
  AND event_date   = ?
  AND status IN ('pending', 'confirmed')
  AND start_time < :new_end_time
  AND end_time   > :new_start_time
```

This is the standard **interval overlap** condition. If any row matches, the confirmation is rejected with `409 Conflict` and a clear error message.

---

## API Endpoints

### Auth (public)
| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive a JWT |

### Performers (JWT required)
| Method | Path | Roles | Description |
|--------|------|-------|-------------|
| `GET` | `/api/v1/performers` | All | List performers (paginated, filterable by `category`) |
| `GET` | `/api/v1/performers/:id` | All | Get performer by ID |
| `POST` | `/api/v1/performers` | Admin | Create performer profile |
| `PUT` | `/api/v1/performers/:id` | Admin, Performer (own) | Update performer profile |
| `DELETE` | `/api/v1/performers/:id` | Admin | Delete performer |

### Bookings (JWT required)
| Method | Path | Roles | Description |
|--------|------|-------|-------------|
| `GET` | `/api/v1/bookings` | All (role-scoped) | List bookings |
| `GET` | `/api/v1/bookings/:id` | All (role-scoped) | Get booking by ID |
| `POST` | `/api/v1/bookings` | Client | Create booking request |

### Admin (JWT + Admin role required)
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/admin/stats` | Dashboard statistics |
| `PUT` | `/api/v1/admin/bookings/:id/status` | Approve / reject / complete a booking |
| `DELETE` | `/api/v1/admin/bookings/:id` | Delete a booking |

---

## Standard Response Envelope

All responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "meta": { "page": 1, "page_size": 20, "total": 42 }
}
```

Errors:

```json
{
  "success": false,
  "error": "booking time conflicts with an existing confirmed booking"
}
```

---

## Setup & Running

### Prerequisites
- Go 1.23+
- PostgreSQL 14+

### Steps

```bash
# 1. Clone and enter backend directory
cd backend

# 2. Copy and configure environment
cp .env.example .env
# Edit .env with your DB credentials and JWT secret

# 3. Download dependencies
go mod download

# 4. Run the server
go run ./cmd/api/main.go
```

The server starts on `http://localhost:8080` by default.

---

## Running Tests

```bash
cd backend

# All tests
go test ./...

# With verbose output
go test ./internal/service/... ./pkg/utils/... -v
```

**16 unit tests** covering:
- Booking conflict detection (no conflict, with conflict, invalid transitions)
- Role-based visibility (CLIENT, PERFORMER, ADMIN)
- JWT generation, parsing, and tamper detection

---

## Security

- **JWT** (HS256) with 24-hour expiry, secret loaded from environment variable.
- **RBAC middleware** enforces role checks at the route level before handlers execute.
- **Passwords** hashed with `bcrypt` (cost 10).
- **Rate limiter**: 10 req/s per IP, burst of 30 (token bucket algorithm).
- **No hardcoded secrets** — all sensitive values read from environment variables.

---

## Database Schema

```sql
users        (id, full_name, email, password, role, created_at, updated_at)
performers   (id, user_id FK→users, name, category, price, description, created_at, updated_at)
bookings     (id, performer_id FK→performers, client_id FK→users,
              event_date, start_time, end_time, status,
              notes, approved_by, approved_at, created_at, updated_at)
```

Indexes on: `bookings.performer_id`, `bookings.client_id`, `bookings.event_date`, `bookings.status`, `users.email` (unique), `performers.user_id` (unique).
