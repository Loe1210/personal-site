# Microservice Error Cache Hardening Design

## Goal

Unify business response semantics across gateway-facing HTTP handlers and internal RPC handlers, harden auth RPC error mapping, keep content caching inside content-service, and expand smoke coverage for the new content route surface.

## Error Semantics

All expected request and business failures use the same response shape:

```json
{"code":0,"msg":"success","data":{}}
```

For errors, HTTP handlers still return the same envelope with HTTP 200 when the route exists:

```json
{"code":20010001,"msg":"invalid article id","data":null}
```

Handler code owns parameter and binding errors. Service code owns business errors such as not found, invalid credentials, expired session, or missing permissions. Handlers translate service errors into the shared envelope. Unexpected failures are not manually invented in handlers; global recover middleware catches panics and returns the same envelope with the internal error code.

Route misses remain transport-level 404, so the deprecated `/api/articles` route can smoke-test as 404 while `/api/content/*` business failures stay envelope-based.

## Shared Proto

Create `idl/base/base.proto` with a small `BaseResp` message and project-local `ErrorCode` enum. Auth and content RPC responses embed `base.BaseResp base_resp = ...`. RPC handlers return nil RPC error for expected business failures and fill `base_resp`; only transport/setup failures remain RPC errors.

## Auth RPC Mapping

Auth service returns `base_resp` for invalid session and permission business failures. Gateway auth client inspects `base_resp` and maps it to typed domain errors. Gateway middleware converts auth business errors to the unified envelope rather than HTTP 401, while RPC transport errors become an upstream-unavailable envelope.

## Content Cache

Only content-service owns content caching. `ArticleService.GetArticleByID` reads local cache first, then Redis, then MySQL repository. Redis and local cache misses fall through; repository hits backfill Redis and local. Cache backend failures do not become business errors when MySQL can answer. Mutations invalidate article detail cache keys.

## Panic Recovery

Each HTTP service installs a base recover middleware that returns the unified envelope for panics. A shared gopool helper registers panic handling for pooled tasks. Direct goroutines should use a `defer` recover helper at the goroutine boundary.

## Smoke Coverage

`make micro-smoke` checks:

- gateway `/healthz`
- content-service direct `/articles`
- gateway `/api/content/articles`
- gateway deprecated `/api/articles` as transport 404
- auth login and `/me` envelope behavior

## Documentation

Update devlog and architecture docs with the response envelope, error ownership, auth RPC mapping, and content-service-only cache boundary.
