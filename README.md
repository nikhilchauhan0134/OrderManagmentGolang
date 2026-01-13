# OrderManagmentGolang
Demo of how we create golang with workpolling with gorouting and connect with sql server
Project Folder Structure (Diagram)
OrderManagementSystem
│
├── cmd
│   └── server
│       └── main.go          ← App entry point (DI, server start)
│
├── internal
│   │
│   ├── contracts            ← Interfaces (abstractions)
│   │   └── order_repository.go
│   │
│   ├── db                   ← DB connection / config
│   │   └── order_repo_db.go
│   │
│   ├── handler              ← HTTP layer (no business logic)
│   │   └── order_handler.go
│   │
│   ├── middleware           ← Auth, logging, rate limit
│   │
│   ├── models               ← Domain models / DTOs
│   │
│   ├── repository           ← DB implementations
│   │   └── order_repo.go
│   │
│   ├── router               ← Route registration
│   │   ├── order_routes.go
│   │   └── router.go
│   │
│   └── service              ← Business logic + concurrency
│       └── order_service.go
│
├── go.mod
└── go.sum
Request → Response Flow

HTTP Request
     │
     ▼
Router (order_routes.go)
     │
     ▼
Handler (order_handler.go)
     │
     ▼
Service (order_service.go)
     │
     ▼
Contracts (interfaces)
     │
     ▼
Repository (order_repo.go)
     │
     ▼
DB (order_repo_db.go)
     │
     ▼
HTTP Response


Dependency Direction (Clean Architecture)

Handler
   ↓
Service
   ↓
Contracts  ←── Repository implements this
   ↑
Repository
   ↓
DB

