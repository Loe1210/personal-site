# Phase 16 - Gateway-Centered Microservice Runtime

## Date

2026-07-22

## Scope

This phase starts the deeper microservice refactor by making `gateway` the single browser-facing backend and moving service-to-service reads/auth checks onto Kitex RPC. It removes `web-bff` from the main runtime path, keeps the static frontend as an npm-built/Nginx-served asset, and introduces Nacos-backed service discovery for the first gateway RPC clients.

The intent is not to complete every business boundary in one pass. This phase establishes the runnable migration shape: frontend calls gateway, gateway calls business services through RPC, and Docker Compose can start and smoke-test the stack.

## Gateway And RPC Integration

- Added a gateway content RPC client backed by Kitex.
- Added first-class article read handlers in gateway:
  - `GET /api/articles`
  - `GET /api/articles/:id`
- Added a gateway auth RPC client backed by Kitex.
- Changed gateway session validation to call `auth-service` through RPC.
- Protected `/api/content/admin/*path` at the gateway layer with the auth middleware.
- Kept `/api/content/*path` as a compatibility proxy while remaining content write/admin routes are migrated.

## Service Discovery

- Added Nacos registry support for `auth-service`.
- Added Nacos registry support for `content-service`.
- Added Nacos resolver support for gateway Kitex clients.
- Gateway now prefers Nacos service discovery when `NACOS_ADDR` is configured.
- Direct RPC addresses remain as local-development and troubleshooting fallbacks:
  - `AUTH_RPC_ADDR`
  - `CONTENT_RPC_ADDR`

## Runtime And Deployment Shape

- Removed `services/web-bff` from the main runtime path.
- Removed `web-bff` from Docker Compose and deployment service lists.
- Removed the gateway `/api/blog` reverse proxy route that depended on `web-bff`.
- Kept the frontend as static assets served by Nginx, with API traffic routed to gateway.
- Exposed RPC ports in Docker Compose:
  - `auth-service`: `9101`
  - `content-service`: `9103`

## Documentation

- Added the gateway-centered implementation plan:
  - `docs/superpowers/plans/2026-07-22-microservice-deep-refactor-gateway-centered.md`
- Added the current microservice architecture snapshot:
  - `docs/architecture/microservices-current-state.md`
- Added the frontend-to-gateway contract:
  - `docs/architecture/frontend-gateway-contract.md`
- Updated local runbook, README, and pending task notes for the new runtime shape.

## Verification

The following checks passed before release:

```bash
go test ./...
make proto-check
docker compose -f deploy/docker/compose.yaml ps
make micro-smoke
```

The smoke test verified:

- gateway health
- anonymous auth `/me`
- content article list
- login cookie flow

## Remaining Work

- Migrate `media-service` gateway access to Kitex RPC.
- Replace remaining `/api/content/*path` compatibility proxy routes with gateway handlers and RPC clients.
- Simplify frontend Nginx API rewrites after compatibility paths are removed.
- Decide whether to split database ownership further by schema, service database, or a dedicated data service.
- Extend proto contracts only when a gateway handler or service-to-service workflow actually needs the interface.
- Add deeper production governance later: rate limits, circuit breaking, unified error codes, tracing dashboards, and CI/CD gates.

## Production Deployment

Code production release uses the existing project command with the active branch:

```bash
make deploy-code BRANCH=feature/microservice-deep-refactor
```

This fetches the branch on the configured server, fast-forwards the remote working tree, and rebuilds the frontend, media-service, content-service, and gateway containers.
