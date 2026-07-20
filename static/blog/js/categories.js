(function () {
    'use strict';

    var TEXT = {
        uncategorized: '\u672a\u5206\u7c7b',
        articles: '\u7bc7',
        empty: '\u8fd8\u6ca1\u6709\u53ef\u5c55\u793a\u7684\u5206\u7c7b\u6587\u7ae0',
        failed: '\u5206\u7c7b\u76ee\u5f55\u52a0\u8f7d\u5931\u8d25',
        end: '-- \u5df2\u7ecf\u5230\u5e95\u4e86 --'
    };

    function escapeHtml(value) {
        return String(value == null ? '' : value).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;');
    }

    function categoryKey(name) {
        return String(name || TEXT.uncategorized).toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/^-+|-+$/g, '') || 'uncategorized';
    }

    function groupPosts(posts, categories) {
        var groups = {};
        (categories || []).forEach(function (category) { groups[category.name || TEXT.uncategorized] = []; });
        (posts || []).forEach(function (post) {
            var name = post.category || TEXT.uncategorized;
            if (!groups[name]) groups[name] = [];
            groups[name].push(post);
        });
        return groups;
    }

    function renderPost(post) {
        var date = window.formatDate(post.published_at || post.created_at || post.updated_at);
        return '<a class="category-directory__post" href="/blog/post/' + encodeURIComponent(post.slug || post.id) + '"><span class="category-directory__post-title">' + escapeHtml(post.title) + '</span><time class="category-directory__post-date">' + escapeHtml(date) + '</time></a>';
    }

    function renderCategoryDirectory(groups, focus) {
        var names = Object.keys(groups).sort(function (left, right) { return groups[right].length - groups[left].length || left.localeCompare(right, 'zh-CN'); });
        if (!names.length) return '<p class="directory-state">' + TEXT.empty + '</p>';
        return names.map(function (name) {
            var posts = groups[name];
            return '<section id="category-' + categoryKey(name) + '" class="category-directory__group' + (focus === name ? ' is-focus' : '') + '"><header class="category-directory__heading"><span class="category-directory__dot" aria-hidden="true"></span><h3 class="category-directory__title">' + escapeHtml(name) + '</h3><span class="category-directory__count">' + posts.length + ' ' + TEXT.articles + '</span><span class="category-directory__line" aria-hidden="true"></span></header><div class="category-directory__posts">' + posts.map(renderPost).join('') + '</div></section>';
        }).join('') + '<p class="directory-end">' + TEXT.end + '</p>';
    }

    function setText(id, value) { var element = document.getElementById(id); if (element) element.textContent = String(value == null ? '' : value); }
    function setStats(posts, categories, tags) { setText('statPosts', posts); setText('statCategories', categories); setText('statTags', tags); }

    function loadBackground() {
        var script = document.createElement('script');
        script.src = '/assets/json/images.js?t=' + Date.now();
        script.onload = function () {
            if (!window.BING_IMAGES || !window.BING_IMAGES.length) return;
            var index = Number(sessionStorage.getItem('blog-bg-index'));
            if (!Number.isInteger(index) || index < 0 || index >= window.BING_IMAGES.length) { index = Math.floor(Math.random() * window.BING_IMAGES.length); sessionStorage.setItem('blog-bg-index', String(index)); }
            document.body.style.backgroundImage = "url('/assets/" + window.BING_IMAGES[index].replace(/['\\]/g, '\\$&') + "')";
        };
        document.body.appendChild(script);
    }

    document.addEventListener('DOMContentLoaded', function () {
        loadBackground();
        var focus = new URLSearchParams(window.location.search).get('focus') || '';
        var summary = document.getElementById('categoriesSummary');
        var directory = document.getElementById('categoryDirectory');
        Promise.all([BlogAPI.getCategories(), BlogAPI.getTags(), BlogAPI.getPosts({ page: 1, limit: 200 })]).then(function (result) {
            var categories = result[0];
            var tags = result[1];
            var postsData = result[2];
            var posts = postsData.posts || [];
            setStats(postsData.total || posts.length, categories.length, tags.length);
            summary.textContent = '\u5f53\u524d\u5171\u6709 ' + categories.length + ' \u4e2a\u5206\u7c7b\uff0c\u6536\u5f55 ' + (postsData.total || posts.length) + ' \u7bc7\u6587\u7ae0';
            directory.innerHTML = renderCategoryDirectory(groupPosts(posts, categories), focus);
            if (focus) { var target = document.getElementById('category-' + categoryKey(focus)); if (target) setTimeout(function () { target.scrollIntoView({ behavior: 'smooth', block: 'start' }); }, 120); }
        }).catch(function (error) {
            summary.textContent = TEXT.failed;
            directory.innerHTML = '<p class="directory-state directory-state--error">' + TEXT.failed + ': ' + escapeHtml(error.message) + '</p>';
        });
    });
}());