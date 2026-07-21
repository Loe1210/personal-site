# Phase 15 - Blog Directory UI, Interaction Polish, And Pet Consolidation

## Date

2026-07-21

## Scope

This phase completes the blog experience refresh across `/blog/`, `/blog/categories`, `/blog/tags`, and article detail pages. The work is frontend-only: it preserves the existing public APIs and article data model while rebuilding presentation, navigation, interaction, and static deployment assets.

## Blog, Category, And Tag Experience

- Unified the desktop sidebar shell across blog, category, and tag pages, including profile content, statistics, navigation, social actions, avatar feedback, and the bottom pet dock.
- Rebuilt the category directory as a grouped article timeline. Category links update the focus URL and smoothly scroll to the matched group.
- Rebuilt the tag page as `Tag.sort()` with an interactive constellation and a focused archive view. Selecting a tag renders its linked articles and article links use numeric IDs accepted by the detail route.
- Updated the category heading to `Category.sort()` and removed directory hero divider lines.
- Added a sliding sidebar navigation indicator and normalized active navigation state per route.

## Article Stream And Detail Navigation

- Changed the blog list to text-first dark article cards; cover images remain available on article detail pages.
- Replaced numbered pagination with incremental loading through `IntersectionObserver`, with a scroll fallback.
- Moved category and tag labels into the card body and kept date plus reading time in the footer.
- Preserved numeric article detail URLs so tag archive items resolve correctly instead of redirecting back to the list.

## Visual System And Motion

- Added `static/assets/img/blog-background.png`, copied from the approved background source, and applied it to all blog routes including article detail.
- Reworked the header canopy into a symmetric `60 x 16` breathing dot grid. Each dot uses a staggered brightness cycle while the header remains visually continuous with the page background.
- Centered the blog, category, and tag hero titles with staged entrance motion.
- Added granular hover feedback for cards, labels, titles, summaries, navigation, social actions, profile avatar, and sidebar text.
- Sidebar characters now behave independently: a hovered character lifts and stays at its peak, then plays a decreasing free-fall-style rebound sequence after pointer leave.
- Placed the blog search control in the desktop header's right-top corner, with a responsive in-flow layout on narrow screens.

## Pet Assets

- Repaired `pet/index.json` as valid JSON.
- Reduced the pet manifest to the single default pet `yuexinmiao` (月薪喵).
- Removed the `zhangfei-tusun` metadata and spritesheet resources.

## Verification

The following checks passed before release:

```bash
node --check static/blog/js/sidebar.js
go test ./...
docker compose -f deploy/docker/compose.yaml up -d --build --no-deps frontend
```

Local deployment verification confirmed that `/blog/` serves the latest versioned blog stylesheet, sidebar interaction script, header-star script, and the salary-cat-only pet manifest.

## Production Deployment

Static-only production release uses the existing project command:

```bash
make deploy-frontend
```

It packages `static/`, `pet/`, frontend Nginx configuration, and the UI regression test file; uploads them to the configured server; rebuilds the remote `frontend` container; then verifies `/blog/`. It does not run database migrations or remove MySQL volumes.