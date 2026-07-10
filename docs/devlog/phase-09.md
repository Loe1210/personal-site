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

## Phase 09.1 - Admin UI Redesign Refresh

### Goal

Bring the admin console into the same `Terminal Gallery` design system as the public frontend, without changing existing admin APIs or scripts.

### Completed

- Reworked shared admin shell:
  - admin layout now reuses the same background atmosphere as the public site
  - redesigned admin sidebar branding and system status block
  - redesigned topbar with stronger page identity and user status area
- Reworked admin login page:
  - two-column hero + login card structure
  - terminal-gallery visual tone
  - preserved existing login form ids and JS behavior
- Reworked admin dashboard page:
  - stronger welcome panel
  - system status box
  - refined stat cards and quick action cards
- Reworked admin articles page:
  - richer toolbar section
  - preserved list container, search input and filter ids
- Reworked admin article editor page:
  - clearer writing surface vs publishing settings split
  - preserved all form ids and upload hooks
- Reworked admin taxonomy page:
  - stronger category/tag management presentation
  - preserved existing form and list ids
- Reworked shared admin stylesheet so the admin console now visually matches the approved public direction.

### Verification

- `go build ./...` passed after the admin redesign changes.
- Updated admin templates on disk were confirmed successfully:
  - admin login template
  - admin dashboard template
  - admin sidebar template
- Runtime validation note:
  - local port `8888` was already occupied by an older running instance during this verification round
  - requests sent to `http://127.0.0.1:8888/admin/...` returned the old shell from that already-running process
  - this means the new admin HTML was written correctly and the project builds, but final visual runtime confirmation requires restarting the running local server instance so it loads the new templates

### Notes

- This redesign preserved the current admin JS contract by keeping the important DOM ids unchanged.
- Backend routes and session-based admin APIs were not changed in this round.
- The admin console now has a frontend structure aligned with the public `Terminal Gallery` system.
