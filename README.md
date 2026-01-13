# Order Management System (Go Backend)

This project follows **Clean Architecture** and **production-grade Go concurrency patterns**, suitable for backend interviews (Swiggy / Flipkart / Uber scale systems).

---

## 📁 Project Structure

```
OrderManagementSystem
│
├── cmd
│   └── server
│       └── main.go          # Application entry point (DI, server start)
│
├── internal
│   ├── contracts            # Interfaces (abstractions)
│   │   └── order_repository.go
│   │
│   ├── db                   # Database connection / setup
│   │   └── order_repo_db.go
│   │
│   ├── handler              # HTTP layer (request/response only)
│   │   └── order_handler.go
│   │
│   ├── middleware           # Auth, logging, rate limiting
│   │
│   ├── models               # Domain models / DTOs
│   │
│   ├── repository           # DB implementations
│   │   └── order_repo.go
│   │
│   ├── router               # Route registration
│   │   ├── order_routes.go
│   │   └── router.go
│   │
│   └── service              # Business logic + concurrency
│       └── order_service.go
│
├── go.mod
├── go.sum
└── README.md
```

---

## 🔁 Request Flow

```
HTTP Request
     ↓
Router
     ↓
Handler
     ↓
Service (Concurrency Layer)
     ↓
Contracts (Interfaces)
     ↓
Repository
     ↓
Database
     ↓
HTTP Response
```

👉 **Concurrency is implemented only in the Service layer**

---

## 🧠 Architecture Principles

* **Clean Architecture**
* **Dependency Inversion** using contracts (interfaces)
* **Loose coupling & testability**
* **Concurrency-safe design**

```
Handler → Service → Contracts ← Repository → DB
```

---

## ⚙️ Concurrency Used in This Project

| Concept     | Where Used               | Purpose                      |
| ----------- | ------------------------ | ---------------------------- |
| Goroutines  | Service                  | Parallel execution           |
| Channels    | Service                  | Coordination & communication |
| WaitGroup   | Service                  | Synchronization              |
| Mutex       | Inventory / Shared state | Data safety                  |
| RWMutex     | Cache                    | Read-heavy optimization      |
| Atomic      | Metrics                  | Lock-free counters           |
| Context     | All APIs                 | Cancellation & timeout       |
| Worker Pool | Bulk APIs                | Controlled concurrency       |
| Semaphore   | DB access                | Resource limiting            |
| sync.Pool   | Object reuse             | Reduce GC pressure           |
| sync.Once   | Initialization           | Run once                     |
| sync.Cond   | Inventory wait           | Advanced coordination        |

---

## 📌 Why This Structure?

* Easy to scale
* Easy to test
* Safe concurrency
* Matches real-world backend systems

**Interview explanation**:

> We isolate HTTP concerns, business logic, and persistence using clean architecture. Concurrency is handled in the service layer with proper synchronization and context propagation.

---

## 🧪 Testing Strategy

```
Handler Tests   → Mock Service
Service Tests   → Mock Repository
Repository Test → Real DB
```

---

## ▶️ Run the Application

```bash
go mod tidy
go run cmd/server/main.go
```

---

## 🚀 Future Improvements

* Redis caching
* Kafka integration
* gRPC APIs
* Graceful shutdown
* Rate limiting middleware
* Distributed tracing

---

## 👨‍💻 Author

**Nikhil Chauhan**
Backend Engineer (Go / .NET / SQL)

---

✅ This README is **interview-ready** and matches the exact project structure.
