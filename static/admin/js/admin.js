(function () {
    'use strict';

    var API_BASE = '/api/admin';
    var toastTimer = null;

    function request(path, options) {
        options = options || {};
        return fetch(API_BASE + path, {
            method: options.method || 'GET',
            headers: Object.assign({
                'Content-Type': 'application/json'
            }, options.headers || {}),
            body: options.body ? JSON.stringify(options.body) : undefined,
            credentials: 'same-origin'
        }).then(function (res) {
            if (res.status === 401) {
                throw new Error('unauthorized');
            }
            if (!res.ok) {
                var apiError = new Error(friendlyHttpError(path, res.status));
                apiError.status = res.status;
                apiError.path = path;
                throw apiError;
            }
            return res.json();
        }).then(function (data) {
            if (data.code !== 0) {
                throw new Error(data.message || '服务暂时没有返回可用结果，请稍后再试');
            }
            return data.data;
        });
    }

    function uploadCoverFile(file, onProgress) {
        if (!file) return Promise.resolve('');
        if (!file.type || file.type.indexOf('image/') !== 0) {
            return Promise.reject(new Error('请选择图片文件作为封面'));
        }

        return new Promise(function (resolve, reject) {
            var xhr = new XMLHttpRequest();
            var formData = new FormData();
            formData.append('file', file);
            formData.append('biz_type', 'article_cover');

            xhr.open('POST', '/api/media/upload', true);
            xhr.withCredentials = true;

            xhr.upload.addEventListener('progress', function (evt) {
                if (!evt.lengthComputable || typeof onProgress !== 'function') return;
                var percent = evt.total > 0 ? Math.round((evt.loaded / evt.total) * 100) : 0;
                onProgress(percent);
            });

            xhr.addEventListener('load', function () {
                if (xhr.status < 200 || xhr.status >= 300) {
                    reject(new Error(friendlyHttpError('/api/media/upload', xhr.status)));
                    return;
                }
                var data;
                try {
                    data = JSON.parse(xhr.responseText || '{}');
                } catch (err) {
                    reject(new Error('封面上传失败，请稍后再试'));
                    return;
                }
                if (data.code !== 0) {
                    reject(new Error(data.message || '封面上传失败，请稍后再试'));
                    return;
                }
                var record = data.data || {};
                if (!record.url) {
                    reject(new Error('封面上传成功，但服务没有返回可访问地址'));
                    return;
                }
                resolve(record.url);
            });

            xhr.addEventListener('error', function () {
                reject(new Error('封面上传失败，请检查网络后重试'));
            });

            xhr.send(formData);
        });
    }

    function friendlyHttpError(path, status) {
        if (status === 404) {
            if (path.indexOf('/articles/') === 0) return '没有找到这篇文章，可能已被删除或接口路径尚未同步。';
            return '没有找到对应的数据，请刷新页面后再试。';
        }
        if (status === 403) return '当前账号没有权限执行这个操作。';
        if (status >= 500) return '服务暂时开小差了，请稍后再试。';
        return '请求失败，请稍后再试。';
    }

    function showToast(message, type) {
        var toast = document.getElementById('adminToast');
        if (!toast) {
            toast = document.createElement('div');
            toast.id = 'adminToast';
            toast.className = 'admin-toast';
            document.body.appendChild(toast);
        }
        toast.textContent = message || '操作失败，请稍后再试。';
        toast.className = 'admin-toast admin-toast--' + (type || 'error') + ' show';
        clearTimeout(toastTimer);
        toastTimer = setTimeout(function () {
            toast.classList.remove('show');
        }, 3600);
    }

    function reportError(prefix, err) {
        if (err && err.message === 'unauthorized') {
            showLogin();
            return;
        }
        showToast(prefix + '：' + ((err && err.message) || '请稍后再试'), 'error');
    }
    function slugify(str) {
        str = str || '';
        return str.toLowerCase()
            .replace(/[^\w\u4e00-\u9fa5\s-]/g, '')
            .replace(/\s+/g, '-')
            .replace(/-+/g, '-')
            .trim('-');
    }

    function metaEndpoint(type) {
        return type === 'category' ? '/categories' : '/tags';
    }

    function generateSlug(title) {
        var s = slugify(title);
        if (!s) s = 'post-' + Date.now();
        return s;
    }

    function formatDate(iso) {
        if (!iso) return '';
        var d = new Date(iso);
        if (isNaN(d.getTime())) {
            return iso;
        }
        var y = d.getFullYear();
        var m = String(d.getMonth() + 1).padStart(2, '0');
        var day = String(d.getDate()).padStart(2, '0');
        var h = String(d.getHours()).padStart(2, '0');
        var mi = String(d.getMinutes()).padStart(2, '0');
        return y + '-' + m + '-' + day + ' ' + h + ':' + mi;
    }

    function escapeHtml(str) {
        if (!str) return '';
        var div = document.createElement('div');
        div.appendChild(document.createTextNode(str));
        return div.innerHTML;
    }


    function setCoverStatus(message) {
        var status = document.getElementById('coverUploadStatus');
        if (status) {
            status.textContent = message || '选择图片后会自动上传，并在保存文章时写入封面地址。';
        }
    }

    function setCoverProgress(percent, message) {
        var panel = document.getElementById('coverUploadProgress');
        var bar = document.getElementById('coverUploadProgressBar');
        var text = document.getElementById('coverUploadProgressText');
        var safePercent = Math.max(0, Math.min(100, parseInt(percent || 0, 10)));
        if (!panel || !bar || !text) return;
        if (safePercent <= 0) {
            panel.style.display = 'none';
            bar.style.width = '0%';
            text.textContent = '';
            return;
        }
        panel.style.display = 'block';
        bar.style.width = safePercent + '%';
        text.textContent = message || ('上传中 ' + safePercent + '%');
    }

    function setCoverImage(url) {
        var input = document.getElementById('postCoverImage');
        var preview = document.getElementById('coverPreview');
        if (input) input.value = url || '';
        if (!preview) return;
        if (!url) {
            preview.style.display = 'none';
            preview.innerHTML = '';
            setCoverStatus('选择图片后会自动上传，并在保存文章时写入封面地址。');
            setCoverProgress(0);
            return;
        }
        preview.style.display = 'block';
        preview.innerHTML = '<img src="' + escapeHtml(url) + '" alt="封面预览">' +
            '<div class="cover-preview__url">' + escapeHtml(url) + '</div>';
        setCoverStatus('封面已上传，保存文章后生效。');
    }

    var cachedCategories = [];
    var cachedTags = [];

    var postList = document.getElementById('postList');
    var statusFilter = 'all';
    var categoryList = document.getElementById('categoryList');
    var tagListEl = document.getElementById('tagList');
    var editorModal = document.getElementById('editorModal');
    var metaEditorModal = document.getElementById('metaEditorModal');
    var confirmModal = document.getElementById('confirmModal');
    var loginModal = document.getElementById('loginModal');
    var loginForm = document.getElementById('loginForm');
    var loginError = document.getElementById('loginError');
    var adminApp = document.getElementById('adminApp');
    var authMask = document.getElementById('authMask');

    function getCategoryNameById(id) {
        for (var i = 0; i < cachedCategories.length; i++) {
            if (cachedCategories[i].id === id) {
                return cachedCategories[i].name;
            }
        }
        return '';
    }

    function getCategoryIdByName(name) {
        for (var i = 0; i < cachedCategories.length; i++) {
            if (cachedCategories[i].name === name) {
                return cachedCategories[i].id;
            }
        }
        return 0;
    }

    function getTagNamesByIds(ids) {
        if (!ids || !Array.isArray(ids)) return [];
        var names = [];
        for (var i = 0; i < ids.length; i++) {
            for (var j = 0; j < cachedTags.length; j++) {
                if (cachedTags[j].id === ids[i]) {
                    names.push(cachedTags[j].name);
                    break;
                }
            }
        }
        return names;
    }

    function getTagIdsByNames(names) {
        var ids = [];
        for (var i = 0; i < names.length; i++) {
            for (var j = 0; j < cachedTags.length; j++) {
                if (cachedTags[j].name === names[i]) {
                    ids.push(cachedTags[j].id);
                    break;
                }
            }
        }
        return ids;
    }

    function countPostsByCategory(catId, posts) {
        var count = 0;
        for (var i = 0; i < posts.length; i++) {
            if (posts[i].category_id === catId) count++;
        }
        return count;
    }

    function countPostsByTag(tagId, posts) {
        var count = 0;
        for (var i = 0; i < posts.length; i++) {
            var tagIds = posts[i].tag_ids || [];
            for (var j = 0; j < tagIds.length; j++) {
                if (tagIds[j] === tagId) {
                    count++;
                    break;
                }
            }
        }
        return count;
    }

    function mapArticle(a) {
        return {
            id: a.id,
            title: a.title,
            slug: a.slug,
            summary: a.summary,
            content_md: a.content_md,
            content_html: a.content_html,
            category_id: a.category_id,
            tag_ids: a.tag_ids || [],
            status: a.status,
            created_at: a.created_at,
            updated_at: a.updated_at,
            published_at: a.published_at,
            cover_image: a.cover_image || ''
        };
    }

    function fetchAllPosts() {
        var allPosts = [];
        var page = 1;
        var pageSize = 100;

        function fetchPage() {
            return request('/articles?page=' + page + '&page_size=' + pageSize).then(function (data) {
                var list = data.list || [];
                allPosts = allPosts.concat(list.map(mapArticle));
                if (data.total > page * pageSize) {
                    page++;
                    return fetchPage();
                }
                return allPosts;
            });
        }
        return fetchPage();
    }

    function loadPost(id) {
        return request('/articles/' + encodeURIComponent(id)).then(function (data) {
            return mapArticle(data || data.article || {});
        });
    }

    function loadCategories() {
        return request('/categories').then(function (data) {
            cachedCategories = (data.list || []).map(function (c) {
                return {
                    id: c.id,
                    name: c.name,
                    slug: c.slug,
                    description: c.description
                };
            });
            return cachedCategories;
        });
    }

    function loadTags() {
        return request('/tags').then(function (data) {
            cachedTags = (data.list || []).map(function (t) {
                return {
                    id: t.id,
                    name: t.name,
                    slug: t.slug,
                    description: t.description
                };
            });
            return cachedTags;
        });
    }

    function renderPosts(posts) {
        var sorted = posts.slice().sort(function (a, b) {
            return new Date(b.created_at) - new Date(a.created_at);
        });

        if (statusFilter !== 'all') {
            sorted = sorted.filter(function (p) {
                if (statusFilter === 'draft') return p.status === 'draft';
                return p.status !== 'draft';
            });
        }

        if (sorted.length === 0) {
            postList.innerHTML = '<div class="admin-empty"><p>还没有文章</p><p style="font-size:.85em;">点击右上角「+ 新建文章」开始创作</p></div>';
            return;
        }

        postList.innerHTML = sorted.map(function (p) {
            var tagNames = getTagNamesByIds(p.tag_ids);
            var tags = tagNames.map(function (t) {
                return '<span class="post-manage-item__tag">' + escapeHtml(t) + '</span>';
            }).join('');
            var catName = getCategoryNameById(p.category_id);
            var cat = catName ? '<span class="post-manage-item__category">' + escapeHtml(catName) + '</span>' : '';
            var isDraft = p.status === 'draft';
            var statusBadge = isDraft
                ? '<span class="post-status-badge post-status-badge--draft">草稿</span>'
                : '<span class="post-status-badge post-status-badge--published">已发布</span>';
            var titleClass = isDraft ? 'post-manage-item__title post-manage-item__title--draft' : 'post-manage-item__title';
            return '<article class="post-manage-item">' +
                '<div class="post-manage-item__info">' +
                    cat +
                    '<h3 class="' + titleClass + '">' + escapeHtml(p.title || '无标题') + statusBadge + '</h3>' +
                    '<p class="post-manage-item__summary">' + escapeHtml(p.summary || '暂无摘要') + '</p>' +
                    '<div class="post-manage-item__meta">' +
                        '<span>创建于 ' + (formatDate(p.created_at) || '刚刚') + '</span>' +
                    '</div>' +
                    (tags ? '<div class="post-manage-item__tags">' + tags + '</div>' : '') +
                '</div>' +
                '<div class="post-manage-item__actions">' +
                    '<button class="btn btn-ghost btn-sm" data-action="edit" data-id="' + p.id + '">编辑</button>' +
                    '<button class="btn btn-danger btn-sm" data-action="delete" data-id="' + p.id + '" data-title="' + escapeHtml(p.title || '') + '">删除</button>' +
                '</div>' +
            '</article>';
        }).join('');
    }

    function renderCategories(posts) {
        if (cachedCategories.length === 0) {
            categoryList.innerHTML = '<div class="admin-empty"><p>还没有分类</p><p style="font-size:.85em;">点击右上角「+ 新建分类」添加</p></div>';
            return;
        }
        categoryList.innerHTML = cachedCategories.map(function (c) {
            var count = countPostsByCategory(c.id, posts);
            return '<div class="tag-manage-item">' +
                '<div>' +
                    '<div class="tag-manage-item__name">' + escapeHtml(c.name) + '</div>' +
                    '<div class="tag-manage-item__count">' + count + ' 篇文章</div>' +
                '</div>' +
                '<div class="tag-manage-item__actions">' +
                    '<button class="btn btn-ghost btn-sm" data-action="edit-meta" data-type="category" data-id="' + c.id + '" data-name="' + escapeHtml(c.name) + '">编辑</button>' +
                    '<button class="btn btn-danger btn-sm" data-action="delete-meta" data-type="category" data-id="' + c.id + '" data-name="' + escapeHtml(c.name) + '" data-count="' + count + '">删除</button>' +
                '</div>' +
            '</div>';
        }).join('');
    }

    function renderTags(posts) {
        if (cachedTags.length === 0) {
            tagListEl.innerHTML = '<div class="admin-empty"><p>还没有标签</p><p style="font-size:.85em;">点击右上角「+ 新建标签」添加</p></div>';
            return;
        }
        tagListEl.innerHTML = cachedTags.map(function (t) {
            var count = countPostsByTag(t.id, posts);
            return '<div class="tag-manage-item">' +
                '<div>' +
                    '<div class="tag-manage-item__name">' + escapeHtml(t.name) + '</div>' +
                    '<div class="tag-manage-item__count">' + count + ' 篇文章</div>' +
                '</div>' +
                '<div class="tag-manage-item__actions">' +
                    '<button class="btn btn-ghost btn-sm" data-action="edit-meta" data-type="tag" data-id="' + t.id + '" data-name="' + escapeHtml(t.name) + '">编辑</button>' +
                    '<button class="btn btn-danger btn-sm" data-action="delete-meta" data-type="tag" data-id="' + t.id + '" data-name="' + escapeHtml(t.name) + '" data-count="' + count + '">删除</button>' +
                '</div>' +
            '</div>';
        }).join('');
    }

    function populateCategoryCheckboxes(selectedId) {
        var container = document.getElementById('postCategoryCheckboxes');
        if (cachedCategories.length === 0) {
            container.innerHTML = '<span style="color:rgba(255,255,255,0.4);font-size:.9em;">请先在「分类」面板创建分类</span>';
            return;
        }
        container.innerHTML = cachedCategories.map(function (c) {
            var checked = selectedId === c.id;
            return '<label class="tag-checkbox category-checkbox' + (checked ? ' checked' : '') + '">' +
                '<input type="radio" name="postCategory" value="' + c.id + '"' + (checked ? ' checked' : '') + '>' +
                '<span>' + escapeHtml(c.name) + '</span>' +
            '</label>';
        }).join('');
    }

    function populateTagCheckboxes(selectedIds) {
        var container = document.getElementById('postTagCheckboxes');
        selectedIds = selectedIds || [];
        if (cachedTags.length === 0) {
            container.innerHTML = '<span style="color:rgba(255,255,255,0.4);font-size:.9em;">请先在「标签」面板创建标签</span>';
            return;
        }
        container.innerHTML = cachedTags.map(function (t) {
            var checked = selectedIds.indexOf(t.id) !== -1;
            return '<label class="tag-checkbox' + (checked ? ' checked' : '') + '">' +
                '<input type="checkbox" value="' + t.id + '"' + (checked ? ' checked' : '') + '>' +
                '<span>' + escapeHtml(t.name) + '</span>' +
            '</label>';
        }).join('');
    }

    var currentTab = 'posts';

    function switchTab(tabName) {
        currentTab = tabName;
        document.querySelectorAll('.admin-tab').forEach(function (btn) {
            btn.classList.toggle('active', btn.dataset.tab === tabName);
        });
        document.querySelectorAll('.admin-panel').forEach(function (panel) {
            panel.classList.toggle('active', panel.id === 'tab-' + tabName);
        });
        refreshAll();
    }

    function refreshAll() {
        postList.innerHTML = '<div class="admin-loading">加载中...</div>';
        categoryList.innerHTML = '<div class="admin-loading">加载中...</div>';
        tagListEl.innerHTML = '<div class="admin-loading">加载中...</div>';

        Promise.all([loadCategories(), loadTags(), fetchAllPosts()]).then(function (results) {
            var posts = results[2];
            renderPosts(posts);
            renderCategories(posts);
            renderTags(posts);
        }).catch(function (err) {
            if (err.message === 'unauthorized') {
                showLogin();
            } else {
                postList.innerHTML = '<div class="admin-error">加载失败：' + escapeHtml(err.message) + '</div>';
                categoryList.innerHTML = '<div class="admin-error">加载失败：' + escapeHtml(err.message) + '</div>';
                tagListEl.innerHTML = '<div class="admin-error">加载失败：' + escapeHtml(err.message) + '</div>';
            }
        });
    }

    var editingPostId = null;

    function openEditor(postId) {
        editingPostId = postId;
        var modalTitle = document.getElementById('modalTitle');
        var form = document.getElementById('postForm');
        form.reset();
        document.getElementById('postId').value = '';
        populateCategoryCheckboxes(0);
        populateTagCheckboxes([]);

        if (postId) {
            modalTitle.textContent = '编辑文章';
            document.getElementById('postId').value = postId;
            loadPost(postId).then(function (post) {
                document.getElementById('postTitle').value = post.title || '';
                populateCategoryCheckboxes(post.category_id);
                populateTagCheckboxes(post.tag_ids || []);
                document.getElementById('postSummary').value = post.summary || '';
                setCoverImage(post.cover_image || '');
                document.getElementById('postContent').value = post.content_md || '';
                var status = post.status || 'published';
                var radio = document.querySelector('input[name="postStatus"][value="' + status + '"]');
                if (radio) radio.checked = true;
            }).catch(function (err) {
                reportError('加载文章失败', err);
            });
        } else {
            modalTitle.textContent = '新建文章';
            var pubRadio = document.querySelector('input[name="postStatus"][value="published"]');
            if (pubRadio) pubRadio.checked = true;
        }

        editorModal.style.display = 'flex';
        document.body.style.overflow = 'hidden';
    }

    function closeEditor() {
        editorModal.style.display = 'none';
        document.body.style.overflow = '';
        editingPostId = null;
        var previewPane = document.getElementById('previewPane');
        var coverFileInput = document.getElementById('postCoverFile');
        var clearCoverBtn = document.getElementById('clearCoverBtn');
        if (coverFileInput) {
            coverFileInput.addEventListener('change', function (e) {
                var file = e.target.files && e.target.files[0];
                if (!file) return;
                setCoverProgress(1, '开始上传...');
                setCoverStatus('封面上传中...');
                uploadCoverFile(file, function (percent) {
                    setCoverProgress(percent, '上传中 ' + percent + '%');
                }).then(function (url) {
                    setCoverImage(url);
                    setCoverProgress(100, '上传完成');
                    showToast('封面上传成功', 'success');
                    setTimeout(function () {
                        setCoverProgress(0);
                    }, 900);
                }).catch(function (err) {
                    setCoverImage('');
                    setCoverProgress(0);
                    reportError('封面上传失败', err);
                }).finally(function () {
                    coverFileInput.value = '';
                });
            });
        }
        if (clearCoverBtn) {
            clearCoverBtn.addEventListener('click', function () {
                setCoverImage('');
                setCoverProgress(0);
                if (coverFileInput) coverFileInput.value = '';
            });
        }

        var previewBtn = document.getElementById('previewToggleBtn');
        if (previewPane) previewPane.style.display = 'none';
        if (previewBtn) {
            previewBtn.textContent = '预览';
            previewBtn.classList.add('btn-ghost');
            previewBtn.classList.remove('btn-primary');
        }
    }

    var editingMetaId = null;
    var editingMetaType = null;

    function openMetaEditor(type, id, name) {
        editingMetaId = id;
        editingMetaType = type;
        var modalTitle = document.getElementById('metaModalTitle');
        var nameLabel = document.getElementById('metaNameLabel');
        document.getElementById('metaType').value = type;
        document.getElementById('metaOldName').value = name || '';
        document.getElementById('metaName').value = name || '';

        if (type === 'category') {
            modalTitle.textContent = id ? '编辑分类' : '新建分类';
            nameLabel.textContent = '分类名称';
        } else {
            modalTitle.textContent = id ? '编辑标签' : '新建标签';
            nameLabel.textContent = '标签名称';
        }

        metaEditorModal.style.display = 'flex';
        document.body.style.overflow = 'hidden';
        setTimeout(function () { document.getElementById('metaName').focus(); }, 100);
    }

    function closeMetaEditor() {
        metaEditorModal.style.display = 'none';
        document.body.style.overflow = '';
        editingMetaId = null;
        editingMetaType = null;
    }

    function confirmAction(title, text, warn, onConfirm) {
        document.getElementById('confirmTitle').textContent = title;
        document.getElementById('confirmText').textContent = text;
        var warnEl = document.getElementById('confirmWarn');
        if (warn) {
            warnEl.textContent = warn;
            warnEl.style.display = 'block';
        } else {
            warnEl.style.display = 'none';
        }
        confirmModal.style.display = 'flex';
        document.body.style.overflow = 'hidden';

        var confirmBtn = document.getElementById('confirmDeleteBtn');
        var cancelBtn = document.getElementById('cancelDeleteBtn');

        function cleanup() {
            confirmModal.style.display = 'none';
            document.body.style.overflow = '';
            confirmBtn.removeEventListener('click', handleConfirm);
            cancelBtn.removeEventListener('click', handleCancel);
        }

        function handleConfirm() {
            cleanup();
            onConfirm();
        }

        function handleCancel() {
            cleanup();
        }

        confirmBtn.addEventListener('click', handleConfirm);
        cancelBtn.addEventListener('click', handleCancel);
    }

    function showLogin() {
        adminApp.style.display = 'none';
        authMask.style.display = 'block';
        loginModal.style.display = 'flex';
        document.body.style.overflow = 'hidden';
        setTimeout(function () { document.getElementById('loginUsername').focus(); }, 100);
    }

    function hideLogin() {
        loginModal.style.display = 'none';
        authMask.style.display = 'none';
        adminApp.style.display = 'block';
        document.body.style.overflow = '';
        loginError.style.display = 'none';
        loginForm.reset();
    }

    function handleLogin(e) {
        e.preventDefault();
        var username = document.getElementById('loginUsername').value.trim();
        var password = document.getElementById('loginPassword').value;
        loginError.style.display = 'none';

        request('/login', {
            method: 'POST',
            body: {
                username: username,
                password: password
            }
        }).then(function () {
            hideLogin();
            refreshAll();
            loadBingWallpaper();
        }).catch(function (err) {
            loginError.textContent = err.message === 'unauthorized' ? '用户名或密码错误' : '登录失败：' + err.message;
            loginError.style.display = 'block';
        });
    }

    function handleLogout() {
        request('/logout', { method: 'POST' }).then(function () {
            showLogin();
        }).catch(function () {
            showLogin();
        });
    }

    function checkAuth() {
        return request('/me').then(function () {
            return true;
        }).catch(function (err) {
            if (err.message === 'unauthorized') {
                return false;
            }
            return false;
        });
    }

    function loadBingWallpaper() {
        var script = document.createElement('script');
        script.src = '../assets/json/images.js?t=' + Date.now();
        script.onload = function () {
            if (window.BING_IMAGES && window.BING_IMAGES.length > 0) {
                var key = 'admin-bg-index';
                var index = parseInt(sessionStorage.getItem(key), 10);
                if (isNaN(index) || index >= window.BING_IMAGES.length) {
                    index = Math.floor(Math.random() * window.BING_IMAGES.length);
                    sessionStorage.setItem(key, index);
                }
                var url = '../assets/' + window.BING_IMAGES[index];
                document.body.style.backgroundImage = "url('" + url.replace(/['\\]/g, '\\$&') + "')";
            }
        };
        document.body.appendChild(script);
    }

    function init() {
        loadBingWallpaper();

        checkAuth().then(function (loggedIn) {
            if (loggedIn) {
                hideLogin();
                refreshAll();
            } else {
                showLogin();
            }
        });

        document.getElementById('newPostBtn').addEventListener('click', function () {
            openEditor(null);
        });

        document.getElementById('closeModalBtn').addEventListener('click', closeEditor);
        document.getElementById('cancelBtn').addEventListener('click', closeEditor);

        document.getElementById('newCategoryBtn').addEventListener('click', function () {
            openMetaEditor('category', null, null);
        });
        document.getElementById('newTagBtn').addEventListener('click', function () {
            openMetaEditor('tag', null, null);
        });
        document.getElementById('closeMetaModalBtn').addEventListener('click', closeMetaEditor);
        document.getElementById('cancelMetaBtn').addEventListener('click', closeMetaEditor);

        document.getElementById('logoutBtn').addEventListener('click', handleLogout);

        document.querySelectorAll('.admin-tab').forEach(function (btn) {
            btn.addEventListener('click', function () {
                switchTab(btn.dataset.tab);
            });
        });

        var previewBtn = document.getElementById('previewToggleBtn');
        var previewPane = document.getElementById('previewPane');
        var previewContent = document.getElementById('previewContent');
        var contentTextarea = document.getElementById('postContent');
        var previewOn = false;
        if (previewBtn) {
            previewBtn.addEventListener('click', function () {
                previewOn = !previewOn;
                if (previewOn) {
                    previewPane.style.display = 'block';
                    previewBtn.textContent = '关闭预览';
                    previewBtn.classList.add('btn-primary');
                    previewBtn.classList.remove('btn-ghost');
                    updatePreview();
                } else {
                    previewPane.style.display = 'none';
                    previewBtn.textContent = '预览';
                    previewBtn.classList.remove('btn-primary');
                    previewBtn.classList.add('btn-ghost');
                }
            });
        }
        if (contentTextarea) {
            contentTextarea.addEventListener('input', function () {
                if (previewOn) updatePreview();
            });
        }
        function updatePreview() {
            var md = contentTextarea.value || '';
            if (window.marked) {
                previewContent.innerHTML = marked.parse(md);
            } else {
                previewContent.innerHTML = '<pre>' + escapeHtml(md) + '</pre>';
            }
        }

        document.querySelectorAll('.status-filter__btn').forEach(function (btn) {
            btn.addEventListener('click', function () {
                document.querySelectorAll('.status-filter__btn').forEach(function (b) {
                    b.classList.remove('active');
                });
                btn.classList.add('active');
                statusFilter = btn.dataset.filter;
                fetchAllPosts().then(function (posts) {
                    renderPosts(posts);
                });
            });
        });

        loginForm.addEventListener('submit', handleLogin);

        document.getElementById('postForm').addEventListener('submit', function (e) {
            e.preventDefault();
            var checkedTagIds = [];
            document.querySelectorAll('#postTagCheckboxes input:checked').forEach(function (cb) {
                checkedTagIds.push(parseInt(cb.value, 10));
            });
            var checkedCategoryId = 0;
            var categoryRadio = document.querySelector('#postCategoryCheckboxes input:checked');
            if (categoryRadio) checkedCategoryId = parseInt(categoryRadio.value, 10);

            var title = document.getElementById('postTitle').value.trim();
            var summary = document.getElementById('postSummary').value.trim();
            var content = document.getElementById('postContent').value;
            var coverImage = document.getElementById('postCoverImage').value.trim();
            var status = document.querySelector('input[name="postStatus"]:checked').value;

            if (!title) {
                showToast('请先填写文章标题。', 'warn');
                return;
            }

            var slug = generateSlug(title);

            var body = {
                title: title,
                slug: slug,
                summary: summary,
                content_md: content,
                content_html: content,
                cover_image: coverImage,
                category_id: checkedCategoryId || 0,
                tag_ids: checkedTagIds,
                status: status
            };

            var requestPromise;
            if (editingPostId) {
                requestPromise = request('/articles/' + editingPostId, {
                    method: 'PUT',
                    body: body
                });
            } else {
                requestPromise = request('/articles', {
                    method: 'POST',
                    body: body
                });
            }

            requestPromise.then(function () {
                closeEditor();
                refreshAll();
            }).catch(function (err) {
                reportError('保存失败', err)
            });
        });

        document.getElementById('metaForm').addEventListener('submit', function (e) {
            e.preventDefault();
            var type = document.getElementById('metaType').value;
            var newName = document.getElementById('metaName').value.trim();
            if (!newName) return;

            var slug = generateSlug(newName);
            var body = {
                name: newName,
                slug: slug,
                description: ''
            };

            var requestPromise;
            if (editingMetaId) {
                requestPromise = request(metaEndpoint(type) + '/' + editingMetaId, {
                    method: 'PUT',
                    body: body
                });
            } else {
                requestPromise = request(metaEndpoint(type), {
                    method: 'POST',
                    body: body
                });
            }

            requestPromise.then(function () {
                closeMetaEditor();
                refreshAll();
            }).catch(function (err) {
                reportError('保存失败', err)
            });
        });

        postList.addEventListener('click', function (e) {
            var btn = e.target.closest('[data-action]');
            if (!btn) return;
            var action = btn.dataset.action;
            var id = parseInt(btn.dataset.id, 10);

            if (action === 'edit') {
                openEditor(id);
            } else if (action === 'delete') {
                var title = btn.dataset.title;
                confirmAction(
                    '删除文章',
                    '确定要删除「' + (title || '这篇文章') + '」吗？',
                    '此操作不可恢复',
                    function () {
                        request('/articles/' + id, { method: 'DELETE' }).then(function () {
                            refreshAll();
                        }).catch(function (err) {
                            reportError('删除失败', err)
                        });
                    }
                );
            }
        });

        function handleMetaClick(e) {
            var btn = e.target.closest('[data-action]');
            if (!btn) return;
            var action = btn.dataset.action;
            var type = btn.dataset.type;
            var id = parseInt(btn.dataset.id, 10);
            var name = btn.dataset.name;

            if (action === 'edit-meta') {
                openMetaEditor(type, id, name);
            } else if (action === 'delete-meta') {
                var count = parseInt(btn.dataset.count || '0', 10);
                var typeName = type === 'category' ? '分类' : '标签';
                var warn = count > 0 ? '该' + typeName + '下有 ' + count + ' 篇文章，删除后可能影响文章显示' : '';
                confirmAction(
                    '删除' + typeName,
                    '确定要删除' + typeName + '「' + name + '」吗？',
                    warn,
                    function () {
                        request(metaEndpoint(type) + '/' + id, { method: 'DELETE' }).then(function () {
                            refreshAll();
                        }).catch(function (err) {
                            reportError('删除失败', err)
                        });
                    }
                );
            }
        }

        categoryList.addEventListener('click', handleMetaClick);
        tagListEl.addEventListener('click', handleMetaClick);

        document.getElementById('postCategoryCheckboxes').addEventListener('change', function (e) {
            var radio = e.target;
            if (radio.type === 'radio') {
                document.querySelectorAll('#postCategoryCheckboxes .tag-checkbox').forEach(function (label) {
                    label.classList.remove('checked');
                });
                radio.closest('.tag-checkbox').classList.add('checked');
            }
        });

        document.getElementById('postTagCheckboxes').addEventListener('change', function (e) {
            var cb = e.target;
            if (cb.type === 'checkbox') {
                cb.closest('.tag-checkbox').classList.toggle('checked', cb.checked);
            }
        });

        document.querySelectorAll('.modal__overlay').forEach(function (overlay) {
            overlay.addEventListener('click', function () {
                if (editorModal.style.display === 'flex') closeEditor();
                if (metaEditorModal.style.display === 'flex') closeMetaEditor();
            });
        });

        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape') {
                if (loginModal.style.display === 'flex') return;
                if (editorModal.style.display === 'flex') closeEditor();
                if (metaEditorModal.style.display === 'flex') closeMetaEditor();
                if (confirmModal.style.display === 'flex') {
                    confirmModal.style.display = 'none';
                    document.body.style.overflow = '';
                }
            }
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
