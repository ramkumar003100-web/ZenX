# ZenX Framework Audit (Incremental Evolution)

## Architecture Map

- **Entry points**: CLI in `cmd/zenx`; core runtime usage in `pkg/examples`.
- **HTTP stack**: radix-style router with global + per-route middleware in `internal/router`.
- **Cross-cutting middleware**: security, timeout, CORS, rate limiting under `internal/middleware`.
- **Auth**: JWT, RBAC, password helpers in `internal/auth`.
- **Persistence**: SQL abstraction + migration executor in `internal/database`.
- **Async**: local priority queue and distributed queue helpers in `internal/jobs`.
- **Realtime**: WebSocket hub + manager under `internal/websocket`.
- **Observability**: Prometheus collectors and logging/tracing primitives in `internal/observability`.
- **Extensibility**: plugin registry/manager under `internal/plugins`; events in `internal/events`.

## Key Findings

### Architecture Issues

1. Router lacked first-class route grouping/versioning API, causing repeated path construction in applications.
2. Middleware utilities duplicated semantic checks (safe method detection).
3. Some packages expose powerful low-level APIs (database query builder) without clear safety boundaries.

### Code Quality Issues

1. `RateLimit` spawned a cleanup ticker goroutine per middleware construction without lifecycle shutdown, which can leak for dynamic router reload scenarios.
2. Sparse tests (no `_test.go` baseline) created high regression risk.
3. Several components rely on in-memory maps without eviction policies in long-lived services.

### Performance Bottlenecks

1. Router path splitting allocates per request (`strings.Split`) and method mismatch scans all trees for 405 checks.
2. Queue timing allocates new timers frequently when jobs are delayed.
3. Database query builder uses string formatting without statement caching.

### Security Weaknesses

1. JWT parsing validated algorithm but did not enforce issuer claim.
2. CSRF logic accepted any non-empty token (presence check only), useful baseline but not full integrity.
3. Request identity propagation was absent, reducing auditability for security events.

### Observability Gaps

1. No built-in request ID middleware for log correlation.
2. No panic recovery middleware in the default middleware package.
3. Limited standard middleware for request latency/trace injection in router pipeline.

## Implemented Improvements in This Iteration

1. Added router **grouping and API versioning primitives** (`Router.Group`, `Group.Version`, nested groups).
2. Added **request ID middleware** with context propagation and response header support.
3. Added **panic recovery middleware** with structured slog fields and stack trace logging.
4. Fixed rate limiter cleanup strategy to avoid background goroutine leaks.
5. Hardened JWT parsing with issuer validation.
6. Added initial automated tests for router grouping and request ID propagation.

## Next Iterations (Proposed)

1. Replace router split-path matching with zero-allocation path scanning and benchmarking.
2. Add explicit middleware interfaces for metrics/tracing with OpenTelemetry integration.
3. Introduce pluggable token validation strategy for CSRF and auth hardening.
4. Add queue dead-letter queue + retry policy abstractions.
5. Add package-level architecture tests to enforce boundary rules.
