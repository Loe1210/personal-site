# Microservices Current State

## Services

- `frontend`: Nginx-served static site and static asset entrypoint.
- `gateway`: external HTTP entrypoint. It owns cross-cutting HTTP middleware such as auth checks and upload guards, but it should not own content-domain business handlers.
- `auth-service`: login, session, user, role, and permission domain.
- `content-service`: article, category, and tag domain.
- `media-service`: upload task, uploaded file metadata, and local file storage domain.

## Current Communication

- Browser traffic enters the frontend container first.
- Frontend API traffic flows through `gateway` through the generic `/api/` Nginx proxy.
- Public blog code calls content APIs with explicit `/api/content/*` paths.
- `gateway` proxies `/api/content/*` to `content-service` without article/category/tag business handling.
- `gateway` protects `/api/content/admin/*` with auth middleware before proxying to `content-service`.
- `gateway` proxies `/api/auth/*` and `/api/media/*` to their domain services.
- `auth-service`, `content-service`, and `media-service` start Kitex RPC servers.
- `gateway` uses Kitex RPC for auth session validation; content-domain browser APIs remain HTTP proxy calls to keep gateway thin.

## Current Contracts

- `idl/auth/auth.proto` exposes `ValidateSession` and `CheckPermission`.
- `idl/content/content.proto` exposes article list, detail, create, update, and delete RPC methods for service-to-service use when a real caller needs them.
- `idl/media/media.proto` exposes `GetFile`.
- Existing proto contracts should be used before adding new methods.
- New proto methods should be added only when a concrete gateway or service caller needs them in the same implementation slice.

## Target Runtime Topology

```text
Browser
  -> frontend Nginx static site
  -> /api/* forwarded to gateway
  -> gateway Hertz middleware and thin routing/proxy layer
  -> auth-service Kitex RPC for gateway session validation
  -> content-service HTTP for content-domain browser APIs
  -> media-service HTTP for upload streaming and media APIs
```

## Data Ownership

- `auth-service` owns `auth_db`.
- `content-service` owns `content_db`.
- `media-service` owns `media_db`.
- There is no generic database service.
- Services must not read another service's database tables directly.

## Immediate Baseline Checks

Run these before implementation changes:

```powershell
go test ./...
make proto-check
make micro-smoke
```

If Docker or Kubernetes is unavailable locally, record that as an environment limitation rather than a code failure.

## Baseline Result on 2026-07-22

- `go test ./...`: passed.
- `make proto-check`: passed.
- `make micro-smoke`: passed after Docker Desktop was available and the Compose stack was running.

## Error Response Contract

Existing HTTP API routes return a stable JSON envelope:

```json
{"code":0,"msg":"success","data":{}}
```

Expected failures keep HTTP transport successful and set `code` and `msg` in the envelope. Parameter and binding errors are handled in HTTP handlers. Business errors are returned by services and translated by handlers. Panics and other unexpected failures are captured by the shared base recover middleware and returned as the internal error envelope.

Route misses are still transport-level 404. This is why `/api/articles` can remain a 404 while `/api/content/articles` returns envelope-shaped content responses.

## Shared RPC Base Response

`idl/base/base.proto` defines `BaseResp` and project-local error codes. Auth and content RPC responses embed `base_resp`; handlers return nil RPC error for expected business failures and fill `base_resp` instead. Gateway auth validation reads `base_resp` and maps it into typed application errors before the HTTP middleware writes the unified envelope.

## Content Cache Boundary

Article detail caching is owned only by `content-service`:

```text
content-service ArticleService
  -> local article cache
  -> Redis article cache
  -> MySQL article repository
```

Gateway remains a thin router/proxy and does not add content cache, article lookup, or fallback business logic. Content writes invalidate article detail cache keys inside content-service.
