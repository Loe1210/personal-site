# Phase 18 - Timeout And Resilience Baseline

## Date

2026-07-23

## Scope

This phase establishes the first shared microservice resilience baseline: timeout hierarchy, limited idempotent-read retry, unified timeout/upstream errors, and failure statistics without introducing a circuit breaker dependency.

## Changes

- Added `internal/xresilience` with default timeout budgets, idempotent read detection, bounded retry, timeout classification, and failure statistics.
- Updated gateway proxy to:
  - use a 6 second proxy timeout,
  - retry only `GET` and `HEAD` once,
  - avoid retry for writes and uploads,
  - return shared envelope errors for proxy timeout/upstream failures.
- Updated gateway auth RPC client to add an 800ms deadline and bounded retry for read-only session/permission checks.
- Added Redis command timeouts for auth session store and content article cache.
- Added Redis connect/read/write timeout options for auth-service and content-service Redis pools.
- Added MySQL repository boundary timeouts in auth-service, content-service, and media-service.
- Added unit tests for shared retry/timeout helpers, gateway proxy retry boundaries, auth RPC deadline, and Redis command timeouts.
- Expanded smoke validation for gateway content response envelope shape.
- Documented timeout hierarchy, retry scope, unified errors, and future breaker principles.

## Constraints Preserved

- Gateway remains a thin proxy and auth middleware layer.
- Login, writes, upload requests, and upload complete do not use automatic retry.
- Timeout and retry exhaustion use shared application error semantics.
- Circuit breaker behavior is intentionally limited to failure counting and documented protection principles.

## Verification

Run before merge:

```bash
go test ./...
make proto-check
make micro-smoke
```
