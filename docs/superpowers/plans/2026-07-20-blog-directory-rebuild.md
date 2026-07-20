# Blog Directory Rebuild Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild the category and tag pages from the blog page shell, with route-correct sidebar state, a reference-matched category directory, and an interactive tag constellation.

**Architecture:** Keep the existing BlogAPI interface and blog background loader. Introduce one small sidebar-state helper shared by the two page scripts, then render each page only into its own main-content mount point. CSS is scoped to the new page classes so the blog home page remains unchanged.

**Tech Stack:** Static HTML, vanilla JavaScript, CSS, Nginx local container deployment.

## Global Constraints

- Preserve `/api/categories`, `/api/tags`, and `/api/articles` contracts.
- Use the blog page sidebar markup and copy unchanged, changing only the active navigation item per route.
- Keep focus links functional through `?focus=<name>`.
- Use ASCII-only source comments and escaped HTML for API-derived strings.

---

### Task 1: Define DOM and navigation regression checks

**Files:**
- Create: `blog_directory_ui_test.go`
- Modify: `static/blog/categories.html`
- Modify: `static/blog/tags.html`

**Interfaces:**
- Consumes: static page HTML served by Nginx.
- Produces: checks that assert each page has exactly one active sidebar entry and the correct main mount point.

- [ ] **Step 1: Write failing checks**

```go
func TestDirectoryPagesUseTheirOwnActiveNavigation(t *testing.T) {
    assertActiveNav(t, "static/blog/categories.html", "分类")
    assertActiveNav(t, "static/blog/tags.html", "标签")
}
```

- [ ] **Step 2: Run checks to verify failure**

Run: `go test ./... -run TestDirectoryPagesUseTheirOwnActiveNavigation`
Expected: FAIL while both pages mark `博客` active.

- [ ] **Step 3: Replace page shells from the blog baseline**

```html
<a class="blog-profile-nav__item is-active" href="/blog/categories">...</a>
<a class="blog-profile-nav__item" href="/blog/">...</a>
```

- [ ] **Step 4: Run checks to verify pass**

Run: `go test ./... -run TestDirectoryPagesUseTheirOwnActiveNavigation`
Expected: PASS.

### Task 2: Rebuild the category directory renderer

**Files:**
- Modify: `static/blog/categories.html`
- Modify: `static/blog/js/categories.js`
- Modify: `static/blog/css/blog.css`
- Test: `blog_directory_ui_test.go`

**Interfaces:**
- Consumes: `BlogAPI.getCategories()`, `BlogAPI.getTags()`, and `BlogAPI.getPosts({ page: 1, limit: 200 })`.
- Produces: `renderCategoryDirectory(groups, focus)` HTML containing `.category-directory__group` and article links to `/blog/post/<slug>`.

- [ ] **Step 1: Write a failing renderer check**

```go
func TestCategoryPageHasReferenceDirectoryMount(t *testing.T) {
    requireContains(t, "static/blog/categories.html", `id="categoryDirectory"`)
    requireContains(t, "static/blog/categories.html", `CATEGORY / DIRECTORY`)
}
```

- [ ] **Step 2: Run it and confirm failure**

Run: `go test ./... -run TestCategoryPageHasReferenceDirectoryMount`
Expected: FAIL before the old mount is removed.

- [ ] **Step 3: Implement the reference directory**

```js
function renderCategoryDirectory(groups, focus) {
  return names.map(renderCategoryGroup).join('')
    + '<p class="directory-end">-- 已经到底了 --</p>';
}
```

- [ ] **Step 4: Scope layout CSS to `.blog-categories-page`**

```css
.category-directory__posts { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); }
.category-directory__group::before { content: ''; border-left: 1px solid var(--blog-line); }
```

- [ ] **Step 5: Run checks and syntax validation**

Run: `go test ./... -run TestCategoryPageHasReferenceDirectoryMount; node --check static/blog/js/categories.js`
Expected: PASS with no JavaScript syntax errors.

### Task 3: Rebuild the tag constellation renderer

**Files:**
- Modify: `static/blog/tags.html`
- Modify: `static/blog/js/tags.js`
- Modify: `static/blog/css/blog.css`
- Test: `blog_directory_ui_test.go`

**Interfaces:**
- Consumes: tags with `{ name, count }` and posts with `tags` arrays.
- Produces: `renderTagConstellation(tags, posts, focus)` and a focused article list mounted into `#tagArticles`.

- [ ] **Step 1: Write a failing page-structure check**

```go
func TestTagPageHasConstellationAndArticleMounts(t *testing.T) {
    requireContains(t, "static/blog/tags.html", `id="tagConstellation"`)
    requireContains(t, "static/blog/tags.html", `id="tagArticles"`)
}
```

- [ ] **Step 2: Run it and confirm failure**

Run: `go test ./... -run TestTagPageHasConstellationAndArticleMounts`
Expected: FAIL before the old tag pile markup is replaced.

- [ ] **Step 3: Implement deterministic constellation placement and focus filtering**

```js
function articleMatchesTag(post, name) {
  return (post.tags || []).indexOf(name) >= 0;
}

function selectTag(name) {
  history.replaceState(null, '', '/blog/tags?focus=' + encodeURIComponent(name));
  renderFocusedArticles(name);
}
```

- [ ] **Step 4: Add constellation and focused-list styles**

```css
.tag-constellation__node { --node-size: 1; transform: translate(var(--x), var(--y)); }
.tag-constellation__node.is-focus { box-shadow: 0 0 36px rgba(125, 184, 246, .82); }
```

- [ ] **Step 5: Run checks and JavaScript validation**

Run: `go test ./... -run TestTagPageHasConstellationAndArticleMounts; node --check static/blog/js/tags.js`
Expected: PASS with no JavaScript syntax errors.

### Task 4: Verify local deployment

**Files:**
- Modify: `static/blog/categories.html`
- Modify: `static/blog/tags.html`
- Modify: `static/blog/css/blog.css`

- [ ] **Step 1: Build and restart the local frontend**

Run: `docker compose -f deploy/docker/compose.yaml up -d --build frontend`
Expected: `frontend` and `gateway` report `Up`.

- [ ] **Step 2: Verify served pages and source contracts**

Run: `curl.exe -fsS http://127.0.0.1:8080/blog/categories; curl.exe -fsS http://127.0.0.1:8080/blog/tags`
Expected: category response includes `categoryDirectory`; tag response includes `tagConstellation` and `tagArticles`.

- [ ] **Step 3: Run final regression checks**

Run: `go test ./... -run 'TestDirectoryPages|TestCategoryPage|TestTagPage'; git diff --check`
Expected: PASS and no whitespace errors.
