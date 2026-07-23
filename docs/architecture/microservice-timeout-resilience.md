# Microservice Timeout And Resilience

## Current Timeout Inventory

- `gateway` proxy previously used a fixed `http.Client{Timeout: 10s}` and returned ad hoc JSON for proxy failures.
- `gateway` upload guard bounds upload request handling with a 2 minute context timeout and a concurrency gate.
- `gateway` auth RPC validation uses a gateway-owned RPC deadline.
- `auth-service` and `content-service` Redis pools set connect/read/write timeouts, and Redis commands use per-command timeouts.
- `auth-service`, `content-service`, and `media-service` MySQL repositories accept `context.Context` and add a service-owned short deadline at the repository boundary.

## Timeout Hierarchy

Timeouts stay layered from outside to inside. Inner dependencies should fail before the caller's whole request budget is exhausted.

```text
client/browser request
  -> gateway proxy timeout: 6s
  -> auth RPC validation timeout: 800ms
  -> MySQL repository timeout: 1.5s
  -> Redis command/connect/read/write timeout: 300ms
  -> upload guard timeout: 2m for upload bodies only
```

The gateway remains thin: it owns proxy deadlines, auth validation deadlines, retry gates, and breaker gates, but it does not own content, media, or auth business fallback logic.

## Retry Policy

Automatic retry is limited to idempotent reads:

- Gateway proxy retries only `GET` and `HEAD`.
- Retry budget is two total attempts, with a short fixed backoff.
- Retryable failures are transport timeouts/errors and upstream 5xx responses for idempotent reads.
- Auth RPC validation and permission checks use the same bounded read retry because they are read-only.

No automatic retry is applied to:

- Login.
- Article/category/tag writes.
- Upload requests.
- Upload complete/finalization.
- Any POST, PUT, PATCH, or DELETE proxy request.

Non-idempotent requests can still fail fast when their upstream breaker is already open, but they are never replayed.

## Circuit Breaker V1

`internal/xresilience.CircuitBreaker` is a small in-repository state machine with three states:

```text
closed -> open -> half_open -> closed
              \-> open on failed probe
```

Default behavior:

- Consecutive infrastructure failures open the breaker after the configured threshold.
- Open breakers reject immediately with `ErrCircuitOpen`.
- After the open timeout, the next call becomes a half-open probe.
- A successful half-open probe closes the breaker and increments recovery count.
- A failed half-open probe opens the breaker again.

Gateway???:

- `gateway -> auth-service` HTTP proxy: `/api/auth/*`.
- `gateway -> content-service` HTTP proxy: `/api/content/*` and `/api/content/admin/*`.
- `gateway -> media-service` HTTP proxy: `/api/media/*`.
- `gateway -> auth-service` Kitex RPC: session validation and permission checks.

Database breaker is intentionally not implemented in this phase. Database calls keep repository timeouts only. If database breaker is added later, it should be service-local and owned by each domain service, not by gateway.

## Content Read Degradation

Public content reads can degrade through content-service's existing cache chain:

```text
content-service ArticleService
  -> local article cache
  -> Redis article cache
  -> MySQL article repository
```

Redis cache failures are ignored when MySQL can answer. Gateway does not serve cached content itself and does not return fake success when content-service is unavailable or circuit-open.

## Error Semantics

Timeout, retry exhaustion, and circuit-open rejection use the shared application envelope through `xhttp.Fail` and `xerrors` codes. Gateway proxy failures do not return ad hoc `message` fields.

Expected service business failures continue to use business codes. Transport or dependency failures are mapped to upstream failure, timeout, or circuit-open codes at the boundary that observes them.

## Metrics And Logs

The breaker snapshot exposes these fields:

- `breaker_state`
- `breaker_rejected`
- `breaker_failures`
- `breaker_recoveries`

Gateway proxy and auth RPC client log these fields on success, failure, and rejection. This is intentionally basic so the next observability phase can wire the same snapshot into metrics without changing breaker behavior.

## Verification

Run before merging this phase:

```powershell
go test ./...
make proto-check
make micro-smoke
```
