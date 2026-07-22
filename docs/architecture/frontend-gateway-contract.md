# Frontend to Gateway Contract

## Current Frontend Shape

The project currently has no `package.json` at the repository root or under `frontend/`.

The frontend is therefore treated as a static asset deployment:

- `static/` contains the site pages, scripts, and styles.
- `pet/` contains pet assets.
- `frontend/nginx/default.conf` serves static files through Nginx.
- `frontend/Dockerfile` packages the static frontend container.

If the frontend later becomes a package-managed app, npm build should produce static files that Nginx can serve. Until then, do not introduce npm only for the microservice refactor.

## Target Contract

Browser requests use two categories of paths:

- static page and asset paths served by Nginx,
- API paths under `/api/*` forwarded to `gateway`.

Target flow:

```text
Browser
  -> frontend Nginx
  -> /api/* proxy to gateway
  -> gateway handlers and middleware
  -> domain services through Kitex RPC
```

## Compatibility Note

`frontend/nginx/default.conf` currently has compatibility locations that rewrite public frontend API paths to gateway's existing `/api/content/*` reverse proxy paths. Examples:

- `/api/articles` -> `/api/content/articles`
- `/api/categories` -> `/api/content/categories`
- `/api/tags` -> `/api/content/tags`

Do not remove these rewrites until `gateway` owns first-class handlers for the public article/category/tag endpoints.

## Migration Target

After gateway content handlers are implemented, Nginx should be simplified so API traffic is forwarded without content-service-specific rewrites:

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

Gateway, not Nginx, should own backend service routing.
