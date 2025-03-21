

## ✅ Goal Summary

- 📎 **Service**: A **Go HTTP server** that returns **globally unique, time-sortable** message IDs.
- ⚙️ **Architecture**: Modular, **clean layered design** following **Dependency Inversion Principle (DIP)**.
- 🧱 **Reusable**: Same architecture can be used for **any future microservice**.
- 📡 **Usage**: Other services call `GET /message-id` → Receive `{"id": "..."}`

---

## 📁 Folder Structure (Clean Architecture)

```
message-id-service/
│
├── cmd/                     # Entry point for app
│   └── server/              # Server startup
│       └── main.go
│
├── internal/
│   ├── app/                 # Core application logic (Usecases)
│   │   └── idgen.go
│   │
│   ├── domain/              # Interfaces, domain models
│   │   └── idgen.go         # ID generator interface
│   │
│   ├── infra/               # External adapters
│   │   └── snowflake/       # Snowflake ID implementation
│   │       └── generator.go
│   │
│   ├── delivery/            # HTTP handlers/controllers
│   │   └── http/
│   │       └── handler.go
│   │
│   └── config/              # Config loading (if needed later)
│       └── config.go
│
├── go.mod
└── README.md
```

> This structure follows **Hexagonal/Clean Architecture**. All core logic depends on abstractions, not implementations.

---

## 🧩 Key Components

### 1. `domain/idgen.go` — Interface

```go
package domain

type IDGenerator interface {
    GenerateID() (string, error)
}
```

---

### 2. `infra/snowflake/generator.go` — Actual Implementation

We'll use a Snowflake-like library for Go: `github.com/bwmarrin/snowflake`.

```go
package snowflake

import (
    "github.com/bwmarrin/snowflake"
    "strconv"
    "message-id-service/internal/domain"
)

type SnowflakeGenerator struct {
    node *snowflake.Node
}

func NewSnowflakeGenerator(nodeID int64) (domain.IDGenerator, error) {
    node, err := snowflake.NewNode(nodeID)
    if err != nil {
        return nil, err
    }
    return &SnowflakeGenerator{node: node}, nil
}

func (s *SnowflakeGenerator) GenerateID() (string, error) {
    id := s.node.Generate()
    return strconv.FormatInt(id.Int64(), 10), nil
}
```

---

### 3. `app/idgen.go` — Usecase Layer

```go
package app

import "message-id-service/internal/domain"

type IDService struct {
    generator domain.IDGenerator
}

func NewIDService(gen domain.IDGenerator) *IDService {
    return &IDService{generator: gen}
}

func (s *IDService) GetID() (string, error) {
    return s.generator.GenerateID()
}
```

---

### 4. `delivery/http/handler.go` — HTTP Controller

```go
package http

import (
    "encoding/json"
    "net/http"
    "message-id-service/internal/app"
)

type Handler struct {
    service *app.IDService
}

func NewHandler(service *app.IDService) *Handler {
    return &Handler{service: service}
}

func (h *Handler) GetID(w http.ResponseWriter, r *http.Request) {
    id, err := h.service.GetID()
    if err != nil {
        http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
        return
    }

    response := map[string]string{"id": id}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

### 5. `cmd/server/main.go` — Entry Point

```go
package main

import (
    "log"
    "net/http"

    "message-id-service/internal/app"
    "message-id-service/internal/delivery/http"
    "message-id-service/internal/infra/snowflake"
)

func main() {
    gen, err := snowflake.NewSnowflakeGenerator(1)
    if err != nil {
        log.Fatalf("Error initializing Snowflake: %v", err)
    }

    service := app.NewIDService(gen)
    handler := http.NewHandler(service)

    http.HandleFunc("/message-id", handler.GetID)

    log.Println("Message ID Service running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 🏗️ How This Structure Helps

| Layer | Responsibility | Depends On |
|-------|----------------|------------|
| `domain` | Interface definitions | None |
| `infra` | Concrete implementations | `domain` |
| `app` | Business logic | `domain` |
| `delivery` | Expose HTTP API | `app` |
| `cmd` | Wire everything | All above |

### 🔄 Dependency Flow

```
infra  --> domain <-- app <-- delivery <-- cmd
```

> All concrete layers depend on interfaces, never the other way. This is **Dependency Inversion Principle (DIP)** in action.

---

## 🔌 Install Required Dependency

In `go.mod`:

```bash
go get github.com/bwmarrin/snowflake
```

---

## ✅ Example Response

```bash
GET /message-id
→ { "id": "275889400987365376" }
```

These IDs:
- Are **globally unique**
- Are **roughly time-sortable**
- Can be **used for ordering** messages

---

## 🧠 Optional Enhancements

| Feature | Description |
|--------|-------------|
| `HealthCheck` | Add `/health` endpoint |
| `Metrics` | Use Prometheus/Grafana |
| `Rate limit` | Avoid abuse (e.g., IP-based throttle) |
| `Distributed Deployment` | Use different nodeIDs for each instance |

---

Would you like me to:
- Create a GitHub-ready version of this?
- Generate test cases for the IDService?
- Dockerize it?

Let me know, I can help you set it up fully!