(function () {
    'use strict';

    function escapeHtml(str) {
        if (str == null) return '';
        return String(str)
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    }

    function slugify(text) {
        return String(text || 'uncategorized').toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/^-+|-+$/g, '');
    }

    function groupPostsByCategory(posts) {
        var groups = {};
        posts.forEach(function (post) {
            var key = post.category || '未分类';
            if (!groups[key]) groups[key] = [];
            groups[key].push(post);
        });
        return groups;
    }

    function renderCategoryPostLink(post) {
        return '<a class="category-group__post-link" href="/blog/post/' + encodeURIComponent(post.id) + '">'
            + '<span class="category-group__post-title">' + escapeHtml(post.title) + '</span>'
            + '<span class="category-group__post-date">' + escapeHtml(window.formatDate(post.created_at) || '') + '</span>'
            + '</a>';
    }

    function renderCategoryGroups(groups, focus) {
        var names = Object.keys(groups).sort(function (a, b) {
            return groups[b].length - groups[a].length || a.localeCompare(b, 'zh-CN');
        });

        return names.map(function (name) {
            var isFocus = focus && focus === name;
            return '<section class="category-group' + (isFocus ? ' is-focus' : '') + '" id="category-' + slugify(name) + '">'
                + '<header class="category-group__head">'
                + '<h3 class="category-group__title">' + escapeHtml(name) + '</h3>'
                + '<span class="category-group__count">' + groups[name].length + ' 篇</span>'
                + '</header>'
                + '<div class="category-group__posts">' + groups[name].map(renderCategoryPostLink).join('') + '</div>'
                + '</section>';
        }).join('');
    }

    function updateStats(postsTotal, categoriesTotal, tagsTotal) {
        setText('statPosts', postsTotal);
        setText('statCategories', categoriesTotal);
        setText('statTags', tagsTotal);
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

    function scrollToFocus(focus) {
        if (!focus) return;
        var target = document.getElementById('category-' + slugify(focus));
        if (target) {
            setTimeout(function () {
                target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            }, 120);
        }
    }

    document.addEventListener('DOMContentLoaded', function () {
        loadBackground();
        var params = new URLSearchParams(window.location.search);
        var focus = params.get('focus') || '';
        var summary = document.getElementById('categoriesSummary');
        var groupsEl = document.getElementById('categoryGroups');

        Promise.all([
            BlogAPI.getCategories().catch(function () { return []; }),
            BlogAPI.getTags().catch(function () { return []; }),
            BlogAPI.getPosts({ page: 1, limit: 200 }).catch(function () { return { posts: [], total: 0 }; })
        ]).then(function (result) {
            var categories = result[0];
            var tags = result[1];
            var postsData = result[2];
            var groups = groupPostsByCategory(postsData.posts || []);

            updateStats(postsData.total || (postsData.posts || []).length, categories.length, tags.length);
            summary.textContent = '当前共有 ' + categories.length + ' 个分类，收录 ' + (postsData.total || (postsData.posts || []).length) + ' 篇文章';
            groupsEl.innerHTML = renderCategoryGroups(groups, focus);
            scrollToFocus(focus);
        }).catch(function (err) {
            summary.textContent = '分类加载失败';
            groupsEl.innerHTML = '<div class="blog-error">加载失败：' + escapeHtml(err.message) + '</div>';
        });
    });
})();
