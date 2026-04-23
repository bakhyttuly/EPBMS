# EPBMS Analysis and Architecture Redesign

## 1. Analysis of Existing Code

### Bugs and Bad Practices Identified
1. **Poor Architecture:** The project currently uses a flat structure (`handlers`, `models`, `services`, `routes`). Business logic is heavily mixed into the HTTP handlers (e.g., `CreateBooking` in `booking_handler.go` handles role checks, database queries, conflict resolution, and HTTP responses).
2. **Weak Validation:** Request bodies are bound using `ShouldBindJSON` but lack structural validation tags (e.g., `binding:"required"`). This allows empty or malformed data to be saved to the database.
3. **Incomplete Role Logic:** The existing roles are `admin`, `organizer`, and `performer`. The new requirements specify `ADMIN`, `PERFORMER`, and `CLIENT`. The current implementation allows organizers to create bookings directly, bypassing the required "pending" state and admin approval flow.
4. **Session Management:** The project uses cookie-based sessions (`github.com/gin-contrib/sessions`) instead of stateless JWTs, which is less scalable for REST APIs and mobile clients.
5. **Database Queries in Handlers:** Handlers directly access `config.DB` to perform CRUD operations, making unit testing impossible without a real database.
6. **Conflict Detection Flaws:** The `HasBookingConflict` function parses times using `time.Parse("15:04", ...)` on every check. It fetches all bookings for a performer on a specific date and iterates through them in memory. This is inefficient and prone to race conditions under concurrent load.
7. **Hardcoded Configuration:** Database credentials and session keys are hardcoded in `config/db.go` and `main.go`.

## 2. New Architecture Design

The project will be refactored using **Clean Architecture** principles to separate concerns, improve testability, and ensure scalability.

### Layers
1. **Domain Layer (`internal/domain`):** Contains core business models (entities) and interface definitions for repositories and services.
2. **Repository Layer (`internal/repository`):** Implements domain interfaces for data persistence using GORM. Handles all database interactions.
3. **Service Layer (`internal/service`):** Contains business logic (e.g., conflict detection, role validation, booking state transitions). Depends on repository interfaces.
4. **Handler Layer (`internal/handler`):** Handles HTTP requests and responses using Gin. Depends on service interfaces.

### Database Schema Improvements
- **Users:** `id`, `email`, `password`, `role` (ADMIN, PERFORMER, CLIENT), `created_at`, `updated_at`.
- **Performers:** `id`, `user_id` (FK to users), `name`, `category`, `price`, `description`, `created_at`, `updated_at`.
- **Bookings:** `id`, `performer_id` (FK to performers), `client_id` (FK to users), `event_date` (Date), `start_time` (Time), `end_time` (Time), `status` (PENDING, CONFIRMED, REJECTED, COMPLETED), `created_at`, `updated_at`.

### Booking Flow Logic
1. **CLIENT** creates a booking request -> Status: `PENDING`.
2. **ADMIN** reviews the request.
   - If approved: Checks for conflicts. If no conflicts -> Status: `CONFIRMED`.
   - If rejected -> Status: `REJECTED`.
3. **PERFORMER** can only view bookings where Status is `CONFIRMED` and `performer_id` matches their ID.

### API Contract (RESTful)
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/performers` (Public/Client)
- `GET /api/v1/performers/:id` (Public/Client)
- `POST /api/v1/bookings` (Client only)
- `GET /api/v1/bookings` (Admin sees all, Performer sees own confirmed, Client sees own)
- `PUT /api/v1/admin/bookings/:id/status` (Admin only - approve/reject)
- `GET /api/v1/performers/schedule` (Performer only)
