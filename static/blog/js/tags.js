(function () {
    'use strict';

    function escapeHtml(value) {
        return String(value == null ? '' : value).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;');
    }

    function buildConstellation(tags) {
        var ordered = (tags || []).slice().sort(function (left, right) { return right.count - left.count || left.name.localeCompare(right.name, 'zh-CN'); });
        var maxCount = Math.max.apply(null, ordered.map(function (tag) { return tag.count; }).concat([1]));
        return ordered.map(function (tag, index) {
            var angle = -Math.PI / 2 + index * 2.399963229728653;
            var ring = 17 + Math.sqrt(index + 1) * 11;
            return { name: tag.name, count: tag.count, x: Math.max(7, Math.min(93, 50 + Math.cos(angle) * ring)), y: Math.max(12, Math.min(88, 50 + Math.sin(angle) * ring * 0.62)), scale: 0.82 + tag.count / maxCount * 0.78, delay: index * 45 };
        });
    }

    function renderTagConstellation(items, focus) {
        if (!items.length) return '<p class="directory-state">\u8fd8\u6ca1\u6709\u53ef\u5c55\u793a\u7684\u6807\u7b7e</p>';
        return items.map(function (item) {
            var focused = item.name === focus;
            var style = '--x:' + item.x.toFixed(2) + '%;--y:' + item.y.toFixed(2) + '%;--scale:' + item.scale.toFixed(2) + ';--delay:' + item.delay + 'ms;';
            return '<button class="tag-constellation__node' + (focused ? ' is-focus' : '') + '" type="button" data-tag="' + escapeHtml(item.name) + '" aria-pressed="' + focused + '" style="' + style + '"><span class="tag-constellation__spark" aria-hidden="true"></span><span class="tag-constellation__name">' + escapeHtml(item.name) + '</span><span class="tag-constellation__count">' + item.count + '</span></button>';
        }).join('');
    }

    function articleHasTag(post, name) {
        return (post.tags || []).some(function (tag) { return (tag.name || tag) === name; });
    }

    function renderFocusedArticles(posts, focus) {
        if (!focus) return '<p class="tag-articles__hint">\u9009\u62e9\u4e00\u9897\u6807\u7b7e\u661f\uff0c\u67e5\u770b\u5b83\u8fde\u63a5\u7684\u6587\u7ae0\u3002</p>';
        var matched = posts.filter(function (post) { return articleHasTag(post, focus); });
        if (!matched.length) return '<p class="tag-articles__hint">\u300c' + escapeHtml(focus) + '\u300d\u6682\u65e0\u516c\u5f00\u6587\u7ae0\u3002</p>';
        return '<div class="tag-articles__heading"><span>TAG / ' + escapeHtml(focus) + '</span><strong>' + matched.length + ' \u7bc7\u6587\u7ae0</strong></div><div class="tag-articles__list">' + matched.map(function (post) {
            var date = window.formatDate(post.published_at || post.created_at || post.updated_at);
            return '<a class="tag-articles__item" href="/blog/post/' + encodeURIComponent(post.slug || post.id) + '"><span>' + escapeHtml(post.title) + '</span><time>' + escapeHtml(date) + '</time></a>';
        }).join('') + '</div>';
    }

    function setText(id, value) { var element = document.getElementById(id); if (element) element.textContent = String(value == null ? '' : value); }
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
        var field = document.getElementById('tagConstellation');
        var articles = document.getElementById('tagArticles');
        var summary = document.getElementById('tagsSummary');
        var state = { focus: new URLSearchParams(window.location.search).get('focus') || '', items: [], posts: [] };
        function render() { field.innerHTML = renderTagConstellation(state.items, state.focus); articles.innerHTML = renderFocusedArticles(state.posts, state.focus); }
        field.addEventListener('click', function (event) {
            var node = event.target.closest('[data-tag]');
            if (!node) return;
            state.focus = node.getAttribute('data-tag') || '';
            var url = new URL(window.location.href);
            url.searchParams.set('focus', state.focus);
            history.replaceState(null, '', url.pathname + url.search);
            render();
        });
        Promise.all([BlogAPI.getTags(), BlogAPI.getCategories(), BlogAPI.getPosts({ page: 1, limit: 200 })]).then(function (result) {
            var tags = result[0];
            var categories = result[1];
            var postsData = result[2];
            state.items = buildConstellation(tags);
            state.posts = postsData.posts || [];
            setText('statPosts', postsData.total || state.posts.length);
            setText('statCategories', categories.length);
            setText('statTags', tags.length);
            summary.textContent = '\u5f53\u524d\u5171\u6709 ' + tags.length + ' \u4e2a\u6807\u7b7e\uff0c\u70b9\u4eae\u5176\u4e2d\u4e00\u9897\u661f\u67e5\u770b\u6587\u7ae0';
            render();
        }).catch(function (error) {
            summary.textContent = '\u6807\u7b7e\u661f\u56fe\u52a0\u8f7d\u5931\u8d25';
            field.innerHTML = '<p class="directory-state directory-state--error">' + escapeHtml(error.message) + '</p>';
        });
    });
}());