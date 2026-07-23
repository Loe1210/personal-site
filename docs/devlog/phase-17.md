# Phase 17 - Thin Gateway Content Routing

## Date

2026-07-22

## Scope

This phase corrects the gateway-centered refactor direction for content APIs. Instead of moving article read business handling into `gateway`, the project keeps content-domain behavior inside `content-service` and makes `gateway` thinner.

Gateway should own cross-cutting HTTP concerns such as auth checks, upload guards, and route-level proxying. It should not own article/category/tag handling logic unless there is a concrete cross-service aggregation need.

## Changes

- Removed gateway's content-specific Kitex client package.
- Removed gateway's content-specific article HTTP handler package.
- Removed first-class gateway routes for:
  - `GET /api/articles`
  - `GET /api/articles/:id`
- Removed gateway startup initialization for the content Kitex client.
- Removed unused gateway content RPC environment variables from Docker Compose:
  - `CONTENT_SERVICE_NAME`
  - `CONTENT_RPC_ADDR`
- Kept gateway's thin content proxy routes:
  - `/api/content/*path`
  - `/api/content/admin/*path`
- Kept gateway auth protection on `/api/content/admin/*path` so admin content APIs are not exposed directly.
- Simplified frontend Nginx by removing public content compatibility rewrites for `/api/articles`, `/api/categories`, and `/api/tags`.
- Updated the public blog frontend API base to call `/api/content/*` directly.

## Rationale

The previous direction made `gateway` start owning content-specific handlers. That was useful for proving RPC wiring, but it made the gateway thicker than intended.

The revised boundary is:

```text
frontend
  -> gateway
  -> content-service
```

with gateway acting as a thin HTTP entrypoint and content-service retaining article/category/tag behavior.

## Verification

Run before release:

```bash
go test ./services/gateway/internal/router
go test ./services/gateway/cmd
go test ./...
make proto-check
make micro-smoke
```

## Follow-Up

- Decide whether auth for content admin APIs should eventually move into `content-service`; only then can gateway's admin content proxy be removed entirely.
- If a future page needs cross-service aggregation, add a dedicated gateway handler for that use case instead of recreating generic content CRUD in gateway.


## Error And Cache Hardening Addendum

This follow-up keeps the thin gateway boundary from Phase 17 and standardizes failure semantics across services.

- Added shared `idl/base/base.proto` with `BaseResp` and project-local error codes.
- Auth and content RPC handlers now use `base_resp` for expected business failures and reserve RPC errors for transport/unexpected failures.
- HTTP handlers return a unified envelope for existing routes: `code`, `msg`, and `data`. Business failures should not become HTTP 401/404/500 responses.
- Handler-owned errors are request parsing and binding failures. Service-owned errors are business outcomes such as expired sessions, permission denial, and missing articles.
- Unexpected panics are handled by shared base recover middleware; services also register a gopool panic handler through `xsafe`.
- `content-service` owns article detail caching with a local -> Redis -> MySQL lookup order. Gateway does not cache or inspect content-domain objects.
- Smoke now covers gateway `/api/content/articles` and confirms deprecated `/api/articles` remains a route-level 404.
