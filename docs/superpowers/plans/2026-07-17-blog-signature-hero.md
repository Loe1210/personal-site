# Blog Signature Hero Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade the blog home hero into a SimonAKing-inspired signature header, improve hover breathing interactions, and add a cute site cursor.

**Architecture:** Keep the existing HTML structure and append a focused CSS override section to avoid destabilizing prior layout work. Add one global cursor stylesheet that can be linked from public pages.

**Tech Stack:** Static HTML, CSS, existing Go UI tests.

## Global Constraints

- Preserve the dark mountain/fog visual direction.
- Do not remove the pet runtime.
- Use cache-busted CSS versions so Docker/nginx browsers fetch the new assets.

---

### Task 1: Tests And Copy Markers

**Files:**
- Modify: `blog_index_ui_test.go`

**Steps:**
- [ ] Add assertions for `STDIN | Think >> /dev/Mind`, `CODE / MICROSERVICE / LEARNING NOTES`, `cute-cursor.css`, and `blog.css?v=15`.
- [ ] Run `go test .` and expect failure before implementation.

### Task 2: Hero Copy And Cursor Asset

**Files:**
- Modify: `static/blog/index.html`
- Modify: all public static pages that should load cursor CSS
- Create: `static/assets/css/cute-cursor.css`

**Steps:**
- [ ] Replace hero copy with the approved English signature copy.
- [ ] Add a cute cursor stylesheet using inline SVG cursor data URLs.
- [ ] Update blog stylesheet links to `blog.css?v=15`.

### Task 3: Visual Breathing And Hover Layer

**Files:**
- Modify: `static/blog/css/blog.css`

**Steps:**
- [ ] Append a dedicated signature hero override section.
- [ ] Add hover transitions for title, subtitle, search, cards, media, links, tags, nav, and sidebar profile actions.
- [ ] Respect reduced motion.
- [ ] Run `go test .`, rebuild local Docker, and verify local HTTP resources.