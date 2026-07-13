/**
 * 博客详情页逻辑
 * 负责：读取 ?id= 参数 → 加载文章 → 渲染 markdown → 代码高亮 → Bing 背景
 */
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

    // 渲染文章主体
    function renderPost(post) {
        var container = document.getElementById('postDetail');

        // 配置 marked
        if (window.marked) {
            marked.setOptions({
                breaks: true,
                gfm: true,
                highlight: function (code, lang) {
                    if (window.hljs && lang && hljs.getLanguage(lang)) {
                        try {
                            return hljs.highlight(code, { language: lang }).value;
                        } catch (e) {}
                    }
                    return code;
                }
            });
        }

        var content = post.content || '';
        // 去除开头重复的 # 标题（避免和页面标题重复显示）
        content = content.replace(/^\s*#\s+[^\n]+\n+/, '');

        var contentHtml = window.marked
            ? marked.parse(content)
            : '<pre>' + escapeHtml(content) + '</pre>';

        var tagsHtml = (post.tags || []).map(function (t) {
            return '<a class="blog-chip" href="/blog/?tag=' + encodeURIComponent(t) + '">'
                + escapeHtml(t) + '</a>';
        }).join('');

        container.innerHTML = ''
            + '<div>'
            + '  <span class="post-detail__category">' + escapeHtml(post.category || '未分类') + '</span>'
            + '</div>'
            + '<h1 class="post-detail__title">' + escapeHtml(post.title) + '</h1>'
            + '<div class="post-detail__meta">'
            + '  <span>&#128197; ' + (window.formatDate(post.created_at) || '') + '</span>'
            + (post.reading_time ? '<span>&#9201; ' + post.reading_time + ' 分钟阅读</span>' : '')
            + (post.author ? '<span>&#9997; ' + escapeHtml(post.author) + '</span>' : '')
            + '  <div class="post-detail__tags">' + tagsHtml + '</div>'
            + '</div>'
            + '<div class="post-content" id="postContent">' + contentHtml + '</div>'
            + '<nav id="postNav" class="post-nav"></nav>';

        // 代码高亮
        if (window.hljs && window.marked) {
            document.querySelectorAll('.post-content pre code').forEach(function (block) {
                hljs.highlightElement(block);
            });
        }

        // 图片懒加载
        document.querySelectorAll('.post-content img').forEach(function (img) {
            img.loading = 'lazy';
        });

        // 添加代码块复制按钮
        addCopyButtons();

        // 生成目录
        generateToc();

        // 阅读进度条
        bindReadingProgress();

        // 返回顶部
        bindBackToTop();

        // 管理按钮滚动隐藏
        bindAdminButton();

        // 加载上下篇导航
        loadPostNav(post.id);
    }

    // 代码块复制按钮
    function addCopyButtons() {
        var pres = document.querySelectorAll('.post-content pre');
        pres.forEach(function (pre) {
            var wrapper = document.createElement('div');
            wrapper.className = 'code-block-wrapper';
            pre.parentNode.insertBefore(wrapper, pre);
            wrapper.appendChild(pre);

            var btn = document.createElement('button');
            btn.className = 'code-copy-btn';
            btn.textContent = '复制';
            btn.addEventListener('click', function () {
                var code = pre.querySelector('code');
                var text = code ? code.textContent : pre.textContent;
                if (navigator.clipboard) {
                    navigator.clipboard.writeText(text).then(function () {
                        btn.textContent = '已复制';
                        btn.classList.add('copied');
                        setTimeout(function () {
                            btn.textContent = '复制';
                            btn.classList.remove('copied');
                        }, 2000);
                    });
                } else {
                    var ta = document.createElement('textarea');
                    ta.value = text;
                    document.body.appendChild(ta);
                    ta.select();
                    document.execCommand('copy');
                    document.body.removeChild(ta);
                    btn.textContent = '已复制';
                    btn.classList.add('copied');
                    setTimeout(function () {
                        btn.textContent = '复制';
                        btn.classList.remove('copied');
                    }, 2000);
                }
            });
            wrapper.appendChild(btn);
        });
    }

    // 生成文章目录 TOC
    function generateToc() {
        var headings = document.querySelectorAll('.post-content h2, .post-content h3');
        var toc = document.getElementById('postToc');
        var overlay = document.getElementById('postTocOverlay');
        if (headings.length < 2) {
            if (toc) toc.style.display = 'none';
            return;
        }

        var html = '<p class="post-toc__title">此页内容</p><ul class="post-toc__list">';
        headings.forEach(function (h, i) {
            var id = 'heading-' + i;
            h.id = id;
            var level = h.tagName === 'H2' ? '' : ' post-toc__item--h3';
            html += '<li class="post-toc__item' + level + '">'
                + '<a class="post-toc__link" href="#' + id + '" data-target="' + id + '">' + escapeHtml(h.textContent) + '</a>'
                + '</li>';
        });
        html += '</ul>';
        toc.innerHTML = html;
        toc.style.display = 'block';

        function isWideScreen() {
            return window.innerWidth > 1280;
        }

        toc.querySelectorAll('.post-toc__link').forEach(function (link) {
            link.addEventListener('click', function (e) {
                e.preventDefault();
                var targetId = this.getAttribute('data-target');
                var target = document.getElementById(targetId);
                if (target) {
                    var offset = target.getBoundingClientRect().top + window.pageYOffset - 90;
                    window.scrollTo({ top: offset, behavior: 'smooth' });
                }
            });
        });

        function checkScreen() {
            if (isWideScreen()) {
                toc.classList.remove('collapsed', 'mobile-open');
                if (overlay) overlay.classList.remove('open');
            } else {
                toc.classList.add('collapsed');
                if (overlay) overlay.classList.remove('open');
            }
        }

        checkScreen();
        window.addEventListener('resize', checkScreen);

        var links = toc.querySelectorAll('.post-toc__link');
        window.addEventListener('scroll', function () {
            var scrollY = window.scrollY + 120;
            var activeIdx = 0;
            headings.forEach(function (h, i) {
                if (h.offsetTop <= scrollY) activeIdx = i;
            });
            links.forEach(function (link, i) {
                link.classList.toggle('active', i === activeIdx);
            });
        });

        window.closeToc = function() {};
    }

    // 阅读进度条
    function bindReadingProgress() {
        var bar = document.getElementById('readingProgress');
        if (!bar) return;
        window.addEventListener('scroll', function () {
            var scrollTop = window.scrollY;
            var docHeight = document.documentElement.scrollHeight - window.innerHeight;
            var progress = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0;
            bar.style.width = Math.min(progress, 100) + '%';
        });
    }

    // 返回顶部
    function bindBackToTop() {
        var btn = document.getElementById('backToTop');
        if (!btn) return;
        function check() {
            btn.classList.toggle('visible', window.scrollY > 400);
        }
        check();
        window.addEventListener('scroll', check);
        btn.addEventListener('click', function () {
            window.scrollTo({ top: 0, behavior: 'smooth' });
        });
    }

    // 管理按钮滚动隐藏
    function bindAdminButton() {
        var adminBtn = document.querySelector('.admin-entry-btn');
        if (!adminBtn) return;

        var lastScroll = 0;
        var threshold = 80;
        function check() {
            var currentScroll = window.scrollY;
            if (currentScroll < threshold) {
                adminBtn.classList.remove('hidden');
            } else if (currentScroll > lastScroll + 10) {
                adminBtn.classList.add('hidden');
            } else if (currentScroll < lastScroll - 10) {
                adminBtn.classList.remove('hidden');
            }
            lastScroll = currentScroll;
        }
        check();
        window.addEventListener('scroll', check);
    }

    // 上下篇导航
    function loadPostNav(currentId) {
        var nav = document.getElementById('postNav');
        if (!nav) return;
        BlogAPI.getAdjacentPosts(currentId).then(function (data) {
            var html = '';
            if (data.prev) {
                html += '<a class="post-nav__link" href="/blog/post/' + encodeURIComponent(data.prev.id) + '">'
                    + '<span class="post-nav__label">&larr; 上一篇</span>'
                    + '<span class="post-nav__title">' + escapeHtml(data.prev.title) + '</span>'
                    + '</a>';
            } else {
                html += '<div class="post-nav__placeholder"></div>';
            }
            if (data.next) {
                html += '<a class="post-nav__link post-nav__link--next" href="/blog/post/' + encodeURIComponent(data.next.id) + '">'
                    + '<span class="post-nav__label">下一篇 &rarr;</span>'
                    + '<span class="post-nav__title">' + escapeHtml(data.next.title) + '</span>'
                    + '</a>';
            } else {
                html += '<div class="post-nav__placeholder"></div>';
            }
            nav.innerHTML = html;
        }).catch(function () {
            nav.innerHTML = '';
        });
    }

    function renderError(msg) {
        showToast(msg);
        document.getElementById('postDetail').innerHTML =
            '<div class="blog-error">' + escapeHtml(msg) + '</div>'
            + '<p style="text-align:center;margin-top:20px;">'
            + '<a href="/blog/" class="post-detail__back">&larr; 返回列表</a>'
            + '</p>';
    }

    // Toast 提示
    function showToast(msg) {
        var existing = document.getElementById('blogToast');
        if (existing) existing.remove();
        var toast = document.createElement('div');
        toast.id = 'blogToast';
        toast.className = 'blog-toast';
        toast.textContent = msg;
        document.body.appendChild(toast);
        setTimeout(function () { toast.classList.add('show'); }, 10);
        setTimeout(function () {
            toast.classList.remove('show');
            setTimeout(function () { toast.remove(); }, 300);
        }, 3000);
    }

    // 加载本地随机背景壁纸（会话内保持不变）
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

    // 初始化
    document.addEventListener('DOMContentLoaded', function () {
        var path = window.location.pathname;
        var id = '';
        var postPathMatch = path.match(/^\/blog\/post\/(\d+)/);
        if (postPathMatch) {
            id = postPathMatch[1];
        }
        if (!id) {
            var params = new URLSearchParams(window.location.search);
            id = params.get('id');
        }

        if (!id) {
            window.location.href = '/blog/';
            return;
        }

        Promise.all([
            BlogAPI.getCategories().catch(function() { return []; }),
            BlogAPI.getTags().catch(function() { return []; })
        ]).then(function() {
            return BlogAPI.getPost(id);
        }).then(renderPost).catch(function (err) {
            renderError('文章加载失败：' + err.message);
        });

        loadBingBackground();
    });
})();
