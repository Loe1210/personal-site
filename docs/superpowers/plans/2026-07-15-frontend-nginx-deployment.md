# Frontend Nginx Deployment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** Add a standalone Dockerized frontend served by Nginx, with `/api/*` proxied to the existing gateway service.

**Architecture:** The browser enters through a new `frontend` container exposed on host port `8080`. Nginx serves the existing static frontend files and forwards API traffic to `http://gateway:8888` over the Docker network. Existing backend services and the gateway routing model remain unchanged.

**Tech Stack:** Docker Compose, Nginx, static HTML/CSS/JS, existing Go microservices.

## Global Constraints

- Frontend browser entrypoint is `http://127.0.0.1:8080`.
- Backend API requests from the browser must use the same origin path prefix `/api`.
- Nginx must proxy `/api/*` to `gateway:8888`.
- Existing backend Docker Compose services must continue to start with `make micro-up`.
- Do not change the existing backend service split or API gateway behavior for this task.

---

## File Structure

- Create `frontend/Dockerfile`: builds the frontend image from Nginx and copies static assets into the Nginx web root.
- Create `frontend/nginx/default.conf`: contains static file serving rules and `/api/` reverse proxy rules.
- Modify `deploy/docker/compose.yaml`: adds the `frontend` service and exposes it as `8080:80`.
- Modify `README.md`: documents the new browser entrypoint and keeps gateway as the internal/API entry.

---

### Task 1: Add Nginx Frontend Container

**Files:**
- Create: `frontend/Dockerfile`
- Create: `frontend/nginx/default.conf`

**Interfaces:**
- Consumes: existing `static/` directory at build time.
- Produces: Docker image serving static files on container port `80`.

- [x] **Step 1: Create the frontend Dockerfile**

Create `frontend/Dockerfile`:

```dockerfile
FROM nginx:1.27-alpine

COPY frontend/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY static/ /usr/share/nginx/html/
```

- [x] **Step 2: Create the Nginx config**

Create `frontend/nginx/default.conf`:

```nginx
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    location /api/ {
        proxy_pass http://gateway:8888/api/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

- [x] **Step 3: Review config syntax locally**

Run:

```powershell
docker run --rm -v "${PWD}/frontend/nginx/default.conf:/etc/nginx/conf.d/default.conf:ro" nginx:1.27-alpine nginx -t
```

Expected: Nginx reports configuration syntax is ok.

---

### Task 2: Wire Frontend Into Docker Compose

**Files:**
- Modify: `deploy/docker/compose.yaml`

**Interfaces:**
- Consumes: frontend image build context from the repository root.
- Produces: `frontend` service exposed at `http://127.0.0.1:8080`.

- [x] **Step 1: Add the compose service**

Add this service under `services`:

```yaml
  frontend:
    build:
      context: ../..
      dockerfile: frontend/Dockerfile
    depends_on:
      gateway:
        condition: service_started
    ports:
      - "8080:80"
```

- [x] **Step 2: Start the stack**

Run:

```powershell
make micro-up
```

Expected: Docker Compose builds and starts the new `frontend` service alongside the backend.

- [x] **Step 3: Check frontend container**

Run:

```powershell
docker compose -f deploy/docker/compose.yaml ps frontend
```

Expected: `frontend` is running and has `0.0.0.0:8080->80/tcp`.

---

### Task 3: Document Access Path

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: the final browser entrypoint from Task 2.
- Produces: clear local run instructions for the user.

- [x] **Step 1: Update local run section**

Document:

```text
浏览器访问前端页面：
http://127.0.0.1:8080

后端 API 网关仍然保留：
http://127.0.0.1:8888
```

- [x] **Step 2: Keep Docker Compose command unchanged**

The existing startup command remains:

```bash
make micro-up
```

---

### Task 4: Verify End-To-End

**Files:**
- No code changes.

**Interfaces:**
- Consumes: running Docker Compose stack.
- Produces: evidence that the frontend and API proxy work.

- [x] **Step 1: Run Go tests**

Run:

```powershell
go test ./...
```

Expected: all packages pass.

- [x] **Step 2: Run backend smoke test**

Run:

```powershell
make micro-smoke
```

Expected: gateway, content, auth/session smoke checks pass.

- [x] **Step 3: Check frontend HTTP response**

Run:

```powershell
curl.exe -I http://127.0.0.1:8080
```

Expected: `HTTP/1.1 200 OK`.

- [x] **Step 4: Check API proxy response**

Run:

```powershell
curl.exe -I http://127.0.0.1:8080/api/content/articles
```

Expected: response comes from the backend through the frontend Nginx proxy.
