(function () {
    'use strict';

    var state = {
        page: 1,
        limit: 4,
        category: '',
        tag: '',
        search: ''
    };

    function readQueryParams() {
        var params = new URLSearchParams(window.location.search);
        if (params.get('category')) state.category = params.get('category');
        if (params.get('tag')) state.tag = params.get('tag');
        if (params.get('page')) state.page = parseInt(params.get('page'), 10) || 1;
        if (params.get('keyword')) state.search = params.get('keyword');
    }

    function buildCategoryHref(categoryName) {
        return '/blog/categories?focus=' + encodeURIComponent(categoryName || '未分类');
    }

    function buildTagHref(tagName) {
        return '/blog/tags?focus=' + encodeURIComponent(tagName || '');
    }

    function renderPosts(data) {
        var container = document.getElementById('postList');
        updateStats({ posts: data.total });

        if (!data.posts || data.posts.length === 0) {
            container.innerHTML = '<div class="blog-empty">暂时没有文章。</div>';
            return;
        }

        container.innerHTML = data.posts.map(function (p) {
            var postHref = '/blog/post/' + encodeURIComponent(p.id);
            var tagsHtml = (p.tags || []).map(function (t) {
                return '<a class="post-card__tag" href="' + buildTagHref(t) + '">' + escapeHtml(t) + '</a>';
            }).join('');

            var categoryHtml = '<a class="post-card__category" href="' + buildCategoryHref(p.category || '') + '">' + escapeHtml(p.category || '未分类') + '</a>';
            var coverHtml = p.cover
                ? '<a class="post-card__media-link" href="' + postHref + '"><div class="post-card__media"><img src="' + escapeHtml(p.cover) + '" alt=""></div></a>'
                : '<a class="post-card__media-link" href="' + postHref + '"><div class="post-card__media post-card__media--placeholder"></div></a>';
            var readingTime = p.reading_time ? p.reading_time + ' 分钟阅读' : '持续更新';

            return ''
                + '<article class="post-card">'
                + '  <div class="post-card__link">'
                +       coverHtml
                + '    <div class="post-card__body">'
                + '      <div class="post-card__topline">' + categoryHtml + '</div>'
                + '      <h2 class="post-card__title"><a class="post-card__title-link" href="' + postHref + '">' + escapeHtml(p.title) + '</a></h2>'
                + '      <div class="post-card__meta">'
                + '        <span>' + (window.formatDate(p.created_at) || '') + '</span>'
                + '        <span>' + readingTime + '</span>'
                + '      </div>'
                + '      <p class="post-card__summary">' + escapeHtml(p.summary || '') + '</p>'
                + '      <div class="post-card__actions">'
                + '        <a class="post-card__readmore" href="' + postHref + '">阅读全文</a>'
                + '        <div class="post-card__tags">' + tagsHtml + '</div>'
                + '      </div>'
                + '    </div>'
                + '  </div>'
                + '</article>';
        }).join('');
        animatePostCards(container);
    }

    function animatePostCards(container) {
        var cards = Array.prototype.slice.call(container.querySelectorAll('.post-card'));
        if (!cards.length) return;

        cards.forEach(function (card, index) {
            card.classList.add('is-entering');
            card.style.setProperty('--enter-delay', Math.min(index * 90, 360) + 'ms');
        });

        window.requestAnimationFrame(function () {
            cards.forEach(function (card) {
                card.classList.add('is-entered');
            });
        });
    }

    function renderPagination(data) {
        var pagination = document.getElementById('pagination');
        var totalPages = Math.ceil(data.total / data.limit);

        if (totalPages <= 1) {
            pagination.style.display = 'none';
            return;
        }

        pagination.style.display = 'flex';
        document.getElementById('prevPage').disabled = (data.page <= 1);
        document.getElementById('nextPage').disabled = (data.page >= totalPages);

        var numbersContainer = document.getElementById('pageNumbers');
        var html = '';
        var maxVisible = 5;
        var start = 1;
        var end = totalPages;

        if (totalPages > maxVisible) {
            var half = Math.floor(maxVisible / 2);
            start = Math.max(1, data.page - half);
            end = Math.min(totalPages, start + maxVisible - 1);
            if (end - start < maxVisible - 1) {
                start = Math.max(1, end - maxVisible + 1);
            }
        }

        if (start > 1) {
            html += '<button class="blog-pagination__num" data-page="1">1</button>';
            if (start > 2) html += '<span class="blog-pagination__ellipsis">...</span>';
        }

        for (var i = start; i <= end; i++) {
            html += '<button class="blog-pagination__num' + (i === data.page ? ' active' : '') + '" data-page="' + i + '">' + i + '</button>';
        }

        if (end < totalPages) {
            if (end < totalPages - 1) html += '<span class="blog-pagination__ellipsis">...</span>';
            html += '<button class="blog-pagination__num" data-page="' + totalPages + '">' + totalPages + '</button>';
        }

        numbersContainer.innerHTML = html;
    }

    function escapeHtml(str) {
        if (str == null) return '';
        return String(str)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    }

    function loadPosts() {
        var skeleton = '';
        for (var i = 0; i < 3; i++) {
            skeleton += '<div class="post-card post-card--skeleton">'
                + '<div class="post-card__media"></div>'
                + '<div class="post-card__body">'
                + '<div class="skeleton-line skeleton-line--short"></div>'
                + '<div class="skeleton-line"></div>'
                + '<div class="skeleton-line"></div>'
                + '</div>'
                + '</div>';
        }
        document.getElementById('postList').innerHTML = skeleton;
        return BlogAPI.getPosts({
            page: state.page,
            limit: state.limit,
            category: state.category,
            tag: state.tag,
            search: state.search
        }).then(function (data) {
            renderPosts(data);
            renderPagination(data);
        }).catch(function (err) {
            document.getElementById('postList').innerHTML = '<div class="blog-error">加载失败：' + escapeHtml(err.message) + '</div>';
        });
    }

    function bindPagination() {
        document.getElementById('prevPage').addEventListener('click', function () {
            if (state.page > 1) {
                state.page--;
                updateUrlParams();
                loadPosts();
                window.scrollTo({ top: 0, behavior: 'smooth' });
            }
        });
        document.getElementById('nextPage').addEventListener('click', function () {
            state.page++;
            updateUrlParams();
            loadPosts();
            window.scrollTo({ top: 0, behavior: 'smooth' });
        });
        document.getElementById('pageNumbers').addEventListener('click', function (e) {
            var btn = e.target.closest('.blog-pagination__num');
            if (!btn) return;
            var page = parseInt(btn.getAttribute('data-page'), 10);
            if (page && page !== state.page) {
                state.page = page;
                updateUrlParams();
                loadPosts();
                window.scrollTo({ top: 0, behavior: 'smooth' });
            }
        });
    }

    function bindSearch() {
        var input = document.getElementById('searchInput');
        if (!input) return;
        input.value = state.search;
        var timer = null;
        input.addEventListener('input', function () {
            clearTimeout(timer);
            timer = setTimeout(function () {
                state.search = input.value.trim();
                state.page = 1;
                updateUrlParams();
                loadPosts();
            }, 300);
        });
    }

    function updateUrlParams() {
        var params = new URLSearchParams();
        if (state.category) params.set('category', state.category);
        if (state.tag) params.set('tag', state.tag);
        if (state.search) params.set('keyword', state.search);
        if (state.page > 1) params.set('page', state.page);
        var url = '/blog/' + (params.toString() ? '?' + params.toString() : '');
        window.history.replaceState({}, '', url);
    }

    function bindBackToTop() {
        var btn = document.getElementById('backToTop');
        if (!btn) return;
        window.addEventListener('scroll', function () {
            btn.classList.toggle('visible', window.scrollY > 400);
        });
        btn.addEventListener('click', function () {
            window.scrollTo({ top: 0, behavior: 'smooth' });
        });
    }

    function updateStats(partial) {
        setText('statPosts', partial.posts);
        setText('statCategories', partial.categories);
        setText('statTags', partial.tags);
    }

    function setText(id, value) {
        if (value == null) return;
        var el = document.getElementById(id);
        if (el) el.textContent = String(value);
    }

    function loadBackground() {
        var script = document.createElement('script');
        script.src = '/assets/json/images.js?t=' + Date.now();
        script.onload = function () {
            if (window.BING_IMAGES && window.BING_IMAGES.length > 0) {
                var key = 'blog-bg-index';
                var index = parseInt(sessionStorage.getItem(key), 10);
                if (isNaN(index) || index >= window.BING_IMAGES.length) {
                    index = Math.floor(Math.random() * window.BING_IMAGES.length);
                    sessionStorage.setItem(key, index);
                }
                var url = '/assets/' + window.BING_IMAGES[index];
                document.body.style.backgroundImage = "url('" + url.replace(/['\\]/g, '\\$&') + "')";
            }
        };
        document.body.appendChild(script);
    }

    document.addEventListener('DOMContentLoaded', function () {
        readQueryParams();
        bindPagination();
        bindSearch();
        bindBackToTop();
        loadBackground();

        var categoriesPromise = BlogAPI.getCategories().then(function (items) {
            updateStats({ categories: items.length });
            return items;
        }).catch(function () { return []; });

        var tagsPromise = BlogAPI.getTags().then(function (items) {
            updateStats({ tags: items.length });
            return items;
        }).catch(function () { return []; });

        Promise.all([categoriesPromise, tagsPromise]).then(function () {
            loadPosts();
        });
    });
})();