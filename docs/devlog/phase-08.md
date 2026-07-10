# Phase 08 - Frontend First Release

## Goal

Implement the first complete frontend slice for the personal site without changing existing business APIs.

## Completed

- Added server-rendered page entry routes for:
  - `/`
  - `/blog`
  - `/blog/:slug`
  - `/about`
- Added template structure for:
  - shared layout
  - shared nav
  - home page
  - blog page
  - article detail page
  - about page
- Added frontend assets for:
  - global reset
  - theme system
  - page-specific styles
  - background motion
  - home/blog/article page scripts
  - avatar SVG
- Connected frontend pages to existing public APIs:
  - `/api/articles`
  - `/api/articles/:slug`
  - `/api/categories`
  - `/api/tags`
- Replaced the unreliable static directory serving path with a custom safe file-serving route for Windows runtime stability.

## Verification

- `go build ./...` passed.
- Temporary frontend verification server launched successfully on alternate ports during validation.
- Verified page shells:
  - Home page render OK
  - Blog page render OK
  - About page render OK
  - Article detail shell render OK
- Verified static assets:
  - theme CSS OK
  - blog JS OK
  - avatar SVG OK
- Verified API integration:
  - public article API reachable
- Current data note:
  - public article count was `0` during verification, so full live article-detail data rendering could not be validated against a real published slug in this round.

## Key Notes

- Frontend implementation reuses the reference site's visual direction, but does not import its build pipeline.
- Runtime validation exposed a Windows static asset issue when using the previous directory mount approach.
- The custom `/static/*filepath` handler is now the stable path for frontend assets in this project.

## Phase 08.1 - Public UI Redesign Refresh

### Goal

Replace the first public frontend release with the approved `Terminal Gallery` direction without changing backend APIs.

### Completed

- Reworked shared public layout with:
  - darker terminal-gallery atmosphere
  - deeper grid background
  - glowing accent pixels
  - stronger glass-panel system
- Reworked public navigation branding:
  - `LX` mark
  - `夏风` wordmark
  - unified top navigation style for blog/about/github
- Reworked homepage:
  - centered hero frame
  - stronger first-screen narrative
  - roadmap panel
  - refined latest-post cards
- Reworked blog list page:
  - stronger hero section
  - richer article card layout
  - integrated cover-image rendering when available
- Reworked article detail page:
  - reading-focused center column
  - left and right side information rails
  - refined content container and chips
- Reworked about page:
  - profile-driven intro block
  - roadmap, skills, internship timeline and site-positioning cards
- Updated frontend scripts so the redesigned cards still render live data from existing APIs.

### Verification

- `go build ./...` passed after the redesign changes.
- Started the local service successfully.
- Verified page responses:
  - `/` returned redesigned homepage shell
  - `/blog` returned redesigned blog shell
  - `/about` returned redesigned about shell
- Verified static asset response:
  - `/static/css/theme.css` returned `200`
- Current live data note:
  - public article API returned no published article slug during this verification round, so the redesigned article detail page was validated at the shell and script level, but not against a live published post in this pass.

### Notes

- This redesign only changed the public frontend layer.
- Existing backend APIs and admin frontend routes were not modified in this round.
- The approved visual direction is now `Terminal Gallery` and can be used as the basis for later admin visual unification.

## Phase 08.2 - Final Frontend Verification Pass

### Goal

Run the redesigned frontend on an isolated verification port and confirm the new public pages are being served from the latest code instead of a stale local process.

### Completed

- Added `APP_HOST_PORT` support in `main.go` so local verification can run on a custom port without conflicting with the default `:8888` process.
- Started the latest site successfully on `:8890`.
- Re-verified public page shells from the latest running instance:
  - `/`
  - `/blog`
  - `/about`
- Confirmed the latest public HTML includes the approved `Terminal Gallery` structure and assets.

### Notes

- Public article API still returned no published articles during this pass, so article detail rendering could not be validated against a real published slug.
- The custom verification port support is intentionally kept because it makes future UI review much easier.

## Phase 08.3 - Blog Real Filter Wiring

### Goal

Replace the previous client-side mock filtering flow on the Blog page with real URL-driven category and tag filtering backed by the public article API.

### Completed

- Updated the Blog page filter summary area so the current filter state is visible in the hero panel.
- Switched Blog list loading from "fetch all then filter in browser" to real API requests:
  - `/api/articles?category=...`
  - `/api/articles?tag=...`
  - combined category + tag query support
- Added URL state synchronization for Blog filters:
  - clicking chips updates the browser URL
  - refreshing the page preserves current filter state
  - browser back/forward navigation restores filter state
- Kept category and tag chip rendering driven by the live taxonomy APIs.

### Verification

- Started a temporary verification server on `:8892`.
- Created a temporary published article for verification.
- Verified filtered public article results:
  - category filter returned the expected article
  - tag filter returned the expected article
- Verified `/blog?category=go-backend&tag=hertz` returned the updated page shell with the live filter summary container.
- Deleted the temporary verification article after validation.
