# Microservices Current State

## Services

- `frontend`: Nginx-served static site and static asset entrypoint.
- `gateway`: external HTTP entrypoint. It currently routes API traffic mostly through HTTP reverse proxy handlers.
- `web-bff`: current page-facing aggregation layer. The target architecture removes it from the main runtime path.
- `auth-service`: login, session, user, role, and permission domain.
- `content-service`: article, category, and tag domain.
- `media-service`: upload task, uploaded file metadata, and local file storage domain.

## Current Communication

- Browser traffic enters the frontend container first.
- Frontend API traffic is expected to flow through `gateway`.
- `gateway` currently proxies `/api/auth/*`, `/api/media/*`, `/api/content/*`, and `/api/blog/*` with HTTP reverse proxy logic.
- `web-bff` currently calls `content-service` through HTTP for blog page aggregation.
- `auth-service`, `content-service`, and `media-service` already start Kitex RPC servers.
- Runtime callers are not yet using Kitex RPC as the primary service-to-service path.

## Current Contracts

- `idl/auth/auth.proto` exposes `ValidateSession` and `CheckPermission`.
- `idl/content/content.proto` exposes article list, detail, create, update, and delete RPC methods.
- `idl/media/media.proto` exposes `GetFile`.
- Existing proto contracts should be used before adding new methods.
- New proto methods should be added only when a concrete gateway or service caller needs them in the same implementation slice.

## Target Change

- Remove `web-bff` from the primary runtime path.
- Keep frontend deployment as static assets served by Nginx, with optional npm build if the frontend source grows into a package-managed app.
- Route all browser API calls through `gateway`.
- Move gateway-to-domain-service calls from HTTP reverse proxying toward `Kitex RPC + Nacos service discovery`.
- Use `gateway -> content-service` as the first sample path.
- Use the existing auth proto for gateway session and permission checks before expanding auth contracts.
- Keep media upload streaming over HTTP unless a concrete metadata RPC need appears.

## Target Runtime Topology

```text
Browser
  -> frontend Nginx static site
  -> /api/* forwarded to gateway
  -> gateway Hertz HTTP handlers and middleware
  -> auth-service Kitex RPC
  -> content-service Kitex RPC
  -> media-service Kitex RPC for metadata, HTTP for upload streaming when appropriate
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
- `make micro-smoke`: failed because gateway was not reachable at `http://127.0.0.1:8888/healthz`.
- `make micro-up`: failed because Docker Desktop was not running; Docker could not connect to `npipe:////./pipe/dockerDesktopLinuxEngine`.

The smoke failure is currently an environment limitation, not a confirmed application failure.
