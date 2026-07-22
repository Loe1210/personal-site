# Microservice Deep Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the project into a cleaner learning-oriented microservice system where the frontend is deployed as static assets, `gateway` is the single API entrypoint, and domain services communicate through Kitex RPC, Nacos, and OpenTelemetry.

**Architecture:** Remove `web-bff` from the main runtime path. Browser requests load static assets from Nginx and call `/api/*` on `gateway`. `gateway` calls `auth-service`, `content-service`, and `media-service` through Kitex RPC with Nacos service discovery. Each domain service owns its own data; there is no generic database service.

**Tech Stack:** Go, Hertz, Kitex, protobuf, Nacos, OpenTelemetry, MySQL, Redis, Docker Compose, Kubernetes, Nginx, npm/static frontend build, Prometheus, Grafana, Jaeger.

## Global Constraints

- Work on branch `feature/microservice-deep-refactor`.
- Do not add a `database-service`; keep data ownership inside `auth-service`, `content-service`, and `media-service`.
- Do not split `article`, `category`, and `tag` into separate services in this phase.
- Remove `web-bff` from the primary request path; keep it only temporarily while migration is incomplete, then delete it.
- Use existing proto contracts first. Expand proto only when a concrete caller needs a missing RPC.
- Keep generated code under `kitex_gen/` read-only; regenerate from `idl/**/*.proto`.
- Keep HTTP DTOs and RPC DTOs separate, with explicit mapping code.
- Every task must leave the project testable with focused tests and, when applicable, `go test ./...`, `make proto-check`, and `make micro-smoke`.

---

## Task List

### Task 1: Capture Baseline and Current Architecture

Create `docs/architecture/microservices-current-state.md`, record current HTTP reverse proxy/BFF state, record target gateway-centered topology, run baseline verification, and commit.

### Task 2: Establish Frontend-to-Gateway Contract

Confirm whether the frontend is npm-built or static-only, ensure Nginx serves static assets and proxies only `/api/*` to `gateway`, and update README.

### Task 3: Add Gateway Content Kitex Client

Create `services/gateway/internal/client/content/` with gateway DTOs, `ArticleClient` interface, Kitex client wrapper, and mapping tests using existing `content.proto`.

### Task 4: Replace Gateway Content Reverse Proxy With Handlers

Create gateway HTTP handlers for `GET /api/articles` and `GET /api/articles/:id`, call the content Kitex client, keep legacy `/api/content/*` only as a temporary compatibility path, and test handlers/router.

### Task 5: Verify Content Service RPC Behavior

Add focused tests for `services/content-service/internal/handler/rpc/content.go`, harden article list/detail error handling and DTO mapping, and keep using the existing content proto.

### Task 6: Add Nacos Registry and Resolver for Kitex

Add service-local Nacos registry/resolver helpers, register `content-service`, let `gateway` discover `content-service`, and keep direct RPC address fallback for local development.

### Task 7: Remove Web BFF From Main Runtime Path

Remove gateway routes and compose dependencies for `web-bff`, migrate frontend API calls to gateway endpoints, then delete `services/web-bff` only after references are gone and tests pass.

### Task 8: Use Existing Auth Proto for Gateway Authentication

Create gateway auth RPC client around existing `ValidateSession` and `CheckPermission`, wire gateway auth middleware through RPC, and add auth-service Nacos registration.

### Task 9: Introduce Media RPC Only for Concrete Metadata Needs

Keep streaming upload over HTTP through gateway upload guard. Use existing `GetFile` RPC only if gateway needs metadata; expand media proto only with a same-task concrete caller.

### Task 10: Standardize Errors, Response Mapping, and Logging

Create gateway-local error response types/codes, replace ad hoc gateway JSON errors, and ensure RPC errors map consistently to HTTP responses.

### Task 11: Complete OpenTelemetry Propagation Across Sample Path

Propagate trace context across `frontend -> gateway -> content-service`, add gateway/content spans, configure Kitex tracing middleware, and write a tracing runbook.

### Task 12: Promote Nacos Configuration Gradually

Define config precedence as `env > Nacos > yaml > defaults`, add service name/environment/Nacos/RPC config fields, and document the rollout.

### Task 13: Align Docker Compose With RPC, Nacos, and OTEL

Expose service RPC ports, remove `web-bff` from compose, set gateway service discovery variables, and extend smoke tests for the new path.

### Task 14: Prepare Kubernetes Manifests

Retire `web-bff` manifests, align Service/Deployment ports, add probes, and move runtime config to ConfigMap/Secret.

### Task 15: Final Documentation and Acceptance

Write the final architecture document, update README and pending tasks, run `go test ./...`, `make proto-check`, `make micro-smoke`, and `kubectl kustomize deploy/k8s/base`.

---

## First Execution Slice

This run starts with Task 1 only:

1. Create `docs/architecture/microservices-current-state.md`.
2. Run baseline verification.
3. Review failures as pre-existing or blockers.
4. Commit the baseline docs if verification reaches a usable state.

## Self-Review

- The plan reflects the user's decision to remove `web-bff` from the main path.
- Proto expansion is constrained to concrete RPC callers.
- The first sample path is `frontend -> gateway -> content-service`.
- The plan does not introduce a database service.
