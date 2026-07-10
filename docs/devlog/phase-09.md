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

## Phase 09.2 - Final Admin Polish Pass

### Goal

Polish the last admin presentation details after the redesign and verify the newest admin UI from a fresh isolated local instance.

### Completed

- Added page-specific visible admin headings and descriptions in `biz/site/handler.go` so browser titles and on-page titles are no longer forced to be identical.
- Updated the shared admin topbar template to use polished visible headings.
- Rebuilt and restarted the latest code on verification port `:8890`.
- Re-verified latest admin page shells from the running `:8890` instance:
  - `/admin/login`
  - `/admin`
  - `/admin/articles/new`
  - `/admin/taxonomy`
- Confirmed the new admin shell, title treatment, sidebar, and page descriptions are now served by the latest running code.

### Notes

- This pass was a visual polish pass, not a behavior rewrite.
- Existing admin DOM ids and JS contracts remain intact.
- `APP_HOST_PORT` support now also helps future admin visual review without disturbing the default local service.
## Phase 09.3 - Taxonomy Management Completion

### Goal

Finish the missing category/tag maintenance loop in the admin console so taxonomy is no longer create-only.

### Completed

- Added backend update/delete APIs for category:
  - `PUT /api/admin/categories/:id`
  - `DELETE /api/admin/categories/:id`
- Added backend update/delete APIs for tag:
  - `PUT /api/admin/tags/:id`
  - `DELETE /api/admin/tags/:id`
- Added conflict and usage handling:
  - duplicate category or tag slug/name returns conflict-style business error
  - deleting a category that is still referenced by articles is blocked
  - deleting a tag that is still referenced by `article_tags` is blocked
- Extended RBAC seed permissions:
  - `category:update`
  - `category:delete`
  - `tag:update`
  - `tag:delete`
- Updated admin taxonomy page:
  - added edit action
  - added delete action
  - added cancel-edit flow
  - preserved existing admin page design language
- Updated taxonomy frontend script:
  - create and update now share the same form
  - editing state is visible in-page
  - delete uses guarded confirmation
  - success feedback remains visible after refresh instead of being cleared immediately

### Validation

- `go build ./...` passed.
- Verified the full taxonomy flow on an isolated local instance:
  - login
  - create category
  - update category
  - delete category
  - create tag
  - update tag
  - delete tag

### Notes

- `hz model` was not available in the current toolchain path during this round, so the thrift contracts were updated first and small supplemental request/response structs were added under:
  - `biz/model/category/custom.go`
  - `biz/model/tag/custom.go`
- After the local `hz` toolchain path is restored, these temporary supplemental structs can be replaced by regenerated model code.

## Phase 09.4 - Admin Article Editor Encoding and Layout Fix

### Goal

Fix the broken admin article editor page where heading copy became garbled, the category selector was hard to read, and the right-side publishing panel overflowed visually.

### Completed

- Repaired invalid page copy source in `biz/site/handler.go`.
- Rewrote the repaired handler and article editor template files explicitly as UTF-8 to remove replacement-character rendering.
- Improved the article editor side-panel layout:
  - wider desktop publishing column
  - single-column tag list to avoid chip overflow
  - safer wrapping for long tag names
- Improved category selector readability by styling option foreground/background colors.
- Improved the upload area presentation:
  - added selected-file name display
  - styled native file input
  - kept existing upload behavior intact

### Verification

- `go build ./...` passed.
- Fresh verification server started on `:8894` after the UTF-8 rewrite.
- Verified the new article page HTML contains:
  - `文章主体`
  - `支持 JPG、JPEG、PNG、WEBP、GIF，单文件不超过 5MB。`
- Verified no replacement character remained in the fresh page response.
- Verified admin category API returned category data.
- Verified a real draft article could be created with `category_id = 1` and then deleted successfully.
