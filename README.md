# ZenX Framework

ZenX is a modular Go backend framework designed for production APIs and enterprise workloads. It includes routing, middleware, authentication, database and cache layers, jobs, distributed patterns, observability, and developer tooling.

---

## Table of Contents

- [What ZenX provides](#what-zenx-provides)
- [Repository layout](#repository-layout)
- [Quick start](#quick-start)
- [CLI usage](#cli-usage)
- [Configuration](#configuration)
- [Core framework modules](#core-framework-modules)
- [Enterprise modules](#enterprise-modules)
- [Observability](#observability)
- [Security](#security)
- [OpenAPI and Swagger UI](#openapi-and-swagger-ui)
- [How to run](#how-to-run)
- [How to test](#how-to-test)
- [How to debug](#how-to-debug)
- [Integration flow example](#integration-flow-example)
- [Production rollout checklist](#production-rollout-checklist)
- [Known environment notes](#known-environment-notes)

---

## What ZenX provides

ZenX currently ships with:

1. **Router & middleware**
   - Radix-style route matching with path params.
   - Middleware chain composition.
   - Request context support and request IDs.
   - Timeout, CORS, in-memory rate limiting, and Redis-backed distributed rate limiting.

2. **Dependency injection**
   - Thread-safe service container.
   - Struct-field dependency resolution for handlers/services.

3. **Config & logging**
   - Environment variable config loader.
   - YAML/JSON/TOML config file loading.
   - Hot-reload watcher for file-based config.
   - Structured JSON logger with request ID propagation.

4. **Database layer**
   - MySQL and PostgreSQL driver support.
   - Migration runner from `migrations/*.sql`.
   - Transaction manager.
   - Query builder + lightweight ORM-like helpers (`Insert`, `Delete`).

5. **AuthN/AuthZ**
   - JWT signing and verification.
   - JWT auth middleware.
   - RBAC role checks.
   - Password hashing and verification (bcrypt).

6. **Validation engine**
   - JSON request binding.
   - Tag-driven validation (`required`, `min`, `max`, `email`, etc. from validator library).
   - Structured validation error responses.

7. **Jobs**
   - In-memory worker queue.
   - Delayed + recurring jobs.
   - Retry with exponential backoff.
   - Redis-backed distributed queue with dead-letter queue.

8. **Cache**
   - Redis integration with pooling.
   - KV, hash, TTL, exists, delete operations.

9. **OpenAPI / Swagger**
   - Programmatic OpenAPI spec builder.
   - Swagger UI endpoint.
   - Security/roles metadata in documented operations.

10. **Developer UX**
    - CLI scaffold: `zenx new <project>`.
    - Generated starter project with `/health` and `/metrics` endpoints.

11. **Enterprise additions**
    - WebSocket real-time hub + room broadcast.
    - Feature flags (global / user / role).
    - Optional GraphQL operation engine.
    - Notification services (SMTP email + webhook push retries).
    - File upload/storage service with optional AES encryption.
    - OpenTelemetry tracer wrapper.
    - Plugin hook manager.

---

## Repository layout

```text
zenx/
├─ cmd/zenx/
│  ├─ main.go
│  └─ new.go
├─ internal/
│  ├─ auth/
│  ├─ cache/
│  ├─ database/
│  ├─ jobs/
│  ├─ middleware/
│  ├─ openapi/
│  ├─ router/
│  ├─ validation/
│  ├─ websocket/
│  ├─ notifications/
│  ├─ featureflags/
│  ├─ storage/
│  ├─ tracing/
│  ├─ graphql/
│  └─ plugins/
├─ migrations/
├─ pkg/
│  ├─ config/
│  ├─ di/
│  ├─ logger/
│  ├─ metrics/
│  └─ examples/
└─ go.mod
```

---

## Quick start

### 1) Build CLI

```bash
go build -o zenx ./cmd/zenx
```

### 2) Scaffold a new project

```bash
./zenx new myservice
```

### 3) Run scaffolded service

```bash
cd myservice
go run ./cmd/server
```

### 4) Verify health and metrics

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

---

## CLI usage

```bash
zenx <command> [options]
```

Supported commands:

- `zenx new <project>`: create a starter application skeleton.

---

## Configuration

ZenX supports both **environment variables** and **file-based config**.

### Environment variables

| Variable | Description | Default |
|---|---|---|
| `ZENX_APP_NAME` | Service name | `ZenX` |
| `ZENX_ENV` | Environment | `development` |
| `ZENX_HTTP_ADDR` | HTTP listen address | `:8080` |
| `ZENX_REQUEST_TIMEOUT_SEC` | Request timeout in seconds | `15` |
| `ZENX_JWT_SECRET` | JWT signing secret | `dev-secret` |
| `ZENX_DATABASE_URL` | DB DSN | empty |
| `ZENX_REDIS_ADDR` | Redis address | `127.0.0.1:6379` |

### File config

`pkg/config` supports:

- `.yaml` / `.yml`
- `.json`
- `.toml`

You can also use `Watcher` for hot reload:

```go
w, err := config.NewWatcher("./configs/app.yaml")
if err != nil { panic(err) }

w.Start(2*time.Second, func(c config.Config) {
    // apply dynamic updates
})
```

---

## Core framework modules

### Router (`internal/router`)

- `Router.Handle(method, path, handler, middlewares...)`
- Path params via `:id` style segments.
- Automatic `404` and `405` semantics.
- Global middleware via `Use`.

### Middleware (`internal/middleware`)

Included middleware:

- `RequestID()`
- `Timeout(duration)`
- `CORS(origins)`
- `RateLimit(rps, burst)` (in-memory)
- `DistributedRateLimitTokenBucket(redis, namespace, burst, refill)`
- `SecureHeaders()`
- `CSRF(headerName)`
- `BruteForceGuard(maxAttempts, period)`

### DI (`pkg/di`)

Register singleton dependencies and resolve into structs:

```go
container := di.New()
container.Register(db)
container.Register(cache)

var handler struct {
    DB    *database.DB
    Cache *cache.RedisCache
}
_ = container.Resolve(&handler)
```

### Logger (`pkg/logger`)

- JSON structured logger.
- Request ID extraction from context for per-request logs.

### Validation (`internal/validation`)

Use binder + validator:

```go
var req CreateUserRequest
if err := validation.BindAndValidate(r, &req, validator); err != nil {
    // return centralized error response
}
```

### Auth (`internal/auth`)

- `JWTManager` for token issue/parse.
- `JWTAuth` middleware.
- `RequireRoles("admin")` middleware.
- `HashPassword` / `VerifyPassword`.

### Database (`internal/database`)

- Driver constants: `MySQLDriver`, `PostgresDriver`.
- `DB.Migrate(ctx, "./migrations")`.
- `TxManager.WithinTransaction`.
- Query builder: `Table("users").Where(...).Limit(...)`.
- ORM-like helpers: `Insert`, `Delete`.

### Cache (`internal/cache`)

`RedisCache` methods:

- `Set`, `Get`
- `HSet`, `HGet`
- `Exists`, `Delete`, `TTL`

### Jobs (`internal/jobs`)

- In-memory queue:
  - `NewQueue(workers)`
  - `Enqueue(job, delay, retries)`
  - `ScheduleRecurring(ctx, interval, job, retries)`
- Distributed queue (`DistributedQueue`): Redis sorted-set based scheduling, retry, DLQ.

---

## Enterprise modules

### WebSocket (`internal/websocket`)

- `Hub` for client lifecycle and room maps.
- `Manager` HTTP upgrader + authz callback.
- Broadcast APIs:
  - All clients
  - Room clients
  - Single client

### GraphQL (`internal/graphql`)

- Operation registration for query/mutation handlers.
- Role/authorization gate callback per operation.
- Lightweight JSON protocol (`type`, `name`, `args`).

### Feature flags (`internal/featureflags`)

`FlagRule` supports:

- `Global`
- per-role enables
- per-user enables

### Notifications (`internal/notifications`)

- SMTP templated email sender with retry/backoff.
- Webhook push sender with retry/backoff.

### Storage (`internal/storage`)

- Multipart upload size validation.
- Local filesystem backend.
- Optional AES encryption for uploaded payloads.

### Tracing (`internal/tracing`)

- Wrapper around OpenTelemetry tracer acquisition and span creation.

### Plugins (`internal/plugins`)

Hook interface:

- `OnRequest`
- `OnJobComplete`
- `OnShutdown`

---

## Observability

### Metrics (`pkg/metrics`)

Prometheus metrics included:

- `zenx_db_queries_total`
- `zenx_cache_operations_total`
- `zenx_jobs_runs_total`
- `zenx_websocket_connections`

Expose metrics endpoint in your HTTP mux:

```go
mux.Handle("/metrics", metrics.Handler())
```

### Health check

Recommended endpoint:

- `GET /health` -> `200 OK`

---

## Security

ZenX includes baseline protections:

- JWT authentication + RBAC authorization.
- bcrypt password hashing.
- CORS controls.
- CSRF middleware.
- Standard secure headers.
- Brute-force guard middleware.
- Per-route/global/user rate limiting patterns.

> For internet-facing production systems, also add TLS termination, secret management (vault/KMS), stricter origin policies, key rotation, and security scanning in CI.

---

## OpenAPI and Swagger UI

- Build specs with `openapi.Builder`.
- Add path metadata with optional bearer roles.
- Serve JSON via `builder.Handler()`.
- Serve UI via `openapi.SwaggerUI("/openapi.json")`.

Example:

```go
spec := openapi.New("ZenX API", "1.0.0")
spec.AddPath("/users", "GET", "List users", true, "admin")

mux.HandleFunc("/openapi.json", spec.Handler())
mux.HandleFunc("/swagger", openapi.SwaggerUI("/openapi.json"))
```

---

## How to use this framework in real projects

This section shows a practical pattern for building a ZenX application with **models, repositories, services, controllers, routes, and middleware**.

### 1) Suggested project structure

When using ZenX in an app repository, a common structure is:

```text
myservice/
├─ cmd/server/main.go
├─ internal/
│  ├─ models/
│  ├─ repositories/
│  ├─ services/
│  ├─ controllers/
│  ├─ middleware/
│  └─ routes/
├─ migrations/
└─ configs/
```

### 2) Define a model

```go
type User struct {
    ID       int64  `json:"id"`
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password,omitempty" validate:"required,min=8"`
}
```

### 3) Create a repository

Use `internal/database` query helpers and transaction support.

```go
type UserRepository struct {
    DB *database.DB
}

func (r *UserRepository) Create(ctx context.Context, u User) error {
    _, err := r.DB.Insert(ctx, "users", map[string]any{
        "email": u.Email,
        "name":  u.Name,
    })
    return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    q, args, _ := database.Table("users").
        Select("id", "email", "name").
        Where("email = ?", email).
        Limit(1).
        Build()

    row := r.DB.QueryRowContext(ctx, q, args...)
    var u User
    if err := row.Scan(&u.ID, &u.Email, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}
```

### 4) Create a service layer

Put business rules here (hashing, cache policy, async jobs, feature checks).

```go
type UserService struct {
    Users *UserRepository
    Cache *cache.RedisCache
    Jobs  *jobs.Queue
}

func (s *UserService) Register(ctx context.Context, u User) error {
    hash, err := auth.HashPassword(u.Password)
    if err != nil {
        return err
    }
    u.Password = hash

    if err := s.Users.Create(ctx, u); err != nil {
        return err
    }

    s.Jobs.Enqueue(jobs.JobFunc{
        JobName: "send_welcome_email",
        Fn: func(context.Context) error {
            // send email via notifications.EmailSender
            return nil
        },
    }, 0, 3)

    return nil
}
```

### 5) Create a controller (handler)

Controllers should focus on HTTP concerns: bind, validate, authorize, call service, respond.

```go
type UserController struct {
    Validator *validation.Validator
    Service   *UserService
}

func (ctl *UserController) Register(c *router.Context) {
    var req User
    if err := validation.BindAndValidate(c.Request, &req, ctl.Validator); err != nil {
        http.Error(c.Writer, err.Error(), http.StatusBadRequest)
        return
    }

    if err := ctl.Service.Register(c.Request.Context(), req); err != nil {
        http.Error(c.Writer, "failed to register user", http.StatusInternalServerError)
        return
    }

    c.Writer.WriteHeader(http.StatusCreated)
    _, _ = c.Writer.Write([]byte(`{"status":"created"}`))
}
```

### 6) Register dependencies with DI

```go
container := di.New()
container.Register(db)
container.Register(redisCache)
container.Register(validation.New())
container.Register(userRepo)
container.Register(userService)

var userController UserController
_ = container.Resolve(&userController)
```

### 7) Build middleware stack

A common order:
1. request metadata (`RequestID`)
2. security headers + CORS
3. timeout and rate limit
4. auth + RBAC for protected routes

```go
r := router.New()
r.Use(
    middleware.RequestID(),
    middleware.SecureHeaders(),
    middleware.CORS([]string{"https://app.example.com"}),
    middleware.Timeout(15*time.Second),
    middleware.RateLimit(20, 40),
)

jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.AppName, 24*time.Hour)

r.Handle("POST", "/auth/register", userController.Register)
r.Handle("GET", "/admin/users", adminListUsers,
    auth.JWTAuth(jwtMgr),
    auth.RequireRoles("admin"),
)
```

### 8) Add custom middleware

ZenX middleware signature:

```go
func MyMiddleware(next router.HandlerFunc) router.HandlerFunc {
    return func(c *router.Context) {
        start := time.Now()
        next(c)
        _ = start // track latency, log, etc.
    }
}
```

### 9) Add routes to OpenAPI

```go
spec := openapi.New("My Service", "1.0.0")
spec.AddPath("/auth/register", "POST", "Register user", false)
spec.AddPath("/admin/users", "GET", "List users", true, "admin")
```

### 10) Add health, metrics, and swagger endpoints

```go
mux := http.NewServeMux()
mux.Handle("/metrics", metrics.Handler())
mux.HandleFunc("/openapi.json", spec.Handler())
mux.HandleFunc("/swagger", openapi.SwaggerUI("/openapi.json"))
mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
})
```

### 11) Model migrations

Place SQL files in `migrations/` and run:

```go
if err := db.Migrate(ctx, "./migrations"); err != nil {
    log.Fatal(err)
}
```

### 12) Debugging flow for new features

When adding a new controller or middleware:

1. create request/response model and validation tags.
2. add repository queries and service methods.
3. wire controller with DI.
4. attach middleware at route or global level.
5. add OpenAPI route metadata.
6. test endpoint via curl/Postman and verify metrics/logs.

---

## How to run

### Local run (framework repo)

```bash
go run ./cmd/zenx
```

### Build and run CLI

```bash
go build -o zenx ./cmd/zenx
./zenx new demo
cd demo
go run ./cmd/server
```

---

## How to test

### Format

```bash
gofmt -w $(rg --files -g '*.go')
```

### Unit / package tests

```bash
go test ./...
```

### Targeted packages (if external dependencies are restricted)

```bash
go test ./internal/router ./internal/openapi ./internal/featureflags ./internal/graphql ./internal/plugins ./internal/storage ./pkg/di ./pkg/logger
```

---

## How to debug

### Basic runtime debugging

- Run with verbose logging in development (`ZENX_ENV=development`).
- Include request IDs in log lines via middleware + context logger helpers.
- Validate middleware order (auth, timeout, rate-limit, CORS, etc.) for expected behavior.

### Delve

```bash
dlv debug ./cmd/zenx
```

### Useful checks

```bash
# verify routing behavior
curl -i http://localhost:8080/unknown

# verify rate-limits
hey -n 200 -c 20 http://localhost:8080/health

# verify websocket handshake
wscat -c ws://localhost:8080/ws
```

---

## Integration flow example

Typical production request flow:

1. Request enters router.
2. RequestID + tracing middleware annotate context.
3. Security middleware (headers, CSRF, brute-force, CORS).
4. Authentication middleware (JWT).
5. Authorization middleware (RBAC role check).
6. Validation binder for request body.
7. Business handler uses DI-managed services.
8. Database/cache operations emit metrics.
9. Async work enqueued to jobs queue.
10. OpenAPI/Swagger and Prometheus endpoints remain available for ops.

You can find sample usage references in `pkg/examples/usage.go`.

---

## Production rollout checklist

- [ ] Use strong JWT secret and rotate regularly.
- [ ] Configure Redis/DB with secure credentials and TLS where applicable.
- [ ] Set strict CORS allow-list.
- [ ] Enable distributed rate limiting for multi-node deployments.
- [ ] Enable tracing exporter and central log aggregation.
- [ ] Configure alerting on Prometheus metrics.
- [ ] Run migrations in CI/CD before deployment.
- [ ] Configure dead-letter queue monitoring for distributed jobs.
- [ ] Add integration and load tests for critical APIs.

---

## Known environment notes

In restricted environments, `go test ./...` may fail if dependencies cannot be fetched to produce `go.sum` entries. In that case:

1. Run targeted package tests that do not require external module resolution.
2. Re-run `go mod tidy` and full tests once network/module access is available.

---

## License

Add your preferred license (MIT/Apache-2.0/etc.) at repository root.
