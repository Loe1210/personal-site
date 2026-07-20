package main

import (
    "os"
    "strings"
    "testing"
)

func TestBlogPagesExposeSplitNavigationShell(t *testing.T) {
    pages := map[string][]string{
        "static/blog/index.html": {
            `class="blog-page`,
            `class="blog-sidebar-shell"`,
            `href="/blog/categories"`,
            `href="/blog/tags"`,
            `class="blog-main blog-home-page"`,
        },
        "static/blog/categories.html": {
            `class="blog-sidebar-shell"`,
            `class="blog-main blog-categories-page"`,
            `id="categoryDirectory"`,
        },
        "static/blog/tags.html": {
            `class="blog-sidebar-shell"`,
            `class="blog-main blog-tags-page"`,
            `id="tagConstellation"`,
            `id="tagArticles"`,
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



func TestDirectoryPagesActivateTheirOwnSidebarEntry(t *testing.T) {
    cases := map[string]string{
        "static/blog/categories.html": `class="blog-profile-nav__item is-active" href="/blog/categories"`,
        "static/blog/tags.html":       `class="blog-profile-nav__item is-active" href="/blog/tags"`,
    }

    for path, activeEntry := range cases {
        content := readFile(t, path)
        if !strings.Contains(content, activeEntry) {
            t.Fatalf("expected %s to activate %q", path, activeEntry)
        }
        if strings.Contains(content, `class="blog-profile-nav__item is-active" href="/blog/"`) {
            t.Fatalf("expected %s not to activate the blog entry", path)
        }
    }
}

func TestBlogPagesExposeWechatModalAndPetMount(t *testing.T) {
    pages := []string{
        "static/blog/index.html",
        "static/blog/categories.html",
        "static/blog/tags.html",
    }

    required := []string{
        `id="blogWechatModal"`,
        `data-wechat-open`,
        `data-wechat-close`,
        `data-wechat-overlay`,
        `id="blogSidebarPet"`,
    }

    for _, path := range pages {
        content := readFile(t, path)
        for _, marker := range required {
            if !strings.Contains(content, marker) {
                t.Fatalf("expected %s to contain %q", path, marker)
            }
        }
    }
}


func TestBlogHomeUsesSoftEntranceAnimations(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    script := readFile(t, "static/blog/js/list.js")
    html := readFile(t, "static/blog/index.html")

    requiredStyles := []string{
        `Soft blog entrance`,
        `@keyframes blogSoftReveal`,
        `.blog-page--home .post-card.is-entering`,
        `.blog-page--home .post-card.is-entered`,
        `prefers-reduced-motion: reduce`,
    }
    for _, marker := range requiredStyles {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected blog styles to contain entrance marker %q", marker)
        }
    }

    requiredScript := []string{
        `animatePostCards`,
        `is-entering`,
        `is-entered`,
        `requestAnimationFrame`,
    }
    for _, marker := range requiredScript {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected blog list runtime to contain entrance marker %q", marker)
        }
    }

    if !strings.Contains(html, `/blog/css/blog.css?v=19`) || !strings.Contains(html, `/blog/js/list.js?v=11`) {
        t.Fatalf("expected homepage to load entrance animation assets")
    }
}
func TestBlogHomepageLoadsPetRuntime(t *testing.T) {
    html := readFile(t, "static/blog/index.html")
    if !strings.Contains(html, `/blog/js/pet.js?v=11`) {
        t.Fatalf("expected homepage to load pet runtime")
    }
}


func TestBlogHeroUsesSignatureCopyAndCuteCursor(t *testing.T) {
    html := readFile(t, "static/blog/index.html")
    cursor := readFile(t, "static/assets/css/cute-cursor.css")
    styles := readFile(t, "static/blog/css/blog.css")

    requiredHtml := []string{
        `Build quietly. Ship clearly.`,
        `SYSTEMS / NOTES / BACKEND`,
        `/assets/css/cute-cursor.css`,
        `/blog/css/blog.css?v=19`,
    }
    for _, marker := range requiredHtml {
        if !strings.Contains(html, marker) {
            t.Fatalf("expected blog homepage to contain %q", marker)
        }
    }

    requiredCursor := []string{
        `--cute-cursor-default`,
        `--cute-cursor-pointer`,
        `cursor: var(--cute-cursor-default), auto`,
    }
    for _, marker := range requiredCursor {
        if !strings.Contains(cursor, marker) {
            t.Fatalf("expected cute cursor stylesheet to contain %q", marker)
        }
    }

    requiredStyles := []string{
        `Signature hero polish`,
        `Reference layout rebalance`,
        `Compact mist article cards`,
        `.blog-page--home .blog-hero::before {`,
        `.blog-page--home .blog-page-title:hover`,
        `.blog-page--home .post-card:hover .post-card__media img`,
        `@media (prefers-reduced-motion: reduce)`,
    }
    for _, marker := range requiredStyles {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected blog styles to contain %q", marker)
        }
    }
}
func TestBlogHomepageUsesUpdatedRefreshAssets(t *testing.T) {
    html := readFile(t, "static/blog/index.html")
    if !strings.Contains(html, `/blog/css/blog.css?v=19`) {
        t.Fatalf("expected refreshed stylesheet version")
    }
    if !strings.Contains(html, `/blog/js/sidebar.js?v=2`) {
        t.Fatalf("expected refreshed sidebar script version")
    }
    if !strings.Contains(html, `/blog/js/list.js?v=11`) {
        t.Fatalf("expected refreshed list script version")
    }
}

func TestAdminEditorSupportsCoverLibrary(t *testing.T) {
    html := readFile(t, "static/admin/index.html")
    script := readFile(t, "static/admin/js/admin.js")
    styles := readFile(t, "static/admin/css/admin.css")

    required := []string{
        `id="openCoverLibraryBtn"`,
        `id="coverLibraryModal"`,
        `id="coverLibraryResults"`,
        `OPENVERSE_IMAGE_ENDPOINT`,
        `searchCoverLibrary`,
        `renderFallbackCoverLibrary`,
        `cover-library__grid`,
    }
    for _, marker := range required {
        if !strings.Contains(html, marker) && !strings.Contains(script, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected admin cover library assets to contain %q", marker)
        }
    }
}

func TestPetRuntimeAvoidsCssAnimationConflict(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")
    styles := readFile(t, "static/blog/css/blog.css")

    required := []string{
        `requestAnimationFrame`,
        `is-visible`,
        `decode`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected pet runtime to contain %q", marker)
        }
    }
    forbidden := []string{
        `.blog-sidebar-pet.is-playing .blog-sidebar-pet__sprite`,
        `pet-click-pop`,
    }
    for _, marker := range forbidden {
        if strings.Contains(styles, marker) {
            t.Fatalf("expected pet styles to remove animation conflict marker %q", marker)
        }
    }
}
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

func TestCategoriesPageRendersReferenceDirectory(t *testing.T) {
    script := readFile(t, "static/blog/js/categories.js")
    required := []string{
        `renderCategoryDirectory`,
        `category-directory__group`,
        `category-directory__posts`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected categories.js to contain %q", marker)
        }
    }
}

func TestTagsPageSupportsConstellationAndFocus(t *testing.T) {
    script := readFile(t, "static/blog/js/tags.js")
    styles := readFile(t, "static/blog/css/blog.css")
    required := []string{
        `renderTagConstellation`,
        `tag-constellation__node`,
        `renderFocusedArticles`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected tag page assets to contain %q", marker)
        }
    }
}


func TestNginxAllowsAdminMediaUploads(t *testing.T) {
    conf := readFile(t, "frontend/nginx/default.conf")
    required := []string{
        `client_max_body_size 16m;`,
        `location /api/ {`,
    }
    for _, marker := range required {
        if !strings.Contains(conf, marker) {
            t.Fatalf("expected nginx config to contain %q", marker)
        }
    }
}

func TestAdminCoverUploadUsesChunkedMediaApi(t *testing.T) {
    html := readFile(t, "static/admin/index.html")
    script := readFile(t, "static/admin/js/admin.js")
    required := []string{
        `ADMIN_COVER_CHUNK_SIZE`,
        `/api/media/upload/tasks/init`,
        `/api/media/upload/tasks/`,
        `/chunks/`,
        `/complete`,
        `/cancel`,
        `formData.append('sha256', sha256)`,
        `sendChunkWithProgress`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected admin upload script to contain %q", marker)
        }
    }
    if !strings.Contains(script, `crypto.subtle.digest`) {
        t.Fatalf("expected admin upload script to calculate a client-side final file hash")
    }
    if strings.Contains(script, `xhr.open('POST', '/api/media/upload', true)`) {
        t.Fatalf("expected admin upload script to stop using whole-file media upload")
    }
    if !strings.Contains(html, `/admin/js/admin.js?v=13`) {
        t.Fatalf("expected admin page to load the refreshed chunk upload script")
    }
}
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

func readFile(t *testing.T, path string) string {
    t.Helper()

    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("read %s: %v", path, err)
    }

    return string(data)
}



func TestPetRuntimeSupportsPetSwitching(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")
    manifest := readFile(t, "pet/index.json")

    requiredScript := []string{
        `PET_MANIFEST_CANDIDATES`,
        `dblclick`,
        `localStorage`,
        `switchToPet`,
        `ensurePetMount`,
    }
    for _, marker := range requiredScript {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to contain %q", marker)
        }
    }

    requiredManifest := []string{
        `"yuexinmiao"`,
        `"zhangfei-tusun"`,
        `"pet/yuexinmiao/pet.json"`,
        `"pet/zhangfei-tusun/pet.json"`,
    }
    for _, marker := range requiredManifest {
        if !strings.Contains(manifest, marker) {
            t.Fatalf("expected pet manifest to contain %q", marker)
        }
    }
}



func TestPetRuntimeSwitchesPetsOnDoubleClickOnly(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")

    required := []string{
        `dblclick`,
        `switchToPet(nextPetId(runtime.activePetId), mount, true)`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to contain double-click switch marker %q", marker)
        }
    }

    forbidden := []string{
        `blog-pet-switch`,
        `appendChild(switcher)`,
        `createElement('button')`,
    }
    for _, marker := range forbidden {
        if strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to remove switch button marker %q", marker)
        }
    }
}
func TestPetRuntimePreventsDuplicateStaleRenders(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")

    required := []string{
        `window.__hinsBlogPetRuntimeBooted`,
        `renderToken`,
        `token !== runtime.renderToken`,
        `cleanupDuplicateMounts`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to contain stale-render guard %q", marker)
        }
    }
}
func TestAllMainPagesLoadGlobalPetRuntime(t *testing.T) {
    pages := []string{
        "static/index.html",
        "static/admin/index.html",
        "static/blog/index.html",
        "static/blog/categories.html",
        "static/blog/tags.html",
        "static/blog/post.html",
    }

    for _, path := range pages {
        content := readFile(t, path)
        if !strings.Contains(content, `/blog/js/pet.js?v=11`) {
            t.Fatalf("expected %s to load global pet runtime v11", path)
        }
    }
}


func TestPetRuntimeUsesPerActionFrameCounts(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")
    required := []string{
        `frameLimitForState`,
        `idleFrames`,
        `wavingFrames`,
        `jumpingFrames`,
    }
    for _, marker := range required {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to contain frame-count marker %q", marker)
        }
    }

    yuexin := readFile(t, "pet/yuexinmiao/pet.json")
    zhangfei := readFile(t, "pet/zhangfei-tusun/pet.json")
    requiredMeta := map[string][]string{
        "yuexinmiao": {`"idleFrames": 6`, `"wavingFrames": 4`, `"jumpingFrames": 5`},
        "zhangfei": {`"idleFrames": 7`, `"wavingFrames": 4`, `"jumpingFrames": 5`},
    }
    for _, marker := range requiredMeta["yuexinmiao"] {
        if !strings.Contains(yuexin, marker) {
            t.Fatalf("expected yuexin meta to contain %q", marker)
        }
    }
    for _, marker := range requiredMeta["zhangfei"] {
        if !strings.Contains(zhangfei, marker) {
            t.Fatalf("expected zhangfei meta to contain %q", marker)
        }
    }
}
func TestUpdatedZhangfeiPetUsesElevenRows(t *testing.T) {
    meta := readFile(t, "pet/zhangfei-tusun/pet.json")
    required := []string{
        `"columns": 8`,
        `"rows": 11`,
    }
    for _, marker := range required {
        if !strings.Contains(meta, marker) {
            t.Fatalf("expected zhangfei pet meta to contain %q", marker)
        }
    }
}
func TestDockerImagesShipPetAssets(t *testing.T) {
    frontendDockerfile := readFile(t, "frontend/Dockerfile")
    rootDockerfile := readFile(t, "Dockerfile")

    if !strings.Contains(frontendDockerfile, `COPY pet/ /usr/share/nginx/html/pet/`) {
        t.Fatalf("expected frontend Dockerfile to ship pet assets")
    }
    if !strings.Contains(rootDockerfile, `COPY pet /app/static/pet`) {
        t.Fatalf("expected root Dockerfile to ship pet assets")
    }
}
