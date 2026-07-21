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

    if !strings.Contains(html, `/blog/css/blog.css?v=42`) || !strings.Contains(html, `/blog/js/list.js?v=15`) {
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
        `/blog/css/blog.css?v=42`,
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
    if !strings.Contains(html, `/blog/css/blog.css?v=42`) {
        t.Fatalf("expected refreshed stylesheet version")
    }
    if !strings.Contains(html, `/blog/js/sidebar.js?v=6`) {
        t.Fatalf("expected refreshed sidebar script version")
    }
    if !strings.Contains(html, `/blog/js/list.js?v=15`) {
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
        `renderTagArchive`,
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



func TestPetRuntimeUsesSalaryCatOnlyManifest(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")
    manifest := readFile(t, "pet/index.json")

    for _, marker := range []string{`PET_MANIFEST_CANDIDATES`, `switchToPet`, `ensurePetMount`} {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to contain %q", marker)
        }
    }
    for _, marker := range []string{`"defaultPetId": "yuexinmiao"`, `"yuexinmiao"`, `"pet/yuexinmiao/pet.json"`} {
        if !strings.Contains(manifest, marker) {
            t.Fatalf("expected salary-cat-only manifest to contain %q", marker)
        }
    }
    if strings.Contains(manifest, `zhangfei-tusun`) {
        t.Fatal("expected rabbit-cat manifest entry to be removed")
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
    }

    for _, path := range pages {
        content := readFile(t, path)
        if !strings.Contains(content, `/blog/js/pet.js?v=11`) {
            t.Fatalf("expected %s to load global pet runtime v11", path)
        }
    }
}


func TestPetRuntimeUsesSalaryCatActionFrameCounts(t *testing.T) {
    script := readFile(t, "static/blog/js/pet.js")
    for _, marker := range []string{`frameLimitForState`, `idleFrames`, `wavingFrames`, `jumpingFrames`} {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected pet runtime to contain frame-count marker %q", marker)
        }
    }
    yuexin := readFile(t, "pet/yuexinmiao/pet.json")
    for _, marker := range []string{`"idleFrames": 6`, `"wavingFrames": 4`, `"jumpingFrames": 5`} {
        if !strings.Contains(yuexin, marker) {
            t.Fatalf("expected salary-cat meta to contain %q", marker)
        }
    }
    if _, err := os.Stat("pet/zhangfei-tusun"); !os.IsNotExist(err) {
        t.Fatal("expected rabbit-cat resources to be removed")
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

func TestDirectoryPagesUseBreathingEntrances(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    categories := readFile(t, "static/blog/js/categories.js")
    tags := readFile(t, "static/blog/js/tags.js")

    for _, marker := range []string{
        `Directory page breathing entrance`,
        `.blog-page--directory .directory-hero`,
        `.category-directory__group.is-revealed`,
        `.tag-constellation__node.is-revealed`,
        `@keyframes directoryBreathReveal`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected breathing styles to contain %q", marker)
        }
    }
    if !strings.Contains(categories, `animateDirectoryElements`) {
        t.Fatal("expected category directory to animate mounted elements")
    }
    if !strings.Contains(tags, `animateConstellationNodes`) {
        t.Fatal("expected tag constellation to animate mounted elements")
    }
}

func TestPostCardsHaveGranularInteractions(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `Granular article interactions`,
        `.blog-page--home .post-card__title-link::after`,
        `.blog-page--home .post-card__summary:hover`,
        `.blog-page--home .post-card__tags .post-tag:hover`,
        `.blog-page--home .post-card__readmore:focus-visible`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected article interaction styles to contain %q", marker)
        }
    }
}

func TestTagConstellationNodesKeepTheirRevealAndCuteCursor(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `@keyframes constellationNodeReveal`,
        `.blog-page--tag-constellation .tag-constellation__node.is-revealed`,
        `cursor: var(--cute-cursor-pointer), pointer;`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected tag constellation visibility styles to contain %q", marker)
        }
    }
}

func TestDirectoryTextUsesExplicitInteractiveTargets(t *testing.T) {
    categories := readFile(t, "static/blog/js/categories.js")
    styles := readFile(t, "static/blog/css/blog.css")

    for _, marker := range []string{
        `category-directory__heading-link`,
        `/blog/categories?focus=`,
        `.category-directory__heading-link:hover`,
        `.directory-hero__title:hover`,
        `.tag-constellation__name`,
    } {
        if !strings.Contains(categories, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected directory text interaction to contain %q", marker)
        }
    }
}

func TestCategoryLinksSmoothlyFocusTheirDirectoryGroup(t *testing.T) {
    categories := readFile(t, "static/blog/js/categories.js")
    styles := readFile(t, "static/blog/css/blog.css")

    for _, marker := range []string{
        `scrollIntoView({ behavior: 'smooth', block: 'center' })`,
        `history.replaceState`,
        `category-directory__group.is-focus`,
        `category-directory__heading-link`,
    } {
        if !strings.Contains(categories, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected category smooth focus behavior to contain %q", marker)
        }
    }
}
func TestBlogDirectoryPagesReuseTheHomeSidebarLayout(t *testing.T) {
    for _, path := range []string{
        "static/blog/index.html",
        "static/blog/categories.html",
        "static/blog/tags.html",
    } {
        if !strings.Contains(readFile(t, path), "blog-page--sidebar-home") || !strings.Contains(readFile(t, path), "blog-page--home") {
            t.Fatalf("expected %s to reuse the shared home sidebar layout", path)
        }
    }

    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `.blog-page--sidebar-home .blog-sidebar-shell`,
        `.blog-page--directory:not(.blog-page--sidebar-home) .blog-sidebar-shell`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected shared sidebar styles to contain %q", marker)
        }
    }
}
func TestBlogPagesUseTheSharedBackgroundAndMatchedSidebar(t *testing.T) {
    if _, err := os.Stat("static/assets/img/blog-background.png"); err != nil {
        t.Fatalf("expected shared blog background asset: %v", err)
    }

    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `url('/assets/img/blog-background.png')`,
        `Blog background and sidebar palette`,
        `Blog sidebar blue gradient`,
        `.blog-page--sidebar-home .blog-sidebar-shell`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected blog background styles to contain %q", marker)
        }
    }

    for _, path := range []string{
        "static/blog/js/list.js",
        "static/blog/js/categories.js",
        "static/blog/js/tags.js",
        "static/blog/js/post.js",
    } {
        if strings.Contains(readFile(t, path), "BING_IMAGES") {
            t.Fatalf("expected %s not to override the shared blog background", path)
        }
    }
}
func TestBlogPagesShareCenteredBloomingEntrance(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `Unified centered page entrance`,
        `@keyframes blogTitleBloom`,
        `@keyframes blogPageContentReveal`,
        `.directory-hero { text-align: center; }`,
        `.blog-page--home .blog-page-title`,
        `prefers-reduced-motion: reduce`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected unified page entrance styles to contain %q", marker)
        }
    }
}
func TestSidebarCharactersAndSocialActionsHaveDistinctFeedback(t *testing.T) {
    sidebarScript := readFile(t, "static/blog/js/sidebar.js")
    styles := readFile(t, "static/blog/css/blog.css")

    for _, marker := range []string{
        `decorateProfileText`,
        `blog-profile-character`,
        `@keyframes profileCharacterDrop`,
        `#07C160`,
        `#F6C445`,
        `#111827`,
    } {
        if !strings.Contains(sidebarScript, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected sidebar interaction to contain %q", marker)
        }
    }

    for _, path := range []string{
        "static/blog/index.html",
        "static/blog/categories.html",
        "static/blog/tags.html",
    } {
        page := readFile(t, path)
        for _, marker := range []string{
            `blog-profile-action--wechat`,
            `blog-profile-action--email`,
            `blog-profile-action--github`,
        } {
            if !strings.Contains(page, marker) {
                t.Fatalf("expected %s to contain %q", path, marker)
            }
        }
    }
}

func TestTagConstellationRevealsNodesSlowlyOneByOne(t *testing.T) {
    tags := readFile(t, "static/blog/js/tags.js")
    for _, marker := range []string{
        `animateConstellationNodes`,
        `index * 120`,
        `2160`,
    } {
        if !strings.Contains(tags, marker) {
            t.Fatalf("expected tag constellation reveal to contain %q", marker)
        }
    }
}
func TestSidebarSocialIconsShareOneDefaultVisualWeight(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `Sidebar social default normalization`,
        `color: #26374F`,
        `.blog-profile-action__icon-image`,
        `font-size: 18px`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected unified default social icon style to contain %q", marker)
        }
    }
}
func TestPostCardsPlaceLabelsBelowTitlesAndMetadataInFooter(t *testing.T) {
    script := readFile(t, "static/blog/js/list.js")
    styles := readFile(t, "static/blog/css/blog.css")

    for _, marker := range []string{
        `<div class="post-card__topline">' + categoryHtml`,
        `<div class="post-card__labels">' + tagsHtml`,
        `post-card__meta--footer`,
    } {
        if !strings.Contains(script, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected post card layout to contain %q", marker)
        }
    }
    if strings.Contains(script, `阅读全文`) || strings.Contains(script, `post-card__readmore`) {
        t.Fatalf("expected post card renderer to remove the read-more action")
    }
}

func TestWechatIconUsesTheSameDefaultDarkColorAsOtherSocialIcons(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `Unified social icon default color`,
        `color: #111820`,
        `.blog-profile-action--wechat .blog-profile-action__icon-image`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected unified social default color to contain %q", marker)
        }
    }
}
func TestPostListHydratesMissingTagsFromArticleDetails(t *testing.T) {
    script := readFile(t, "static/blog/js/list.js")
    for _, marker := range []string{
        `hydratePostTags`,
        `BlogAPI.getPost(post.id)`,
        `Promise.all`,
    } {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected post list to hydrate missing tags with %q", marker)
        }
    }
}
func TestTagPageProvidesFocusedArchiveAndHydratesMissingTagData(t *testing.T) {
    script := readFile(t, "static/blog/js/tags.js")
    for _, marker := range []string{
        `renderTagArchive`,
        `hydrateTagPosts`,
        `BlogAPI.getPost(post.id)`,
        `data-tag-return`,
    } {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected focused tag archive to contain %q", marker)
        }
    }
}

func TestDirectoryPagesUseSortStyleTitles(t *testing.T) {
    if !strings.Contains(readFile(t, "static/blog/tags.html"), `Tag.sort()`) {
        t.Fatal("expected Tag.sort() heading")
    }
    if !strings.Contains(readFile(t, "static/blog/categories.html"), `Category.sort()`) {
        t.Fatal("expected Category.sort() heading")
    }
}

func TestSidebarNavigationUsesSlidingIndicator(t *testing.T) {
    script := readFile(t, "static/blog/js/sidebar.js")
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `initSidebarNavigationIndicator`,
        `blog-profile-nav__indicator`,
        `pointerenter`,
        `prefers-reduced-motion`,
    } {
        if !strings.Contains(script, marker) && !strings.Contains(styles, marker) {
            t.Fatalf("expected sidebar sliding indicator to contain %q", marker)
        }
    }
}
func TestTagPageLoadsArticleTotalForSidebarStats(t *testing.T) {
    script := readFile(t, "static/blog/js/tags.js")
    for _, marker := range []string{
        `BlogAPI.getPosts({ page: 1, limit: 1 })`,
        `setText('statPosts', postsData.total`,
    } {
        if !strings.Contains(script, marker) {
            t.Fatalf("expected tag page sidebar article total to contain %q", marker)
        }
    }
}
func TestTagArchiveUsesNumericDetailLinksAndHomeUsesInfiniteTextOnlyCards(t *testing.T) {
    tags := readFile(t, "static/blog/js/tags.js")
    list := readFile(t, "static/blog/js/list.js")
    home := readFile(t, "static/blog/index.html")

    for _, marker := range []string{
        `encodeURIComponent(post.id)`,
        `IntersectionObserver`,
        `infiniteScrollSentinel`,
        `appendPosts`,
    } {
        if !strings.Contains(tags, marker) && !strings.Contains(list, marker) && !strings.Contains(home, marker) {
            t.Fatalf("expected tag/detail or infinite-scroll behavior to contain %q", marker)
        }
    }
    if strings.Contains(list, `post-card__media`) || strings.Contains(list, `renderPagination`) || strings.Contains(list, `bindPagination`) {
        t.Fatal("expected home list to remove cover media and button pagination")
    }
}
func TestAllBlogPagesUseBackgroundAndPostDetailUsesBlogPageClass(t *testing.T) {
    post := readFile(t, "static/blog/post.html")
    styles := readFile(t, "static/blog/css/blog.css")
    if !strings.Contains(post, `class="blog-page blog-page--post"`) {
        t.Fatal("expected post detail to share blog page background")
    }
    if !strings.Contains(styles, `blog-background.png`) {
        t.Fatal("expected shared blog background asset")
    }
}

func TestBlogUsesAnimatedStarCanopyAvatarHoverAndDividerFreeDirectoryHero(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    for _, marker := range []string{
        `headerDotBreathe`,
        `.header-star-field`,
        `.blog-profile-avatar:hover`,
        `.blog-page--directory .directory-hero { border-bottom: 0; }`,
    } {
        if !strings.Contains(styles, marker) {
            t.Fatalf("expected blog polish style to contain %q", marker)
        }
    }
}
func TestBlogHeadersRenderIndependentBreathingDotCanopy(t *testing.T) {
    styles := readFile(t, "static/blog/css/blog.css")
    stars := readFile(t, "static/blog/js/header-stars.js")
    for _, marker := range []string{
        `.header-star-field`,
        `headerDotBreathe`,
        `index < 960`,
        `--star-delay`,
    } {
        if !strings.Contains(styles, marker) && !strings.Contains(stars, marker) {
            t.Fatalf("expected independent breathing dot canopy to contain %q", marker)
        }
    }
}
