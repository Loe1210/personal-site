# Microservice Timeout And Resilience

## Current Timeout Inventory

- `gateway` proxy previously used a fixed `http.Client{Timeout: 10s}` and returned ad hoc JSON for proxy failures.
- `gateway` upload guard already bounds upload request handling with a 2 minute context timeout and a concurrency gate.
- `gateway` auth RPC validation previously forwarded the caller context without a gateway-owned RPC deadline.
- `auth-service` and `content-service` Redis pools previously set only idle timeout; Redis commands did not have per-command read bounds.
- `auth-service`, `content-service`, and `media-service` MySQL repositories already accepted `context.Context`, but did not add a service-owned short deadline at the repository boundary.

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

The gateway remains thin: it owns proxy deadlines and auth validation deadlines, but does not own content, media, or auth business fallback logic.

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

## Error Semantics

Timeout and retry exhaustion use the shared application envelope through `xhttp.Fail` and `xerrors` codes. Gateway proxy failures no longer return ad hoc `message` fields.

Expected service business failures continue to use business codes. Transport or dependency failures are mapped to upstream failure or timeout codes at the boundary that observes them.

## Failure Statistics And Protection Principles

This phase does not introduce a circuit breaker library. The shared `xresilience.FailureStats` counter records total failures, consecutive failures, and last failure time so future protection can be added without changing call sites.

Protection principles for a future breaker:

- Break only at infrastructure boundaries, not inside domain services.
- Use conservative thresholds and short half-open probes.
- Never retry or replay non-idempotent writes.
- Prefer returning the unified timeout/upstream envelope over queueing unbounded work.
- Keep gateway protection generic; domain-specific fallback belongs in the owning service.

## Verification

Run before merging this phase:

```powershell
go test ./...
make proto-check
make micro-smoke
```
