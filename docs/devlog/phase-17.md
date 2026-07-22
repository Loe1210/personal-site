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

