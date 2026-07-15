/**
 * 博客列表页逻辑
 * 负责：渲染文章卡片、分类/标签筛选、分页、Bing 背景加载
 */
(function () {
    'use strict';

    // ====== 分页 & 筛选状态 ======
    var state = {
        page: 1,
        limit: 4,
        category: '',
        tag: '',
        search: ''
    };

    // 从 URL 读取初始筛选条件（支持从侧边栏锚点跳转）
    function readQueryParams() {
        var params = new URLSearchParams(window.location.search);
        if (params.get('category')) state.category = params.get('category');
        if (params.get('tag')) state.tag = params.get('tag');
        if (params.get('page')) state.page = parseInt(params.get('page'), 10) || 1;
    }

    // ====== 渲染文章列表 ======
    function renderPosts(data) {
        var container = document.getElementById('postList');
        if (!data.posts || data.posts.length === 0) {
            container.innerHTML = '<div class="blog-empty">暂时没有文章。</div>';
            return;
        }

        container.innerHTML = data.posts.map(function (p) {
            var tagsHtml = (p.tags || []).map(function (t) {
                return '<span class="post-card__tag" data-tag="' + escapeHtml(t) + '">' + escapeHtml(t) + '</span>';
            }).join('');

            var categoryHtml = '<span class="post-card__category" data-category="' + escapeHtml(p.category || '') + '">' + escapeHtml(p.category || '未分类') + '</span>';

            var postId = p.id;
            var coverHtml = p.cover ? '<div class="post-card__cover"><img src="' + escapeHtml(p.cover) + '" alt=""></div>' : '';
            return ''
                + '<div class="post-card">'
                + coverHtml
                + '  <a class="post-card__link" href="/blog/post/' + encodeURIComponent(postId) + '">'
                + '    ' + categoryHtml
                + '    <h2 class="post-card__title">' + escapeHtml(p.title) + '</h2>'
                + '    <p class="post-card__summary">' + escapeHtml(p.summary || '') + '</p>'
                + '  </a>'
                + '  <div class="post-card__tags">' + tagsHtml + '</div>'
                + '  <div class="post-card__meta">'
                + '    <span>' + (window.formatDate(p.created_at) || '') + '</span>'
                + '    <span>' + (p.reading_time ? p.reading_time + ' 分钟阅读' : '') + '</span>'
                + '  </div>'
                + '</div>';
        }).join('');
    }

    // ====== 渲染分页 ======
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

        // 渲染页码按钮
        var numbersContainer = document.getElementById('pageNumbers');
        var html = '';
        var maxVisible = 7; // 最多显示7个页码
        var start = 1, end = totalPages;

        if (totalPages > maxVisible) {
            var half = Math.floor(maxVisible / 2);
            start = Math.max(1, data.page - half);
            end = Math.min(totalPages, start + maxVisible - 1);
            if (end - start < maxVisible - 1) {
                start = Math.max(1, end - maxVisible + 1);
            }
        }

        // 第一页 + 省略号
        if (start > 1) {
            html += '<button class="blog-pagination__num" data-page="1">1</button>';
            if (start > 2) html += '<span class="blog-pagination__ellipsis">...</span>';
        }

        for (var i = start; i <= end; i++) {
            html += '<button class="blog-pagination__num' + (i === data.page ? ' active' : '') + '" data-page="' + i + '">' + i + '</button>';
        }

        // 最后一页 + 省略号
        if (end < totalPages) {
            if (end < totalPages - 1) html += '<span class="blog-pagination__ellipsis">...</span>';
            html += '<button class="blog-pagination__num" data-page="' + totalPages + '">' + totalPages + '</button>';
        }

        numbersContainer.innerHTML = html;
    }

    // ====== 渲染分类/标签侧边栏 ======
    function renderCategories(categories) {
        var container = document.getElementById('categoryChips');
        if (!categories || categories.length === 0) {
            container.innerHTML = '<span class="blog-empty">暂无分类</span>';
            return;
        }
        var allActive = !state.category;
        var allHref = '/blog/' + (state.tag ? '?tag=' + encodeURIComponent(state.tag) : '');
        var html = '<a class="blog-chip' + (allActive ? ' active' : '') + '" href="' + allHref + '">全部</a>';
        html += categories.map(function (c) {
            var active = state.category === c.name ? ' active' : '';
            var href = '/blog/?category=' + encodeURIComponent(c.name);
            if (state.tag) {
                href += '&tag=' + encodeURIComponent(state.tag);
            }
            return '<a class="blog-chip' + active + '" href="' + href + '">'
                + escapeHtml(c.name)
                + '<span class="blog-chip__count">' + c.count + '</span>'
                + '</a>';
        }).join('');
        container.innerHTML = html;
    }

    function renderTags(tags) {
        var container = document.getElementById('tagChips');
        if (!tags || tags.length === 0) {
            container.innerHTML = '<span class="blog-empty">暂无标签</span>';
            return;
        }
        var allActive = !state.tag;
        var html = '<a class="blog-chip' + (allActive ? ' active' : '') + '" href="/blog/' + (state.category ? '?category=' + encodeURIComponent(state.category) : '') + '">全部</a>';
        html += tags.map(function (t) {
            var active = state.tag === t.name ? ' active' : '';
            var href = '/blog/?tag=' + encodeURIComponent(t.name);
            if (state.category) {
                href += '&category=' + encodeURIComponent(state.category);
            }
            return '<a class="blog-chip' + active + '" href="' + href + '">'
                + escapeHtml(t.name)
                + '<span class="blog-chip__count">' + t.count + '</span>'
                + '</a>';
        }).join('');
        container.innerHTML = html;
    }

    // ====== 转义 HTML，防止 XSS ======
    function escapeHtml(str) {
        if (str == null) return '';
        return String(str)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    }

    // ====== 加载文章列表 ======
    function loadPosts() {
        var skeleton = '';
        for (var i = 0; i < 6; i++) {
            skeleton += '<div class="post-card post-card--skeleton">'
                + '<div class="post-card__link">'
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
            showToast('加载失败，请重试');
            document.getElementById('postList').innerHTML =
                '<div class="blog-error">加载失败：' + escapeHtml(err.message) + '</div>';
        });
    }

    // ====== Toast 提示 ======
    function showToast(msg) {
        var existing = document.getElementById('blogToast');
        if (existing) existing.remove();
        var toast = document.createElement('div');
        toast.id = 'blogToast';
        toast.className = 'blog-toast';
        toast.textContent = msg;
        document.body.appendChild(toast);
        setTimeout(function () {
            toast.classList.add('show');
        }, 10);
        setTimeout(function () {
            toast.classList.remove('show');
            setTimeout(function () { toast.remove(); }, 300);
        }, 3000);
    }

    // ====== 文章卡片内分类/标签点击 ======
    function bindPostCardClicks() {
        var container = document.getElementById('postList');
        if (!container) return;
        container.addEventListener('click', function (e) {
            var tagEl = e.target.closest('.post-card__tag');
            var catEl = e.target.closest('.post-card__category');
            if (tagEl) {
                e.preventDefault();
                e.stopPropagation();
                var tag = tagEl.getAttribute('data-tag');
                if (tag) {
                    state.tag = tag;
                    state.category = '';
                    state.page = 1;
                    updateUrlParams();
                    loadPosts();
                }
            } else if (catEl) {
                e.preventDefault();
                e.stopPropagation();
                var cat = catEl.getAttribute('data-category');
                if (cat) {
                    state.category = cat;
                    state.tag = '';
                    state.page = 1;
                    updateUrlParams();
                    loadPosts();
                }
            }
        });
    }

    // ====== 更新 URL 参数 ======
    function updateUrlParams() {
        var params = new URLSearchParams();
        if (state.category) params.set('category', state.category);
        if (state.tag) params.set('tag', state.tag);
        if (state.page > 1) params.set('page', state.page);
        var url = '/blog/' + (params.toString() ? '?' + params.toString() : '');
        window.history.replaceState({}, '', url);
    }

    // ====== 加载本地随机背景壁纸（会话内保持不变） ======
    function loadBingBackground() {
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

    // ====== 分页按钮绑定 ======
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
        // 页码点击（事件委托）
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

    // ====== 搜索（防抖） ======
    function bindSearch() {
        var input = document.getElementById('searchInput');
        if (!input) return;
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

    // ====== 返回顶部 ======
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

    // ====== 管理按钮滚动隐藏 ======
    function bindAdminButton() {
        var adminBtn = document.querySelector('.admin-entry-btn');
        if (!adminBtn) return;

        var lastScroll = 0;
        var threshold = 80;
        window.addEventListener('scroll', function () {
            var currentScroll = window.scrollY;
            if (currentScroll < threshold) {
                adminBtn.classList.remove('hidden');
            } else if (currentScroll > lastScroll + 10) {
                adminBtn.classList.add('hidden');
            } else if (currentScroll < lastScroll - 10) {
                adminBtn.classList.remove('hidden');
            }
            lastScroll = currentScroll;
        });
    }

    // ====== 初始化 ======
    document.addEventListener('DOMContentLoaded', function () {
        readQueryParams();
        bindPagination();
        bindSearch();
        bindBackToTop();
        bindAdminButton();
        bindPostCardClicks();

        var catPromise = BlogAPI.getCategories();
        var tagPromise = BlogAPI.getTags();

        catPromise.then(renderCategories).catch(function () {
            document.getElementById('categoryChips').innerHTML = '<span class="blog-empty">加载失败</span>';
        });
        tagPromise.then(renderTags).catch(function () {
            document.getElementById('tagChips').innerHTML = '<span class="blog-empty">加载失败</span>';
        });
        Promise.all([
            catPromise.catch(function() { return []; }),
            tagPromise.catch(function() { return []; })
        ]).then(function() {
            loadPosts();
        });
        loadBingBackground();
    });
})();
