# Phase 14 - Blog UI Refresh, Pet Runtime, and Upload Hardening

## Date

2026-07-20

## Summary

This phase prepares the current local version for production deployment. The work combines the blog visual refresh, the global animated pet, admin cover selection, and a safer chunked upload path for article covers.

## Changes

- Refreshed the blog home experience with a fixed side navigation, softer article cards, category/tag split pages, breathing room, page fade-in, hover interactions, and a custom cursor.
- Added the global pet runtime across pages, including double-click pet switching and shared static pet assets under `pet/`.
- Added admin cover selection from a remote image library while keeping local upload as a second option.
- Switched admin cover upload to the backend chunked upload API and kept end-to-end SHA-256 verification.
- Fixed media chunk handling by reading the buffered Hertz request body before writing chunk files.
- Fixed uploaded merged files and generated thumbnails so nginx can read them after upload.
- Documented safe deployment commands and ignored local database backups/deployment archives to avoid accidental commits.

## Data Safety

The production deployment for this phase must be code-only. It must not run SQL imports, must not run `docker compose down -v`, and must not delete or recreate the MySQL volume. Existing uploaded images are preserved through Docker volumes; only file permissions may be normalized to make existing media readable by nginx.

## Verification

Before deployment, run:

```bash
go test . -count=1
go test ./services/media-service/biz/upload ./services/media-service/internal/service -count=1
node --check static/admin/js/admin.js
```

After deployment, verify:

```bash
curl -I --max-time 10 http://117.72.95.156:8080/blog/
curl -s --max-time 10 http://117.72.95.156:8080/blog/ | grep -E 'blog.css|pet.js'
curl -s --max-time 10 http://117.72.95.156:8080/admin/ | grep 'admin.js?v=13'
```