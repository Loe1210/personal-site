# Frontend to Gateway Contract

## Current Frontend Shape

The project currently has no `package.json` at the repository root or under `frontend/`.

The frontend is therefore treated as a static asset deployment:

- `static/` contains the site pages, scripts, and styles.
- `pet/` contains pet assets.
- `frontend/nginx/default.conf` serves static files through Nginx.
- `frontend/Dockerfile` packages the static frontend container.

If the frontend later becomes a package-managed app, npm build should produce static files that Nginx can serve. Until then, do not introduce npm only for the microservice refactor.

## Current Contract

Browser requests use two categories of paths:

- static page and asset paths served by Nginx,
- API paths under `/api/*` forwarded to `gateway`.

Current flow:

```text
Browser
  -> frontend Nginx
  -> /api/* proxy to gateway
  -> gateway cross-cutting middleware and thin proxy routing
  -> domain services
```

## Content API Paths

Public blog pages now call content-domain APIs through explicit `/api/content/*` paths:

- `/api/content/articles`
- `/api/content/articles/:id`
- `/api/content/articles/:id/adjacent`
- `/api/content/categories`
- `/api/content/tags`

Nginx no longer rewrites `/api/articles`, `/api/categories`, or `/api/tags` to content paths. This keeps the frontend contract honest and lets gateway remain a thin routing layer instead of owning content-specific handlers.

## Admin API Paths

Admin pages keep the existing browser-facing `/api/admin/*` paths. Nginx rewrites those paths to gateway's protected content-admin proxy:

- `/api/admin/*` -> `/api/content/admin/*`

Gateway applies session validation before forwarding admin content traffic to `content-service`. Do not bypass gateway for content admin routes unless auth responsibility is moved into `content-service` first.

## Generic Nginx Proxy

All other API traffic goes through the generic proxy:

```nginx
location /api/ {
    proxy_pass http://gateway:8888/api/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

Gateway, not Nginx, owns cross-cutting checks such as auth and upload protection. Content-domain business behavior stays in `content-service`.
