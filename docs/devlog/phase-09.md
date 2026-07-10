# Phase 09 - Admin Frontend First Release

## Goal

Build the first usable admin console for the personal site so the existing backend capabilities can be used from browser pages instead of Swagger only.

## Completed

- Added admin plan document:
  - `docs/personal-site-admin-frontend-plan.md`
- Added admin page routes:
  - `/admin/login`
  - `/admin`
  - `/admin/articles`
  - `/admin/articles/new`
  - `/admin/articles/:id/edit`
  - `/admin/taxonomy`
- Added admin templates:
  - shared admin layout
  - admin sidebar
  - admin topbar
  - login page
  - dashboard page
  - article list page
  - article editor page
  - taxonomy page
- Added admin frontend assets:
  - shared admin stylesheet
  - login stylesheet
  - admin common JS
  - login JS
  - dashboard JS
  - article list JS
  - article editor JS
  - taxonomy JS
- Reused the existing session-based auth APIs and admin CRUD APIs.
- Fixed `/api/admin/me` so unauthenticated requests safely return `401` instead of relying on unsafe session type assertions.

## Validation

### Build

- `go build ./...` passed.

### Page shell checks

Validated successfully on a temporary verification server:

- admin login page shell
- admin dashboard page shell
- admin articles page shell
- admin taxonomy page shell
- admin stylesheet asset

### Session checks

Validated successfully:

- unauthenticated `GET /api/admin/me` returns `401`
- admin login works with seeded account `admin / admin`
- authenticated `GET /api/admin/me` returns current user
- logout works

### Authenticated admin CRUD checks

Validated successfully in one end-to-end flow:

- create category
- create tag
- upload article cover
- create article
- update article
- delete article
- logout

## Key Notes

- The admin frontend still uses the current article list API to prefill the edit page by selecting the matching article id on the client side.
- Upload is integrated into the article editor instead of a dedicated admin upload page in this phase.
- The admin console follows the existing session architecture and does not introduce token storage in the browser.
