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

    function buildTagLayout(tags) {
        return tags.map(function (tag, index) {
            var column = index % 5;
            var row = Math.floor(index / 5);
            return {
                name: tag.name,
                count: tag.count,
                rotate: ((index % 7) - 3) * 4,
                delay: index * 35,
                x: column * 180 + (row % 2) * 18,
                y: row * 88 + (index % 3) * 10
            };
        });
    }

    function renderTagPile(items, focus) {
        return items.map(function (item) {
            var isFocus = focus && focus === item.name;
            return '<a class="tag-pile__item' + (isFocus ? ' is-focus' : '') + '" href="/blog/?tag=' + encodeURIComponent(item.name) + '"'
                + ' style="--tag-rotate:' + item.rotate + 'deg;--tag-delay:' + item.delay + 'ms;left:' + item.x + 'px;top:' + item.y + 'px;">'
                + '<span class="tag-pile__name">' + escapeHtml(item.name) + '</span>'
                + '<span class="tag-pile__count">(' + item.count + ')</span></a>';
        }).join('');
    }

    function setText(id, value) {
        if (value == null) return;
        var el = document.getElementById(id);
        if (el) el.textContent = String(value);
    }

    function updateStats(postsTotal, categoriesTotal, tagsTotal) {
        setText('statPosts', postsTotal);
        setText('statCategories', categoriesTotal);
        setText('statTags', tagsTotal);
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
        loadBackground();
        var params = new URLSearchParams(window.location.search);
        var focus = params.get('focus') || '';
        var summary = document.getElementById('tagsSummary');
        var pile = document.getElementById('tagPile');

        Promise.all([
            BlogAPI.getTags().catch(function () { return []; }),
            BlogAPI.getCategories().catch(function () { return []; }),
            BlogAPI.getPosts({ page: 1, limit: 200 }).catch(function () { return { posts: [], total: 0 }; })
        ]).then(function (result) {
            var tags = result[0];
            var categories = result[1];
            var postsData = result[2];
            var items = buildTagLayout(tags);

            updateStats(postsData.total || (postsData.posts || []).length, categories.length, tags.length);
            summary.textContent = '当前共有 ' + tags.length + ' 个标签';
            pile.innerHTML = renderTagPile(items, focus);
            if (window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
                pile.classList.add('is-reduced-motion');
            }
        }).catch(function (err) {
            summary.textContent = '标签加载失败';
            pile.innerHTML = '<div class="blog-error">加载失败：' + escapeHtml(err.message) + '</div>';
        });
    });
})();
