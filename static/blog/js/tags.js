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

    function renderTagConstellation(items) {
        if (!items.length) return '<p class="directory-state">还没有可展示的标签</p>';
        return items.map(function (item) {
            var style = '--x:' + item.x.toFixed(2) + '%;--y:' + item.y.toFixed(2) + '%;--scale:' + item.scale.toFixed(2) + ';--delay:' + item.delay + 'ms;';
            return '<button class="tag-constellation__node" type="button" data-tag="' + escapeHtml(item.name) + '" style="' + style + '"><span class="tag-constellation__spark" aria-hidden="true"></span><span class="tag-constellation__name">' + escapeHtml(item.name) + '</span><span class="tag-constellation__count">' + item.count + '</span></button>';
        }).join('');
    }

    function articleHasTag(post, name) {
        return (post.tags || []).some(function (tag) { return (tag.name || tag) === name; });
    }

    // The collection endpoint can omit tags, so archive mode resolves only missing records.
    function hydrateTagPosts(posts) {
        var missingTags = (posts || []).filter(function (post) {
            return post.id && (!post.tags || post.tags.length === 0);
        });
        if (!missingTags.length) return Promise.resolve(posts);
        return Promise.all(missingTags.map(function (post) {
            return BlogAPI.getPost(post.id).then(function (detail) {
                post.tags = detail.tags || [];
                return post;
            }).catch(function () {
                return post;
            });
        })).then(function () {
            return posts;
        });
    }

    function renderTagArchive(focus, posts) {
        var matched = (posts || []).filter(function (post) { return articleHasTag(post, focus); });
        var heading = '<div class="tag-archive__header"><button class="tag-archive__return" type="button" data-tag-return>← 返回 Tag.sort()</button><p class="tag-archive__eyebrow">TAG ARCHIVE / ' + escapeHtml(focus) + '</p><h3 class="tag-archive__title">' + escapeHtml(focus) + '</h3><span class="tag-archive__count">' + matched.length + ' 篇文章</span></div>';
        if (!matched.length) return heading + '<p class="tag-archive__state">「' + escapeHtml(focus) + '」暂无公开文章。</p>';
        return heading + '<div class="tag-archive__list">' + matched.map(function (post, index) {
            var date = window.formatDate(post.published_at || post.created_at || post.updated_at);
            var relatedTags = (post.tags || []).filter(function (tag) { return (tag.name || tag) !== focus; }).map(function (tag) {
                var name = tag.name || tag;
                return '<button type="button" class="tag-archive__chip" data-tag-switch="' + escapeHtml(name) + '">' + escapeHtml(name) + '</button>';
            }).join('');
            return '<article class="tag-archive__item" style="--archive-delay:' + Math.min(index * 70, 420) + 'ms"><a class="tag-archive__article-link" href="/blog/post/' + encodeURIComponent(post.id) + '"><h4>' + escapeHtml(post.title) + '</h4><p>' + escapeHtml(post.summary || '') + '</p></a><footer><time>' + escapeHtml(date) + '</time><div class="tag-archive__chips">' + relatedTags + '</div></footer></article>';
        }).join('') + '</div>';
    }

    function animateConstellationNodes(container) {
        var elements = Array.prototype.slice.call(container.querySelectorAll('.tag-constellation__node'));
        if (!elements.length) return;
        if (window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
            elements.forEach(function (element) { element.classList.add('is-revealed'); });
            return;
        }
        elements.forEach(function (element, index) {
            element.classList.add('is-revealing');
            element.style.setProperty('--reveal-delay', Math.min(index * 120, 2160) + 'ms');
        });
        window.requestAnimationFrame(function () {
            elements.forEach(function (element) { element.classList.add('is-revealed'); });
        });
    }

    function setText(id, value) {
        var element = document.getElementById(id);
        if (element) element.textContent = String(value == null ? '' : value);
    }

    document.addEventListener('DOMContentLoaded', function () {
        var field = document.getElementById('tagConstellation');
        var constellation = document.querySelector('.tag-constellation');
        var articles = document.getElementById('tagArticles');
        var summary = document.getElementById('tagsSummary');
        var state = { focus: new URLSearchParams(window.location.search).get('focus') || '', items: [], posts: null, loading: false };

        function updateUrl() {
            var url = new URL(window.location.href);
            if (state.focus) url.searchParams.set('focus', state.focus);
            else url.searchParams.delete('focus');
            history.replaceState(null, '', url.pathname + url.search);
        }

        function render() {
            var archiveMode = Boolean(state.focus);
            document.body.classList.toggle('blog-page--tag-archive', archiveMode);
            constellation.hidden = archiveMode;
            if (!archiveMode) {
                field.innerHTML = renderTagConstellation(state.items);
                articles.innerHTML = '<p class="tag-articles__hint">选择一颗标签星，查看它连接的文章。</p>';
                animateConstellationNodes(field);
                return;
            }
            if (state.loading) {
                articles.innerHTML = '<p class="tag-archive__state">正在整理「' + escapeHtml(state.focus) + '」的文章...</p>';
                return;
            }
            articles.innerHTML = renderTagArchive(state.focus, state.posts || []);
        }

        function selectTag(name) {
            state.focus = name || '';
            updateUrl();
            if (!state.focus) {
                render();
                return;
            }
            if (state.posts) {
                render();
                return;
            }
            state.loading = true;
            render();
            BlogAPI.getPosts({ page: 1, limit: 200 }).then(function (data) {
                return hydrateTagPosts(data.posts || []);
            }).then(function (posts) {
                state.posts = posts;
                state.loading = false;
                render();
            }).catch(function () {
                state.posts = [];
                state.loading = false;
                render();
            });
        }

        field.addEventListener('click', function (event) {
            var node = event.target.closest('[data-tag]');
            if (node) selectTag(node.getAttribute('data-tag') || '');
        });
        articles.addEventListener('click', function (event) {
            var returnButton = event.target.closest('[data-tag-return]');
            if (returnButton) {
                selectTag('');
                return;
            }
            var tagButton = event.target.closest('[data-tag-switch]');
            if (tagButton) selectTag(tagButton.getAttribute('data-tag-switch') || '');
        });

        Promise.all([BlogAPI.getTags(), BlogAPI.getCategories(), BlogAPI.getPosts({ page: 1, limit: 1 })]).then(function (result) {
            var tags = result[0];
            var categories = result[1];
            var postsData = result[2];
            state.items = buildConstellation(tags);
            setText('statPosts', postsData.total || 0);
            setText('statCategories', categories.length);
            setText('statTags', tags.length);
            summary.textContent = '当前共有 ' + tags.length + ' 个标签，点亮其中一颗星查看文章';
            render();
            if (state.focus) selectTag(state.focus);
        }).catch(function (error) {
            summary.textContent = '标签星图加载失败';
            field.innerHTML = '<p class="directory-state directory-state--error">' + escapeHtml(error.message) + '</p>';
        });
    });
}());