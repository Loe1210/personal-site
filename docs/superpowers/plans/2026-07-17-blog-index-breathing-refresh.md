# Blog Index Breathing Refresh Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refresh the blog index so the sidebar feels lighter and closer to the reference, the WeChat action opens a centered modal, the provided animated pet stands at the bottom of the sidebar, and the article stream gains more breathing room.

**Architecture:** Keep the current split-page blog structure and dark mountain visual language, but refactor the sidebar into three clear UI zones: social actions, lighter nav, and a dedicated pet dock. Implement the WeChat modal and pet runtime in the sidebar script layer, while the article-size refinements stay in the homepage list + stylesheet layer.

**Tech Stack:** Static HTML, vanilla JavaScript, CSS, Go test suite, Docker Compose frontend deployment.

## Global Constraints

- Preserve the existing dark mountain atmosphere and current split-page information architecture.
- Keep the admin entry inside the sidebar navigation list.
- Use the provided `/pet/pet.json` and `/pet/spritesheet.webp` assets instead of introducing a heavy external pet framework.
- The WeChat modal must ship now with polished placeholder copy and a placeholder visual area, even before final real QR content arrives.
- Homepage cards must become smaller and gain more breathing room without losing readability.
- Rebuild the Docker frontend and verify the updated asset versions after implementation.

---

## File Map

- Modify: `static/blog/index.html`
  Purpose: Homepage shell, sidebar action markup, modal container, pet mount point, asset version bumps.
- Modify: `static/blog/categories.html`
  Purpose: Keep sidebar structure consistent with the homepage for social actions, modal, nav, and pet mount point.
- Modify: `static/blog/tags.html`
  Purpose: Keep sidebar structure consistent with the homepage for social actions, modal, nav, and pet mount point.
- Modify: `static/blog/css/blog.css`
  Purpose: Button restyle, modal styling, nav downsizing, pet dock layout, and homepage card breathing refinements.
- Modify: `static/blog/js/sidebar.js`
  Purpose: WeChat modal open/close behavior, placeholder image handling, and pet animation bootstrapping.
- Modify: `static/blog/js/list.js`
  Purpose: Further narrow homepage card layout assumptions only if needed for CTA spacing or runtime hooks.
- Create: `static/blog/js/pet.js`
  Purpose: Small isolated runtime that reads `/pet/pet.json`, loads the spritesheet, and animates the sidebar pet.
- Test: `blog_index_ui_test.go`
  Purpose: Assert new homepage/sidebar markers, modal mount, pet mount, and updated script/style references.

### Task 1: Sidebar Social Refresh And WeChat Modal

**Files:**
- Modify: `static/blog/index.html`
- Modify: `static/blog/categories.html`
- Modify: `static/blog/tags.html`
- Modify: `static/blog/css/blog.css`
- Modify: `static/blog/js/sidebar.js`
- Test: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: Existing `.blog-profile-actions`, `.blog-profile-nav`, and current `sidebar.js` WeChat-toggle behavior.
- Produces:
  - `#blogWechatModal` modal root present on all three blog pages.
  - `[data-wechat-open]`, `[data-wechat-close]`, `[data-wechat-overlay]` selectors handled by `sidebar.js`.
  - `initSidebarInteractions()` in `static/blog/js/sidebar.js` bootstraps modal behavior on DOM ready.

- [ ] **Step 1: Write the failing homepage UI test assertions**

```go
func TestBlogPagesExposeWechatModalAndPetMount(t *testing.T) {
    html := readStaticFile(t, "static/blog/index.html")
    if !strings.Contains(html, `id="blogWechatModal"`) {
        t.Fatalf("expected homepage to include WeChat modal root")
    }
    if !strings.Contains(html, `data-wechat-open`) {
        t.Fatalf("expected homepage to include WeChat open trigger")
    }
    if !strings.Contains(html, `id="blogSidebarPet"`) {
        t.Fatalf("expected homepage to include sidebar pet mount")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./...`
Expected: FAIL in `blog_index_ui_test.go` because the modal root and pet mount do not exist yet.

- [ ] **Step 3: Update all three blog HTML files with modal and pet shell markup**

```html
<div class="blog-profile-actions">
    <button class="blog-profile-action blog-profile-action--wechat" type="button" data-wechat-open title="微信">
        <span class="blog-profile-action__glyph">微</span>
    </button>
    <a class="blog-profile-action" href="mailto:2509717573@qq.com" title="邮箱">
        <i class="social iconfont icon-email"></i>
    </a>
    <a class="blog-profile-action" href="https://github.com/Loe1210" target="_blank" rel="noreferrer" title="GitHub">
        <i class="social iconfont icon-github"></i>
    </a>
</div>

<div class="blog-sidebar-pet-dock">
    <div id="blogSidebarPet" class="blog-sidebar-pet" aria-hidden="true"></div>
</div>

<div id="blogWechatModal" class="blog-wechat-modal" hidden>
    <button class="blog-wechat-modal__overlay" type="button" data-wechat-overlay aria-label="关闭微信弹窗"></button>
    <section class="blog-wechat-modal__card" role="dialog" aria-modal="true" aria-labelledby="blogWechatTitle">
        <button class="blog-wechat-modal__close" type="button" data-wechat-close aria-label="关闭">×</button>
        <h3 id="blogWechatTitle" class="blog-wechat-modal__title">欢迎通过微信交流</h3>
        <p class="blog-wechat-modal__desc">这里会放你的最终介绍文字、二维码或联系引导。当前先使用完整弹窗骨架，后续可以直接替换内容。</p>
        <div class="blog-wechat-modal__image" role="img" aria-label="微信二维码占位图"></div>
        <p class="blog-wechat-modal__hint">后续可替换为真实二维码、公众号图或联系说明图。</p>
    </section>
</div>
```

- [ ] **Step 4: Implement the modal interaction in `static/blog/js/sidebar.js`**

```js
(function () {
    function initSidebarInteractions() {
        var modal = document.getElementById('blogWechatModal');
        var openTrigger = document.querySelector('[data-wechat-open]');
        var overlay = document.querySelector('[data-wechat-overlay]');
        var closeButton = document.querySelector('[data-wechat-close]');

        if (!modal || !openTrigger) return;

        function setModalOpen(isOpen) {
            modal.hidden = !isOpen;
            document.body.classList.toggle('blog-wechat-modal-open', isOpen);
        }

        openTrigger.addEventListener('click', function () {
            setModalOpen(true);
        });

        [overlay, closeButton].forEach(function (node) {
            if (!node) return;
            node.addEventListener('click', function () {
                setModalOpen(false);
            });
        });

        document.addEventListener('keydown', function (event) {
            if (event.key === 'Escape') setModalOpen(false);
        });
    }

    document.addEventListener('DOMContentLoaded', initSidebarInteractions);
})();
```

- [ ] **Step 5: Style the buttons and modal in `static/blog/css/blog.css`**

```css
.blog-profile-actions {
    grid-template-columns: repeat(3, 44px);
    justify-content: center;
    gap: 14px 12px;
}

.blog-profile-action {
    width: 44px;
    height: 44px;
    border-radius: 999px;
    background: #fff;
    color: #5b7090;
}

.blog-profile-action--wechat {
    background: linear-gradient(135deg, #43d17a, #2fb768);
    color: #fff;
}

.blog-wechat-modal {
    position: fixed;
    inset: 0;
    z-index: 40;
    display: grid;
    place-items: center;
}

.blog-wechat-modal__card {
    width: min(100%, 360px);
    padding: 28px 28px 24px;
    border-radius: 24px;
    background: rgba(255,255,255,0.96);
    color: #1f2630;
}
```

- [ ] **Step 6: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS for the new sidebar/modal HTML markers and no regression in existing blog tests.

- [ ] **Step 7: Commit**

```bash
git add static/blog/index.html static/blog/categories.html static/blog/tags.html static/blog/css/blog.css static/blog/js/sidebar.js blog_index_ui_test.go
git commit -m "feat: add wechat modal and lighter sidebar actions"
```

### Task 2: Sidebar Pet Integration

**Files:**
- Create: `static/blog/js/pet.js`
- Modify: `static/blog/index.html`
- Modify: `static/blog/categories.html`
- Modify: `static/blog/tags.html`
- Modify: `static/blog/css/blog.css`
- Test: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: `pet/pet.json` with `displayName`, `description`, and `spritesheetPath` fields.
- Produces:
  - `window.BlogSidebarPet.init()` in `static/blog/js/pet.js`.
  - `#blogSidebarPet` mount rendered at the bottom of the sidebar.
  - Pages load `/blog/js/pet.js?v=1` after `sidebar.js`.

- [ ] **Step 1: Write the failing test for pet runtime reference**

```go
func TestBlogHomepageLoadsPetRuntime(t *testing.T) {
    html := readStaticFile(t, "static/blog/index.html")
    if !strings.Contains(html, `/blog/js/pet.js?v=1`) {
        t.Fatalf("expected homepage to load pet runtime")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./...`
Expected: FAIL because the pet runtime script is not referenced yet.

- [ ] **Step 3: Create `static/blog/js/pet.js` with a minimal sprite runtime**

```js
(function () {
    function initPet() {
        var mount = document.getElementById('blogSidebarPet');
        if (!mount) return;

        fetch('/pet/pet.json')
            .then(function (response) { return response.json(); })
            .then(function (meta) {
                var pet = document.createElement('div');
                pet.className = 'blog-sidebar-pet__sprite';
                pet.style.backgroundImage = 'url(/pet/' + meta.spritesheetPath + ')';
                mount.appendChild(pet);
            })
            .catch(function () {
                mount.setAttribute('data-pet-fallback', 'true');
            });
    }

    window.BlogSidebarPet = { init: initPet };
    document.addEventListener('DOMContentLoaded', initPet);
})();
```

- [ ] **Step 4: Add the script tag and sidebar pet styling**

```html
<script src="/blog/js/sidebar.js?v=2"></script>
<script src="/blog/js/pet.js?v=1"></script>
```

```css
.blog-sidebar-shell {
    padding-bottom: 0;
}

.blog-sidebar-pet-dock {
    margin-top: auto;
    min-height: 180px;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    overflow: hidden;
}

.blog-sidebar-pet__sprite {
    width: 148px;
    height: 148px;
    background-repeat: no-repeat;
    background-size: cover;
    animation: blog-pet-float 2.8s ease-in-out infinite;
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./...`
Expected: PASS and the homepage references `/blog/js/pet.js?v=1`.

- [ ] **Step 6: Commit**

```bash
git add static/blog/js/pet.js static/blog/index.html static/blog/categories.html static/blog/tags.html static/blog/css/blog.css blog_index_ui_test.go
git commit -m "feat: mount animated sidebar pet"
```

### Task 3: Homepage Card Breathing Room Refinement And Final Verification

**Files:**
- Modify: `static/blog/css/blog.css`
- Modify: `static/blog/js/list.js`
- Test: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: Existing homepage article stream markup using `.post-card__media-link`, `.post-card__title-link`, and `.post-card__readmore`.
- Produces:
  - Narrower homepage content width.
  - Smaller cover area, lighter title spacing, and more vertical breathing room.
  - Final asset versions aligned across homepage and split pages.

- [ ] **Step 1: Write the failing test for updated asset versions**

```go
func TestBlogHomepageUsesUpdatedRefreshAssets(t *testing.T) {
    html := readStaticFile(t, "static/blog/index.html")
    if !strings.Contains(html, `/blog/css/blog.css?v=8`) {
        t.Fatalf("expected refreshed stylesheet version")
    }
    if !strings.Contains(html, `/blog/js/sidebar.js?v=2`) {
        t.Fatalf("expected refreshed sidebar script version")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./...`
Expected: FAIL because the current asset versions are still older values.

- [ ] **Step 3: Tighten the homepage card sizing in `static/blog/css/blog.css`**

```css
.blog-page--home .post-list--stream,
.blog-page--home .blog-pagination {
    width: min(100%, 760px);
}

.blog-page--home .post-list--stream {
    gap: 34px;
}

.blog-page--home .post-card__media,
.blog-page--home .post-card--skeleton .post-card__media {
    min-height: 154px;
    height: 154px;
}

.blog-page--home .post-card__body {
    padding: 18px 24px 22px;
}

.blog-page--home .post-card__title {
    font-size: 1.58em;
    margin-bottom: 14px;
}

.blog-page--home .post-card__summary {
    max-width: 600px;
    line-height: 1.95;
}
```

- [ ] **Step 4: Bump asset versions and keep runtime hooks stable**

```html
<link rel="stylesheet" href="/blog/css/blog.css?v=8">
<script src="/blog/js/sidebar.js?v=2"></script>
<script src="/blog/js/pet.js?v=1"></script>
<script src="/blog/js/list.js?v=10"></script>
```

- [ ] **Step 5: Run the full verification suite**

Run: `go test ./...`
Expected: PASS

Run: `docker compose -f deploy/docker/compose.yaml up -d --build frontend`
Expected: frontend container rebuilds successfully and restarts without errors.

Run: `Invoke-WebRequest -UseBasicParsing http://127.0.0.1:8080/blog/`
Expected: response HTML contains `blog.css?v=8`, `sidebar.js?v=2`, and `/blog/js/pet.js?v=1`.

- [ ] **Step 6: Commit**

```bash
git add static/blog/css/blog.css static/blog/index.html static/blog/categories.html static/blog/tags.html static/blog/js/list.js static/blog/js/sidebar.js static/blog/js/pet.js blog_index_ui_test.go
git commit -m "feat: refine blog breathing room and deploy pet sidebar"
```

## Self-Review

- Spec coverage: Task 1 covers the social button refresh and WeChat modal. Task 2 covers bottom-mounted animated pet integration from `/pet`. Task 3 covers smaller homepage cards, breathing room, asset versioning, and final verification.
- Placeholder scan: No TODO or TBD placeholders remain in the plan steps. The placeholder content requirement is intentionally explicit and limited to the WeChat modal copy/image area described in the approved spec.
- Type consistency: The plan consistently uses `#blogWechatModal`, `#blogSidebarPet`, `window.BlogSidebarPet.init()`, `/blog/js/pet.js?v=1`, and `[data-wechat-open]` / `[data-wechat-close]` / `[data-wechat-overlay]` across tasks.