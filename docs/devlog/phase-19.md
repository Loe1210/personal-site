# Phase 19 - Gateway Circuit Breaker Baseline

## Date

2026-07-23

## Scope

This phase adds a testable in-repository circuit breaker and wires it into the gateway's upstream boundaries without introducing a third-party breaker library or moving domain logic into gateway.

## Changes

- Added `internal/xresilience.CircuitBreaker` with closed, open, and half-open states.
- Added breaker snapshots for state, rejected calls, failure count, and recovery count.
- Added unified `CodeGatewayCircuitOpen` error semantics for circuit-open rejection.
- Wired gateway HTTP proxy breakers for:
  - `/api/auth/*` to auth-service,
  - `/api/content/*` and `/api/content/admin/*` to content-service,
  - `/api/media/*` to media-service.
- Wired the gateway auth Kitex client with a breaker around session validation and permission checks.
- Kept retry limited to idempotent reads. Non-idempotent writes are never replayed, but open breakers can reject them quickly.
- Converted upstream 5xx responses into the gateway upstream-failure envelope so write/upload failures do not look like successful proxy responses.
- Added breaker log fields: state, rejected calls, failures, and recoveries.
- Expanded smoke coverage for the gateway auth proxy path.

## Content Read Degradation

Public content reads can still degrade inside content-service through local cache, Redis cache, then MySQL. Gateway does not implement content cache fallback and does not return fake success when the content upstream is circuit-open.

## Out Of Scope

- No database circuit breaker in this phase.
- No automatic retry for login, writes, uploads, or upload complete.
- No deployment.

## Verification

Run before merge:

```bash
go test ./...
make proto-check
make micro-smoke
```
