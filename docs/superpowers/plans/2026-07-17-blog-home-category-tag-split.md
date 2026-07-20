# 博客主页 / 分类 / 标签拆分实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将当前博客前端拆分为独立的博客主页、分类页和标签页，并统一为一体式侧栏导航体验。

**Architecture:** 继续使用当前静态页面 + Nginx 分发方案，不新增后端接口。`/blog/` 保留文章流主页，新增 `/blog/categories` 与 `/blog/tags` 两个静态页面，各自通过独立脚本复用现有文章、分类、标签接口，在前端完成分组、标签墙和导航高亮逻辑。

**Tech Stack:** 静态 HTML、Vanilla JavaScript、CSS、Go 轻量回归测试、Docker + Nginx 静态分发

## Global Constraints

- 左侧导航固定为 `首页 / 博客 / 分类 / 标签 / 关于`。
- 取消独立悬浮“返回首页”按钮，把首页入口并入侧栏导航。
- `/blog/` 主页不再展示分类区和标签区，只保留 hero 与文章流。
- `/blog/categories` 做成类似归档页的“按分类分组文章列表”。
- `/blog/tags` 做成带“掉落堆砌”进入动效的标签墙，并提供减少动画时的降级。
- 不新增后端接口，不在本轮扩展详情页、分类详情页或标签详情页二级路由。
- 修改完成后需要验证本地测试通过，并重建 Docker `frontend` 服务确认页面实际生效。

---

### Task 1: 锁定三页拆分后的静态入口骨架

**Files:**
- Modify: `blog_index_ui_test.go`
- Create: `static/blog/categories.html`
- Create: `static/blog/tags.html`
- Modify: `static/blog/index.html`

**Interfaces:**
- Consumes: 当前 `static/blog/index.html`、现有博客侧栏和文章流结构
- Produces: 三个静态页面都包含统一侧栏导航、页面主区域容器与各自页面标识

- [ ] **Step 1: 写失败测试，约束三页入口与统一导航骨架**

```go
func TestBlogPagesExposeSplitNavigationShell(t *testing.T) {
	pages := map[string][]string{
		"static/blog/index.html": {
			`<body class="blog-page">`,
			`class="blog-sidebar-shell"`,
			`href="/blog/categories"`,
			`href="/blog/tags"`,
			`class="blog-main blog-home-page"`,
		},
		"static/blog/categories.html": {
			`class="blog-sidebar-shell"`,
			`class="blog-main blog-categories-page"`,
			`id="categoryGroups"`,
		},
		"static/blog/tags.html": {
			`class="blog-sidebar-shell"`,
			`class="blog-main blog-tags-page"`,
			`id="tagPile"`,
		},
	}

	for path, markers := range pages {
		content := readFile(t, path)
		for _, marker := range markers {
			if !strings.Contains(content, marker) {
				t.Fatalf("expected %s to contain %q", path, marker)
			}
		}
	}
}
```

- [ ] **Step 2: 运行测试，确认它先失败**

Run: `go test ./...`
Expected: FAIL，提示缺少 `categories.html`、`tags.html` 或新导航标记。

- [ ] **Step 3: 更新 `static/blog/index.html` 为博客主页骨架**

```html
<main class="blog-container">
  <div class="blog-shell">
    <aside class="blog-sidebar-shell">
      <nav class="blog-profile-nav" aria-label="博客导航">
        <a class="blog-profile-nav__item" href="/">首页</a>
        <a class="blog-profile-nav__item is-active" href="/blog/">博客</a>
        <a class="blog-profile-nav__item" href="/blog/categories">分类</a>
        <a class="blog-profile-nav__item" href="/blog/tags">标签</a>
        <a class="blog-profile-nav__item" href="/">关于</a>
      </nav>
    </aside>
    <section class="blog-main blog-home-page">
      <header class="blog-hero">...</header>
      <div id="postList" class="post-list post-list--stream"></div>
    </section>
  </div>
</main>
```

- [ ] **Step 4: 新建 `static/blog/categories.html` 页面骨架**

```html
<section class="blog-main blog-categories-page">
  <header class="blog-section-hero">
    <h2 class="blog-section-title">分类</h2>
    <p class="blog-section-subtitle" id="categoriesSummary">加载中...</p>
  </header>
  <section class="blog-directory-panel">
    <div id="categoryGroups" class="blog-category-groups">
      <div class="blog-loading">加载中...</div>
    </div>
  </section>
</section>
```

- [ ] **Step 5: 新建 `static/blog/tags.html` 页面骨架**

```html
<section class="blog-main blog-tags-page">
  <header class="blog-section-hero">
    <h2 class="blog-section-title">标签</h2>
    <p class="blog-section-subtitle" id="tagsSummary">加载中...</p>
  </header>
  <section class="blog-tags-panel">
    <div id="tagPile" class="blog-tag-pile">
      <div class="blog-loading">加载中...</div>
    </div>
  </section>
</section>
```

- [ ] **Step 6: 运行测试，确认入口骨架存在**

Run: `go test ./...`
Expected: 仍可能 FAIL，但不再是缺页面或缺入口骨架。

### Task 2: 抽出统一侧栏与博客主页逻辑

**Files:**
- Modify: `static/blog/index.html`
- Modify: `static/blog/js/list.js`
- Modify: `static/blog/css/blog.css`
- Modify: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: `BlogAPI.getPosts`、`BlogAPI.getCategories`、`BlogAPI.getTags`
- Produces: 统一侧栏统计、博客主页文章流、分类/标签点击跳转到独立页

- [ ] **Step 1: 写失败测试，约束主页不再渲染分类/标签模块**

```go
func TestBlogHomeRemovesInlineCategoryAndTagPanels(t *testing.T) {
	content := readFile(t, "static/blog/index.html")
	forbidden := []string{
		`id="categoryChips"`,
		`id="tagChips"`,
		`class="blog-filter-grid"`,
	}
	for _, marker := range forbidden {
		if strings.Contains(content, marker) {
			t.Fatalf("expected blog home to remove %q", marker)
		}
	}
}
```

- [ ] **Step 2: 运行测试，确认主页还未移除旧模块**

Run: `go test ./...`
Expected: FAIL，提示主页仍含旧分类/标签区域。

- [ ] **Step 3: 精简 `static/blog/index.html`，只保留 hero 与文章流**

```html
<section class="blog-main blog-home-page">
  <header class="blog-hero">...</header>
  <div id="postList" class="post-list post-list--stream">
    <div class="blog-loading">加载中...</div>
  </div>
  <div id="pagination" class="blog-pagination" style="display:none;"></div>
</section>
```

- [ ] **Step 4: 更新 `static/blog/js/list.js` 的点击跳转逻辑**

```javascript
function buildCategoryHref(categoryName) {
    return '/blog/categories?focus=' + encodeURIComponent(categoryName);
}

function buildTagHref(tagName) {
    return '/blog/tags?focus=' + encodeURIComponent(tagName);
}
```

- [ ] **Step 5: 在文章卡片渲染中改用独立页链接**

```javascript
var categoryHtml = '<a class="post-card__category" href="' + buildCategoryHref(p.category || '') + '">' + escapeHtml(p.category || '未分类') + '</a>';
var tagsHtml = (p.tags || []).map(function (t) {
    return '<a class="post-card__tag" href="' + buildTagHref(t) + '">' + escapeHtml(t) + '</a>';
}).join('');
```

- [ ] **Step 6: 在 `static/blog/css/blog.css` 中把侧栏选择器统一为一体式外壳**

```css
.blog-sidebar-shell {
    position: sticky;
    top: 0;
    min-height: 100vh;
    padding: 28px 20px;
    background: rgba(7, 15, 24, 0.82);
    border-right: 1px solid rgba(126, 178, 232, 0.18);
}
```

- [ ] **Step 7: 运行测试，确认主页职责收口**

Run: `go test ./...`
Expected: PASS 当前新增主页结构断言。

### Task 3: 实现分类页的归档式文章目录

**Files:**
- Create: `static/blog/js/categories.js`
- Modify: `static/blog/css/blog.css`
- Modify: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: `BlogAPI.getCategories()`、`BlogAPI.getPosts({ limit })`
- Produces: 分类页分组目录数据结构与渲染结果

- [ ] **Step 1: 写失败测试，约束分类页脚本与目录结构**

```go
func TestCategoriesPageRendersDirectoryGroups(t *testing.T) {
	script := readFile(t, "static/blog/js/categories.js")
	required := []string{
		`renderCategoryGroups`,
		`category-group`,
		`category-group__posts`,
	}
	for _, marker := range required {
		if !strings.Contains(script, marker) {
			t.Fatalf("expected categories.js to contain %q", marker)
		}
	}
}
```

- [ ] **Step 2: 运行测试，确认 `categories.js` 还不存在**

Run: `go test ./...`
Expected: FAIL，提示缺少分类页脚本。

- [ ] **Step 3: 新建 `static/blog/js/categories.js`，拉取文章并按分类聚合**

```javascript
function groupPostsByCategory(posts) {
    var groups = {};
    posts.forEach(function (post) {
        var key = post.category || '未分类';
        if (!groups[key]) groups[key] = [];
        groups[key].push(post);
    });
    return groups;
}
```

- [ ] **Step 4: 渲染归档式分类组**

```javascript
function renderCategoryGroups(groups) {
    return Object.keys(groups).sort().map(function (name) {
        return '<section class="category-group" id="category-' + slugify(name) + '">'
            + '<header class="category-group__head">'
            + '<h3 class="category-group__title">' + escapeHtml(name) + '</h3>'
            + '<span class="category-group__count">' + groups[name].length + '</span>'
            + '</header>'
            + '<div class="category-group__posts">' + groups[name].map(renderCategoryPostLink).join('') + '</div>'
            + '</section>';
    }).join('');
}
```

- [ ] **Step 5: 在 `static/blog/css/blog.css` 中补分类目录样式**

```css
.category-group {
    padding: 22px 0 28px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.category-group__posts {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px 36px;
}
```

- [ ] **Step 6: 运行测试，确认分类页结构通过**

Run: `go test ./...`
Expected: PASS 当前分类页结构断言。

### Task 4: 实现标签页堆叠墙与掉落动画

**Files:**
- Create: `static/blog/js/tags.js`
- Modify: `static/blog/css/blog.css`
- Modify: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: `BlogAPI.getTags()`
- Produces: 标签墙布局、掉落动效类名、减少动画时的降级样式

- [ ] **Step 1: 写失败测试，约束标签页脚本与动效类名**

```go
func TestTagsPageSupportsPileAnimation(t *testing.T) {
	script := readFile(t, "static/blog/js/tags.js")
	styles := readFile(t, "static/blog/css/blog.css")
	requiredScript := []string{
		`renderTagPile`,
		`tag-pile__item`,
		`prefers-reduced-motion`,
	}
	for _, marker := range requiredScript {
		if !strings.Contains(script, marker) && !strings.Contains(styles, marker) {
			t.Fatalf("expected tag page assets to contain %q", marker)
		}
	}
}
```

- [ ] **Step 2: 运行测试，确认标签页脚本还未实现**

Run: `go test ./...`
Expected: FAIL，提示缺少 `tags.js` 或动画样式。

- [ ] **Step 3: 新建 `static/blog/js/tags.js`，为每个标签生成可控位置与旋转**

```javascript
function buildTagLayout(tags) {
    return tags.map(function (tag, index) {
        return {
            name: tag.name,
            count: tag.count,
            rotate: ((index % 7) - 3) * 4,
            delay: index * 35,
            column: index % 5,
            row: Math.floor(index / 5)
        };
    });
}
```

- [ ] **Step 4: 渲染标签堆叠墙**

```javascript
function renderTagPile(items) {
    return items.map(function (item) {
        return '<a class="tag-pile__item" href="/blog/?tag=' + encodeURIComponent(item.name) + '"'
            + ' style="--tag-rotate:' + item.rotate + 'deg;--tag-delay:' + item.delay + 'ms;--tag-column:' + item.column + ';--tag-row:' + item.row + ';">'
            + escapeHtml(item.name) + '<span>' + item.count + '</span></a>';
    }).join('');
}
```

- [ ] **Step 5: 在 `static/blog/css/blog.css` 中补掉落动画与降级**

```css
.tag-pile__item {
    animation: tag-fall 680ms cubic-bezier(0.24, 0.9, 0.28, 1.08) both;
    animation-delay: var(--tag-delay);
    transform: translateY(-120px) rotate(var(--tag-rotate));
}

@media (prefers-reduced-motion: reduce) {
    .tag-pile__item {
        animation: none;
        transform: rotate(var(--tag-rotate));
    }
}
```

- [ ] **Step 6: 运行测试，确认标签页结构与动画类名通过**

Run: `go test ./...`
Expected: PASS 当前标签页结构断言。

### Task 5: 配置静态路由并验证 Docker 实际页面

**Files:**
- Modify: `frontend/nginx/default.conf`
- Modify: `blog_index_ui_test.go`

**Interfaces:**
- Consumes: Nginx 静态分发规则、Docker `frontend` 构建流程
- Produces: `/blog/categories`、`/blog/tags` 可直接访问，Docker 重建后页面生效

- [ ] **Step 1: 写失败测试，约束 Nginx 支持新的静态入口**

```go
func TestNginxRoutesSupportBlogSplitPages(t *testing.T) {
	conf := readFile(t, "frontend/nginx/default.conf")
	required := []string{
		`location = /blog/categories`,
		`location = /blog/tags`,
		`try_files /blog/categories.html =404;`,
		`try_files /blog/tags.html =404;`,
	}
	for _, marker := range required {
		if !strings.Contains(conf, marker) {
			t.Fatalf("expected nginx config to contain %q", marker)
		}
	}
}
```

- [ ] **Step 2: 运行测试，确认路由规则先失败**

Run: `go test ./...`
Expected: FAIL，提示 Nginx 规则缺失。

- [ ] **Step 3: 在 `frontend/nginx/default.conf` 增加两条静态页面路由**

```nginx
location = /blog/categories {
    try_files /blog/categories.html =404;
}

location = /blog/tags {
    try_files /blog/tags.html =404;
}
```

- [ ] **Step 4: 运行完整测试，确认所有断言通过**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: 重建前端 Docker 服务并验证线上 HTML**

Run: `docker compose -f deploy/docker/compose.yaml up -d --build frontend`
Expected: `docker-frontend-1` recreated successfully

- [ ] **Step 6: 抓取线上页面确认新入口生效**

Run: `powershell -Command "Invoke-WebRequest -UseBasicParsing 'http://127.0.0.1:8080/blog/categories' | Select-Object -ExpandProperty Content"`
Expected: HTML 中包含 `blog-categories-page`

- [ ] **Step 7: 再抓取标签页确认新入口生效**

Run: `powershell -Command "Invoke-WebRequest -UseBasicParsing 'http://127.0.0.1:8080/blog/tags' | Select-Object -ExpandProperty Content"`
Expected: HTML 中包含 `blog-tags-page`
