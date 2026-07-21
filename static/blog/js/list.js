(function () {
    'use strict';

    var state = {
        page: 1,
        limit: 4,
        category: '',
        tag: '',
        search: '',
        total: 0,
        posts: [],
        hasMore: true,
        loading: false
    };

    function readQueryParams() {
        var params = new URLSearchParams(window.location.search);
        if (params.get('category')) state.category = params.get('category');
        if (params.get('tag')) state.tag = params.get('tag');
        if (params.get('keyword')) state.search = params.get('keyword');
    }

    function buildCategoryHref(categoryName) {
        return '/blog/categories?focus=' + encodeURIComponent(categoryName || 'uncategorized');
    }

    function buildTagHref(tagName) {
        return '/blog/tags?focus=' + encodeURIComponent(tagName || '');
    }

    // The collection endpoint can omit tags, so resolve only cards that need them.
    function hydratePostTags(posts) {
        var missingTags = (posts || []).filter(function (post) {
            return post.id && (!post.tags || post.tags.length === 0);
        });
        if (missingTags.length === 0) return Promise.resolve(posts);
        return Promise.all(missingTags.map(function (post) {
            return BlogAPI.getPost(post.id).then(function (detail) {
                if (detail && detail.tags && detail.tags.length > 0) post.tags = detail.tags;
                return post;
            }).catch(function () {
                return post;
            });
        })).then(function () {
            return posts;
        });
    }

    function escapeHtml(value) {
        if (value == null) return '';
        return String(value).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;');
    }

    function postCardHtml(post) {
        var postHref = '/blog/post/' + encodeURIComponent(post.id);
        var tagsHtml = (post.tags || []).map(function (tag) {
            return '<a class="post-card__tag" href="' + buildTagHref(tag) + '">' + escapeHtml(tag) + '</a>';
        }).join('');
        var categoryHtml = '<a class="post-card__category" href="' + buildCategoryHref(post.category || '') + '">' + escapeHtml(post.category || 'Uncategorized') + '</a>';
        var readingTime = post.reading_time ? post.reading_time + ' min read' : 'Updating';

        return '<article class="post-card post-card--text-only">'
            + '<div class="post-card__body">'
            + '<div class="post-card__topline">' + categoryHtml + '</div>'
            + '<h2 class="post-card__title"><a class="post-card__title-link" href="' + postHref + '">' + escapeHtml(post.title) + '</a></h2>'
            + '<div class="post-card__labels">' + tagsHtml + '</div>'
            + '<p class="post-card__summary">' + escapeHtml(post.summary || '') + '</p>'
            + '<div class="post-card__actions"><div class="post-card__meta post-card__meta--footer">'
            + '<span>' + (window.formatDate(post.created_at) || '') + '</span>'
            + '<span>' + readingTime + '</span>'
            + '</div></div></div></article>';
    }

    function appendPosts(posts, replace) {
        var container = document.getElementById('postList');
        if (replace) container.innerHTML = '';
        if (!posts.length && replace) {
            container.innerHTML = '<div class="blog-empty">No articles yet.</div>';
            return;
        }
        container.insertAdjacentHTML('beforeend', posts.map(postCardHtml).join(''));
        animatePostCards(container);
    }

    function animatePostCards(container) {
        Array.prototype.slice.call(container.querySelectorAll('.post-card:not(.is-entered)')).forEach(function (card, index) {
            card.classList.add('is-entering');
            card.style.setProperty('--enter-delay', Math.min(index * 90, 360) + 'ms');
            window.requestAnimationFrame(function () { card.classList.add('is-entered'); });
        });
    }

    function renderLoadingSkeleton() {
        document.getElementById('postList').innerHTML = '<div class="post-card post-card--skeleton"><div class="post-card__body"><div class="skeleton-line skeleton-line--short"></div><div class="skeleton-line"></div><div class="skeleton-line"></div></div></div>';
    }

    function updateLoadStatus(message) {
        var status = document.getElementById('infiniteScrollStatus');
        if (status) status.textContent = message;
    }

    function loadPosts(reset) {
        if (state.loading || (!reset && !state.hasMore)) return;
        if (reset) {
            state.page = 1;
            state.posts = [];
            state.total = 0;
            state.hasMore = true;
            renderLoadingSkeleton();
        }

        state.loading = true;
        updateLoadStatus(state.posts.length ? 'Loading more articles...' : 'Loading articles...');
        return BlogAPI.getPosts({
            page: state.page,
            limit: state.limit,
            category: state.category,
            tag: state.tag,
            search: state.search
        }).then(function (data) {
            return hydratePostTags(data.posts || []).then(function (posts) {
                return { data: data, posts: posts };
            });
        }).then(function (result) {
            var data = result.data;
            var posts = result.posts;
            appendPosts(posts, state.posts.length === 0);
            state.posts = state.posts.concat(posts);
            state.total = data.total || state.posts.length;
            state.page += 1;
            state.hasMore = posts.length > 0 && state.posts.length < state.total;
            updateStats({ posts: state.total });
            updateLoadStatus(state.hasMore ? 'Scroll to load more articles' : '-- End of articles --');
        }).catch(function (error) {
            if (state.posts.length === 0) document.getElementById('postList').innerHTML = '<div class="blog-error">Load failed: ' + escapeHtml(error.message) + '</div>';
            updateLoadStatus('Unable to load more articles');
        }).then(function () {
            state.loading = false;
        });
    }

    function bindInfiniteScroll() {
        var sentinel = document.getElementById('infiniteScrollSentinel');
        if (!sentinel) return;
        if ('IntersectionObserver' in window) {
            new IntersectionObserver(function (entries) {
                if (entries.some(function (entry) { return entry.isIntersecting; })) loadPosts(false);
            }, { rootMargin: '420px 0px' }).observe(sentinel);
            return;
        }
        window.addEventListener('scroll', function () {
            if (window.innerHeight + window.scrollY >= document.documentElement.scrollHeight - 520) loadPosts(false);
        });
    }

    function bindSearch() {
        var input = document.getElementById('searchInput');
        if (!input) return;
        input.value = state.search;
        var timer;
        input.addEventListener('input', function () {
            window.clearTimeout(timer);
            timer = window.setTimeout(function () {
                state.search = input.value.trim();
                updateUrlParams();
                loadPosts(true);
            }, 300);
        });
    }

    function updateUrlParams() {
        var params = new URLSearchParams();
        if (state.category) params.set('category', state.category);
        if (state.tag) params.set('tag', state.tag);
        if (state.search) params.set('keyword', state.search);
        window.history.replaceState({}, '', '/blog/' + (params.toString() ? '?' + params.toString() : ''));
    }

    function bindBackToTop() {
        var button = document.getElementById('backToTop');
        if (!button) return;
        window.addEventListener('scroll', function () { button.classList.toggle('visible', window.scrollY > 400); });
        button.addEventListener('click', function () { window.scrollTo({ top: 0, behavior: 'smooth' }); });
    }

    function updateStats(partial) {
        setText('statPosts', partial.posts);
        setText('statCategories', partial.categories);
        setText('statTags', partial.tags);
    }

    function setText(id, value) {
        if (value == null) return;
        var element = document.getElementById(id);
        if (element) element.textContent = String(value);
    }

    document.addEventListener('DOMContentLoaded', function () {
        readQueryParams();
        bindSearch();
        bindBackToTop();
        bindInfiniteScroll();
        Promise.all([
            BlogAPI.getCategories().then(function (items) { updateStats({ categories: items.length }); return items; }).catch(function () { return []; }),
            BlogAPI.getTags().then(function (items) { updateStats({ tags: items.length }); return items; }).catch(function () { return []; })
        ]).then(function () { loadPosts(true); });
    });
}());