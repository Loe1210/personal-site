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
