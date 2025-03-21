

# 🧾 Golang Microservice Architecture Guide

## 📘 Introduction

In a scalable microservice-based system, **consistency** and **separation of concerns** are critical. This document defines a **standardized folder architecture** for any Go microservice written in our team. It is based on **Clean Architecture (Hexagonal Architecture)** and ensures:
- Clear modularity
- Easy testing
- Loose coupling (through interfaces)
- Better onboarding for new devs
- Tech-agnostic business logic (e.g., switch MongoDB → PostgreSQL easily)

> ✅ No matter **what the service does** — message generator, chat processor, user auth, etc. — it will follow this **exact folder structure**.

---

## 📂 Folder Structure

```
your-service/
│
├── cmd/                     # App entry point and boot logic
│   └── server/main.go
│
├── internal/
│   ├── app/                 # Core business logic (Usecases)
│   ├── domain/              # Data models and interfaces (contracts)
│   ├── infra/               # External technologies (DB, Redis, MQ)
│   ├── delivery/            # HTTP, gRPC, MQTT handlers
│   └── config/              # Env loading and configuration
│
├── go.mod / go.sum
└── README.md
```

---

## 📚 Deep Dive into Each Layer

### 1. `cmd/server/main.go`
**What it does**:
- Loads config
- Creates connections (Mongo, Redis, etc.)
- Injects dependencies
- Starts HTTP server

> 🔁 You **wire** everything here like Lego blocks.

---

### 2. `internal/config/`
**What it does**:
- Loads environment variables or config files (JSON/YAML)
- Holds typed configuration

```go
type Config struct {
    Port string
    MongoURI string
    RedisAddr string
}
```

---

### 3. `internal/domain/`
**What it does**:
- Contains all **interfaces** and **data models**
- No technology-specific code
- Business logic depends only on this layer

```go
// MessageRepo defines DB contract
type MessageRepo interface {
    Save(chatID string, msg Message) error
    GetMessages(chatID string) ([]Message, error)
}
```

---

### 4. `internal/app/`
**What it does**:
- Implements business logic using only `domain` interfaces
- No direct dependency on Mongo, Redis, etc.
- Can be fully tested with mocks

```go
type ChatService struct {
    repo domain.MessageRepo
}

func (c *ChatService) SendMessage(chatID string, msg domain.Message) error {
    return c.repo.Save(chatID, msg)
}
```

---

### 5. `internal/infra/`
**What it does**:
- Implements the interfaces from `domain` using tech like MongoDB, Redis, RabbitMQ
- Nothing from `app/` should leak into `infra/`

```go
type MongoMessageRepo struct {
    coll *mongo.Collection
}

func (m *MongoMessageRepo) Save(chatID string, msg domain.Message) error {
    // Save msg in MongoDB
}
```

---

### 6. `internal/delivery/`
**What it does**:
- HTTP, gRPC, MQTT etc. endpoints
- Parses input, calls `app/` logic, returns response

```go
type Handler struct {
    chatService *app.ChatService
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
    // Parse, validate and call chatService.SendMessage()
}
```

---

## 🛠️ Example: Writing a HTTP Server

### 🏗 Folder layout
```
your-service/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── app/chat_service.go
│   ├── domain/message_repo.go
│   ├── infra/mongo/message_repo.go
│   ├── delivery/http/handler.go
│   └── config/config.go
```

---

### ✨ Step-by-step Flow

1. **`domain/`**
```go
// domain/message_repo.go
type MessageRepo interface {
    Save(chatID string, msg Message) error
}
```

2. **`app/`**
```go
// app/chat_service.go
type ChatService struct {
    repo domain.MessageRepo
}

func (c *ChatService) SendMessage(chatID string, msg domain.Message) error {
    return c.repo.Save(chatID, msg)
}
```

3. **`infra/`**
```go
// infra/mongo/message_repo.go
type MongoMessageRepo struct {
    coll *mongo.Collection
}

func (m *MongoMessageRepo) Save(chatID string, msg domain.Message) error {
    // Use m.coll.InsertOne()
}
```

4. **`delivery/`**
```go
// delivery/http/handler.go
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
    // Call chatService.SendMessage(...)
}
```

5. **`cmd/`**
```go
// cmd/server/main.go
func main() {
    cfg := config.Load()

    mongoClient := ConnectMongo(cfg.MongoURI)
    repo := mongo.NewMongoMessageRepo(mongoClient)

    chatService := app.NewChatService(repo)
    handler := delivery.NewHandler(chatService)

    http.ListenAndServe(":8080", handler.Routes())
}
```

---

## 🔌 Examples of Integrating External Services

### ✅ MongoDB
- Implement `MessageRepo` interface in `infra/mongo/`
- Initialize collection in `main.go`
- Inject into `ChatService`

### ✅ Redis
- Use for sessions or pub-sub
- Implement `SessionStore` interface in `infra/redis/`

### ✅ RabbitMQ
- For message queuing / background processing
- Define `MessageQueue` interface in `domain/`
- Implement in `infra/rabbitmq/`
- Use in `app/` for queuing jobs

---

## 📏 Rules of Thumb

| Rule | Description |
|------|-------------|
| No tech code in `app/` | Only interfaces & logic |
| No logic in `delivery/` | Just request parsing and response |
| No direct Mongo/Redis calls in `app/` | Go through `domain` |
| All dependencies created in `main.go` | One place to wire everything |
| Make interfaces in `domain/`, implement in `infra/` | Easy to mock, test, replace |

---

## 🧪 Optional Enhancements

| Feature | Where |
|--------|------|
| Middleware (auth, logging) | `delivery/http/middleware.go` |
| Unit tests | Inside each layer, mock `domain` |
| Docker | Add `Dockerfile` and `docker-compose.yml` |
| Makefile | Automate `run`, `lint`, `test` |

---

## 📚 Conclusion

By following this structure:
- Each developer knows where to add what
- All services stay consistent
- Business logic is decoupled from tech
- Code is reusable, testable, and scalable

> 🔁 Whether you're building a **simple message ID generator** or a **complex chat processor**, use this exact architecture.

---

## 🧰 Want This as a Template Repo?

Let me know — I can generate a **GitHub template repo** for you and your team to clone for any new microservice.

---

Let me know if you want this in PDF, Markdown, or README format for sharing with your team.
