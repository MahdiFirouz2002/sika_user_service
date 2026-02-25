# Concurrent User Importer & API (Golang + Echo + PostgreSQL)

This project implements a concurrent system that:

* Reads **10,000 users** from a JSON file
* Inserts them into a **PostgreSQL** database
* Maintains a **one-to-many relationship** between `users` and `addresses`
* Limits concurrency to **10 workers**
* Exposes an **HTTP API** to fetch a user with all addresses by ID

The project follows a **Clean Architecture (3-layered structure)** approach and is production-ready.

---

# Tech Stack

* **Golang**
* **PostgreSQL**
* **GORM**
* **Echo (HTTP framework)**
* **Channel-based worker pool**
* **Graceful shutdown**
* **Environment-based configuration**

---

# Architecture

The project is structured based on Clean Architecture principles, separating responsibilities into independent layers.

```
cmd/
  └── api/                → Application entry point

internal/
  ├── database/           → DB connection & GORM models
  │   └── models/
  ├── domain/             → Business entities & errors
  ├── repository/         → Data access layer
  ├── service/            → Business logic layer
  ├── server/
  │   ├── controller/     → HTTP handlers
  │   └── http/           → HTTP server setup
pkg/
  └── config/             → Configuration loader
```

---

## Layer Responsibilities

### 1️⃣ Domain Layer (`internal/domain`)

* Defines core business entities:

  * `User`
  * `Address`
* Contains business-level errors (e.g. `ErrUserNotFound`)
* Has **no dependency** on external libraries

This is the heart of the application.

---

### 2️⃣ Repository Layer (`internal/repository`)

* Responsible for database operations
* Uses **GORM**
* Converts between:

  * `domain` models
  * `database` models

Example:

* `Create(user)`
* `Get(id)` (with `Preload("Addresses")`)

This layer depends on:

* Database
* Domain

---

### 3️⃣ Service Layer (`internal/service`)

Contains business logic:

* `UserService`

  * Create user
  * Get user by ID
* `Producer`

  * Reads JSON file
  * Sends users to channel
  * Manages worker pool (10 concurrent workers)

This layer orchestrates repositories but does not know about HTTP.

---

### 4️⃣ HTTP Layer (`internal/server`)

Built using **Echo**.

* `controller` → Handles HTTP requests
* `http` → Starts server & registers routes

Example endpoint:

```
GET /:id
```

Returns:

```json
{
  "user": {
    "id": "...",
    "name": "...",
    "email": "...",
    "phone_number": "...",
    "addresses": [...]
  }
}
```

Includes:

* UUID validation
* Proper error handling
* Graceful shutdown

---

# Concurrency Design

The insertion logic uses a **Worker Pool Pattern**:

1. JSON file is read
2. Users are unmarshaled
3. Users are sent into a buffered channel
4. Exactly **10 workers** consume from the channel
5. Each worker inserts users into the database

```go
producer := service.NewProducer(userService, 10)
producer.RunInsert(ctx)
```

This guarantees:

* Maximum 10 concurrent database operations
* Remaining jobs wait in queue
* Controlled resource usage
* Safe concurrency using `sync.WaitGroup`

---

# Database Design

## users table

| Column       | Type |
| ------------ | ---- |
| id (UUID)    | PK   |
| name         | TEXT |
| email        | TEXT |
| phone_number | TEXT |

## addresses table

| Column   | Type |
| -------- | ---- |
| id       | PK   |
| user_id  | FK   |
| street   | TEXT |
| city     | TEXT |
| state    | TEXT |
| zip_code | TEXT |
| country  | TEXT |

Relationship:

```
User 1 ---- * Address
```

Implemented via GORM associations.

---

# Running the Project

## 1️⃣ Set Environment Variables

Create `.env` file:

```
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=users_db
DATABASE_PORT=5432

HTTP_SERVER_HOST=0.0.0.0
HTTP_SERVER_PORT=8080
```

---

## 2️⃣ Run Database Migration

Migration runs automatically on startup.

---

## 3️⃣ Insert Users (JSON → DB)

Place your file as:

```
users_data.json
```

Run:

```
go run cmd/api/main.go --mode=insert
```

This will:

* Read JSON file
* Start 10 workers
* Insert all users concurrently

---

## 4️⃣ Run HTTP Server

```
go run cmd/api/main.go --mode=server
```

Server runs with graceful shutdown support.

---

# Example API Call

```
GET http://localhost:8080/{user_id}
```

Example:

```
GET http://localhost:8080/83aab3ca-b0fc-409c-9cb8-60916e381c03
```

---

# Production Considerations

* Environment-based config loading
* Graceful shutdown with context
* Controlled concurrency (worker pool)
* Clear separation of concerns
* Database migration on startup
* Proper error handling
* UUID validation

---

# Design Decisions

### Why Worker Pool?

To strictly enforce:

> Only 10 concurrent database insert operations at a time.

### Why Clean Architecture?

* Independent layers
* Easy to test
* Easy to replace DB or HTTP framework
* Scalable structure for large systems

### Why Channel?

Channels provide:

* Safe communication between producer and workers
* Natural job queue
* Built-in concurrency control
