# Blog Index Sidebar Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild `/blog` into a left-profile-sidebar plus right single-column article flow while preserving the site's dark mountain and glassmorphism style.

**Architecture:** Keep the existing static `/blog` page entrypoint and API flow, but replace the list-page DOM structure with a profile-first shell. Reuse real site assets and copy from the homepage, keep list data rendering in `static/blog/js/list.js`, and move the visual lift into `static/blog/css/blog.css` without disturbing article detail behaviors.

**Tech Stack:** Static HTML, vanilla JavaScript, CSS, Go test for lightweight UI regression checks

## Global Constraints

- Preserve the current dark mountain background, translucent glass panels, and cool blue accent language.
- Match the approved structure: fixed left personal sidebar, right content hero, single-column stacked article cards.
- Reuse existing site profile data where possible: `Hins`, homepage avatar asset, GitHub/email style.
- Keep existing blog data loading, category/tag filtering, pagination, and article link behavior working.
- Limit scope to the `/blog` article list page; do not redesign `/blog/post/:id` in this pass.

---

### Task 1: Lock the new blog shell with a regression test

**Files:**
- Create: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: `static/blog/index.html`, `static/blog/css/blog.css`, `static/blog/js/list.js`
- Produces: A lightweight regression test that checks the presence of the new sidebar, hero, and single-column article list hooks

- [ ] Add a Go test that reads the three front-end files and asserts the expected new layout markers exist.
- [ ] Run `go test ./...` to confirm the new assertions fail before implementation.

### Task 2: Rebuild the list page DOM and rendering hooks

**Files:**
- Modify: `static/blog/index.html`
- Modify: `static/blog/js/list.js`

**Interfaces:**
- Consumes: Existing `BlogAPI.getPosts`, `BlogAPI.getCategories`, `BlogAPI.getTags`, current article payload fields
- Produces: Left sidebar profile/navigation markup and single-column card rendering with new metadata and hero text regions

- [ ] Update the list page HTML to introduce the sidebar shell, profile block, hero block, and inline filter panels.
- [ ] Update list rendering so cards fit the new horizontal/article-stream layout while preserving links and filters.
- [ ] Keep category, tag, search, and pagination behavior intact after the DOM changes.

### Task 3: Apply the approved visual system

**Files:**
- Modify: `static/blog/css/blog.css`

**Interfaces:**
- Consumes: New class hooks from `static/blog/index.html` and `static/blog/js/list.js`
- Produces: Approved dark cinematic sidebar layout, glass panels, single-column cards, and responsive behavior

- [ ] Replace the old two-column list-page styling with the new left-sidebar and right-content layout.
- [ ] Style the profile section, compact stats, nav items, hero area, inline filter panels, and article cards.
- [ ] Add responsive rules that collapse the sidebar cleanly on smaller screens while keeping cards readable.

### Task 4: Verify and polish

**Files:**
- Modify: `blog_index_ui_test.go` if assertions need tightening after implementation

**Interfaces:**
- Consumes: Updated static assets
- Produces: Passing regression checks and final verification notes

- [ ] Run `go test ./...` again and confirm the new test passes.
- [ ] Manually review the resulting markup/CSS consistency and remove any stale selectors or incompatible assumptions from the old layout.
