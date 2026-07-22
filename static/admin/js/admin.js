(function () {
    'use strict';

    var toastTimer = null;
    var OPENVERSE_IMAGE_ENDPOINT = 'https://api.openverse.org/v1/images/';

    function adminApiURL(path) {
        if (path === '/login' || path === '/logout' || path === '/me') {
            return '/api/auth' + path;
        }
        return '/api/content/admin' + path;
    }

    function request(path, options) {
        options = options || {};
        return fetch(adminApiURL(path), {
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
                throw new Error(data.msg || data.message || '服务暂时没有返回可用结果，请稍后再试');
            }
            return data.data;
        });
    }

    var ADMIN_COVER_CHUNK_SIZE = 1024 * 1024;
    var ADMIN_UPLOAD_USER_ID = 1;

    function uploadCoverFile(file, onProgress) {
        if (!file) return Promise.resolve('');
        if (!file.type || file.type.indexOf('image/') !== 0) {
            return Promise.reject(new Error('请选择图片文件作为封面'));
        }
        if (typeof onProgress === 'function') onProgress(0);
        return createFileSha256(file)
            .then(function (sha256) {
                return initChunkUpload(file, sha256);
            })
            .then(function (task) {
                return uploadFileChunks(file, task, onProgress)
                    .then(function () {
                        return completeChunkUpload(task.uploadID);
                    })
                    .catch(function (err) {
                        return cancelChunkUpload(task.uploadID).then(function () {
                            throw err;
                        }, function () {
                            throw err;
                        });
                    });
            })
            .then(function (record) {
                if (!record || !record.url) {
                    throw new Error('封面上传成功，但服务没有返回可访问地址');
                }
                if (typeof onProgress === 'function') onProgress(100);
                return record.url;
            });
    }

    function initChunkUpload(file, sha256) {
        var formData = new FormData();
        formData.append('user_id', String(ADMIN_UPLOAD_USER_ID));
        formData.append('file_name', file.name || 'cover-image');
        formData.append('file_size', String(file.size));
        formData.append('chunk_size', String(ADMIN_COVER_CHUNK_SIZE));
        formData.append('content_type', file.type || 'application/octet-stream');
        formData.append('biz_type', 'article_cover');
        formData.append('sha256', sha256);
        return postMediaForm('/api/media/upload/tasks/init', formData).then(function (data) {
            var uploadID = data.upload_id || data.UploadID;
            var chunkSize = Number(data.chunk_size || data.ChunkSize || ADMIN_COVER_CHUNK_SIZE);
            var chunkCount = Number(data.chunk_count || data.ChunkCount || Math.ceil(file.size / chunkSize));
            if (!uploadID) throw new Error('分片上传初始化失败：服务没有返回上传任务');
            return {
                uploadID: uploadID,
                chunkSize: chunkSize > 0 ? chunkSize : ADMIN_COVER_CHUNK_SIZE,
                chunkCount: chunkCount > 0 ? chunkCount : Math.max(1, Math.ceil(file.size / ADMIN_COVER_CHUNK_SIZE))
            };
        });
    }

    function uploadFileChunks(file, task, onProgress) {
        var uploaded = 0;
        var promise = Promise.resolve();
        for (var index = 0; index < task.chunkCount; index += 1) {
            promise = promise.then((function (chunkIndex) {
                return function () {
                    var start = chunkIndex * task.chunkSize;
                    var end = Math.min(file.size, start + task.chunkSize);
                    var chunk = file.slice(start, end);
                    return sendChunkWithProgress(task.uploadID, chunkIndex, chunk, function (loaded) {
                        if (typeof onProgress !== 'function') return;
                        var percent = file.size > 0 ? Math.round(((uploaded + loaded) / file.size) * 98) : 98;
                        onProgress(Math.min(98, Math.max(0, percent)));
                    }).then(function () {
                        uploaded += chunk.size;
                        if (typeof onProgress === 'function') {
                            var percent = file.size > 0 ? Math.round((uploaded / file.size) * 98) : 98;
                            onProgress(Math.min(98, percent));
                        }
                    });
                };
            })(index));
        }
        return promise;
    }

    function sendChunkWithProgress(uploadID, chunkIndex, chunk, onChunkProgress) {
        return new Promise(function (resolve, reject) {
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/api/media/upload/tasks/' + encodeURIComponent(uploadID) + '/chunks/' + chunkIndex + '?user_id=' + ADMIN_UPLOAD_USER_ID, true);
            xhr.withCredentials = true;
            xhr.setRequestHeader('Content-Type', 'application/octet-stream');
            xhr.upload.addEventListener('progress', function (evt) {
                if (!evt.lengthComputable || typeof onChunkProgress !== 'function') return;
                onChunkProgress(evt.loaded);
            });
            xhr.addEventListener('load', function () {
                var data = parseMediaResponse(xhr, '/api/media/upload/tasks/' + uploadID + '/chunks/' + chunkIndex);
                if (data.error) {
                    reject(data.error);
                    return;
                }
                resolve(data.payload);
            });
            xhr.addEventListener('error', function () {
                reject(new Error('封面分片上传失败，请检查网络后重试'));
            });
            xhr.send(chunk);
        });
    }

    function completeChunkUpload(uploadID) {
        var formData = new FormData();
        formData.append('user_id', String(ADMIN_UPLOAD_USER_ID));
        return postMediaForm('/api/media/upload/tasks/' + encodeURIComponent(uploadID) + '/complete', formData);
    }

    function cancelChunkUpload(uploadID) {
        if (!uploadID) return Promise.resolve();
        var formData = new FormData();
        formData.append('user_id', String(ADMIN_UPLOAD_USER_ID));
        return postMediaForm('/api/media/upload/tasks/' + encodeURIComponent(uploadID) + '/cancel', formData).catch(function () {});
    }

    function postMediaForm(path, formData) {
        return fetch(path, {
            method: 'POST',
            body: formData,
            credentials: 'same-origin'
        }).then(function (res) {
            return res.text().then(function (text) {
                var payload = {};
                try {
                    payload = text ? JSON.parse(text) : {};
                } catch (err) {
                    throw new Error('封面上传失败，请稍后再试');
                }
                if (!res.ok) {
                    throw new Error(friendlyHttpError(path, res.status));
                }
                if (payload.code !== 0) {
                    throw new Error(payload.msg || payload.message || '封面上传失败，请稍后再试');
                }
                return payload.data || {};
            });
        });
    }

    function parseMediaResponse(xhr, path) {
        if (xhr.status < 200 || xhr.status >= 300) {
            return { error: new Error(friendlyHttpError(path, xhr.status)) };
        }
        try {
            var payload = JSON.parse(xhr.responseText || '{}');
            if (payload.code !== 0) {
                return { error: new Error(payload.msg || payload.message || '封面上传失败，请稍后再试') };
            }
            return { payload: payload.data || {} };
        } catch (err) {
            return { error: new Error('封面上传失败，请稍后再试') };
        }
    }


    function createFileSha256(file) {
        if (!window.crypto || !crypto.subtle || typeof crypto.subtle.digest !== 'function') {
            return Promise.reject(new Error('当前浏览器不支持文件完整性校验，请换用新版浏览器后重试'));
        }
        return file.arrayBuffer().then(function (buffer) {
            return crypto.subtle.digest('SHA-256', buffer);
        }).then(function (digest) {
            var bytes = Array.prototype.slice.call(new Uint8Array(digest));
            return bytes.map(function (byte) {
                return byte.toString(16).padStart(2, '0');
            }).join('');
        });
    }
    function friendlyHttpError(path, status) {
        if (status === 404) {
            if (path.indexOf('/articles/') === 0) return '没有找到这篇文章，可能已被删除或接口路径尚未同步。';
            return '没有找到对应的数据，请刷新页面后再试。';
        }
        if (status === 401) return '登录状态已过期，请重新登录。';
        if (status === 403) return '当前账号没有权限执行这个操作。';
        if (status === 413) return '文件太大，请换一张更小的图片或使用分片上传。';
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
        showToast(prefix + '：' + ((err && err.message) || '请稍后再试。'), 'error');
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
            status.textContent = message || '选择图片后会自动分片上传，并在保存文章时写入封面地址。';
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
            setCoverStatus('选择图片后会自动分片上传，并在保存文章时写入封面地址。');
            setCoverProgress(0);
            return;
        }
        preview.style.display = 'block';
        preview.innerHTML = '<img src="' + escapeHtml(url) + '" alt="封面预览">' +
            '<div class="cover-preview__url">' + escapeHtml(url) + '</div>';
        setCoverStatus('封面已上传，保存文章后生效。');
    }

    function setCoverLibraryStatus(message) {
        var status = document.getElementById('coverLibraryStatus');
        if (status) status.textContent = message || '选择一张图片后会写入封面地址，保存文章后生效。';
    }
    function openCoverLibrary() {
        var modal = document.getElementById('coverLibraryModal');
        if (!modal) return;
        modal.style.display = 'flex';
        setCoverLibraryStatus('正在加载推荐素材...');
        var results = document.getElementById('coverLibraryResults');
        if (results && !results.innerHTML) {
            searchCoverLibrary();
        }
    }

    function closeCoverLibrary() {
        var modal = document.getElementById('coverLibraryModal');
        if (modal) modal.style.display = 'none';
    }

    function searchCoverLibrary() {
        var queryInput = document.getElementById('coverLibraryQuery');
        var results = document.getElementById('coverLibraryResults');
        var query = queryInput && queryInput.value.trim() ? queryInput.value.trim() : 'technology abstract';
        if (!results) return Promise.resolve();

        results.innerHTML = '<div class="cover-library__empty">正在从素材库加载图片...</div>';
        setCoverLibraryStatus('正在搜索：' + query);

        var params = new URLSearchParams({
            q: query,
            page: '1',
            page_size: '18'
        });

        return fetch(OPENVERSE_IMAGE_ENDPOINT + '?' + params.toString())
            .then(function (res) {
                if (!res.ok) throw new Error('素材库暂时没有响应');
                return res.json();
            })
            .then(function (data) {
                renderCoverLibraryResults(data.results || []);
                setCoverLibraryStatus('点击图片下方按钮即可使用；请按素材来源保留必要署名。');
            })
            .catch(function () {
                renderFallbackCoverLibrary(query);
                setCoverLibraryStatus('Openverse 响应较慢，已切换到备用封面素材。');
            });
    }

    function renderFallbackCoverLibrary(query) {
        var normalized = (query || 'technology abstract').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '') || 'cover';
        var fallbackItems = [];
        for (var i = 1; i <= 18; i++) {
            var seed = normalized + '-' + i;
            fallbackItems.push({
                title: '备用封面素材 ' + i,
                creator: 'Lorem Picsum',
                license: 'placeholder photo',
                thumbnail: 'https://picsum.photos/seed/' + encodeURIComponent(seed) + '/360/160',
                url: 'https://picsum.photos/seed/' + encodeURIComponent(seed) + '/1200/420'
            });
        }
        renderCoverLibraryResults(fallbackItems);
    }

    function renderCoverLibraryResults(items) {
        var results = document.getElementById('coverLibraryResults');
        if (!results) return;
        if (!items.length) {
            results.innerHTML = '<div class="cover-library__empty">没有找到合适图片，试试 cloud、code、mountain、abstract。</div>';
            return;
        }
        results.innerHTML = items.map(function (item) {
            var imageUrl = item.url || item.thumbnail || '';
            var thumbUrl = item.thumbnail || item.url || '';
            var title = item.title || 'Untitled image';
            var creator = item.creator || item.provider || 'Openverse';
            var license = item.license || 'open license';
            if (!imageUrl || !thumbUrl) return '';
            return '<article class="cover-library__card">' +
                '<img class="cover-library__thumb" src="' + escapeHtml(thumbUrl) + '" alt="' + escapeHtml(title) + '" loading="lazy">' +
                '<div class="cover-library__meta">' +
                    '<p class="cover-library__title">' + escapeHtml(title) + '</p>' +
                    '<p class="cover-library__credit">' + escapeHtml(creator) + ' · ' + escapeHtml(license) + '</p>' +
                    '<button type="button" class="btn btn-primary btn-sm cover-library__use" data-cover-url="' + escapeHtml(imageUrl) + '">使用这张</button>' +
                '</div>' +
            '</article>';
        }).join('');
    }

    function initCoverControls() {
        var coverFileInput = document.getElementById('postCoverFile');
        var clearCoverBtn = document.getElementById('clearCoverBtn');
        var openLibraryBtn = document.getElementById('openCoverLibraryBtn');
        var closeLibraryBtn = document.getElementById('closeCoverLibraryBtn');
        var searchLibraryBtn = document.getElementById('searchCoverLibraryBtn');
        var libraryQuery = document.getElementById('coverLibraryQuery');
        var libraryResults = document.getElementById('coverLibraryResults');
        var libraryOverlay = document.querySelector('[data-cover-library-close]');

        if (coverFileInput) {
            coverFileInput.addEventListener('change', function (e) {
                var file = e.target.files && e.target.files[0];
                if (!file) return;
                setCoverProgress(1, '开始上传...');
                setCoverStatus('封面分片上传中...');
                uploadCoverFile(file, function (percent) {
                    setCoverProgress(percent, '上传中 ' + percent + '%');
                }).then(function (url) {
                    setCoverImage(url);
                    setCoverProgress(100, '上传完成');
                    showToast('封面上传成功', 'success');
                    setTimeout(function () { setCoverProgress(0); }, 900);
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
        if (openLibraryBtn) openLibraryBtn.addEventListener('click', openCoverLibrary);
        if (closeLibraryBtn) closeLibraryBtn.addEventListener('click', closeCoverLibrary);
        if (libraryOverlay) libraryOverlay.addEventListener('click', closeCoverLibrary);
        if (searchLibraryBtn) searchLibraryBtn.addEventListener('click', searchCoverLibrary);
        if (libraryQuery) {
            libraryQuery.addEventListener('keydown', function (event) {
                if (event.key === 'Enter') {
                    event.preventDefault();
                    searchCoverLibrary();
                }
            });
        }
        if (libraryResults) {
            libraryResults.addEventListener('click', function (event) {
                var button = event.target.closest('[data-cover-url]');
                if (!button) return;
                setCoverImage(button.getAttribute('data-cover-url'));
                setCoverStatus('已从素材库选择封面，保存文章后生效。');
                closeCoverLibrary();
            });
        }
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
        warnEl.textContent = warn || '';
        warnEl.style.display = warn ? 'block' : 'none';
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
        function handleConfirm() { cleanup(); onConfirm(); }
        function handleCancel() { cleanup(); }
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
        request('/login', { method: 'POST', body: { username: username, password: password } }).then(function () {
            hideLogin();
            refreshAll();
            loadBingWallpaper();
        }).catch(function (err) {
            loginError.textContent = err.message === 'unauthorized' ? '用户名或密码错误' : '登录失败：' + err.message;
            loginError.style.display = 'block';
        });
    }

    function handleLogout() {
        request('/logout', { method: 'POST' }).then(showLogin).catch(showLogin);
    }

    function checkAuth() {
        return request('/me').then(function () { return true; }).catch(function () { return false; });
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
            if (loggedIn) { hideLogin(); refreshAll(); } else { showLogin(); }
        });
        document.getElementById('newPostBtn').addEventListener('click', function () { openEditor(null); });
        document.getElementById('closeModalBtn').addEventListener('click', closeEditor);
        document.getElementById('cancelBtn').addEventListener('click', closeEditor);
        initCoverControls();
        document.getElementById('newCategoryBtn').addEventListener('click', function () { openMetaEditor('category', null, null); });
        document.getElementById('newTagBtn').addEventListener('click', function () { openMetaEditor('tag', null, null); });
        document.getElementById('closeMetaModalBtn').addEventListener('click', closeMetaEditor);
        document.getElementById('cancelMetaBtn').addEventListener('click', closeMetaEditor);
        document.getElementById('logoutBtn').addEventListener('click', handleLogout);
        document.querySelectorAll('.admin-tab').forEach(function (btn) { btn.addEventListener('click', function () { switchTab(btn.dataset.tab); }); });

        var previewBtn = document.getElementById('previewToggleBtn');
        var previewPane = document.getElementById('previewPane');
        var previewContent = document.getElementById('previewContent');
        var contentTextarea = document.getElementById('postContent');
        var previewOn = false;
        function updatePreview() {
            var md = contentTextarea.value || '';
            previewContent.innerHTML = window.marked ? marked.parse(md) : '<pre>' + escapeHtml(md) + '</pre>';
        }
        if (previewBtn) {
            previewBtn.addEventListener('click', function () {
                previewOn = !previewOn;
                previewPane.style.display = previewOn ? 'block' : 'none';
                previewBtn.textContent = previewOn ? '关闭预览' : '预览';
                previewBtn.classList.toggle('btn-primary', previewOn);
                previewBtn.classList.toggle('btn-ghost', !previewOn);
                if (previewOn) updatePreview();
            });
        }
        if (contentTextarea) contentTextarea.addEventListener('input', function () { if (previewOn) updatePreview(); });

        document.querySelectorAll('.status-filter__btn').forEach(function (btn) {
            btn.addEventListener('click', function () {
                document.querySelectorAll('.status-filter__btn').forEach(function (b) { b.classList.remove('active'); });
                btn.classList.add('active');
                statusFilter = btn.dataset.filter;
                fetchAllPosts().then(renderPosts);
            });
        });
        loginForm.addEventListener('submit', handleLogin);
        document.getElementById('postForm').addEventListener('submit', handlePostSubmit);
        document.getElementById('metaForm').addEventListener('submit', handleMetaSubmit);
        postList.addEventListener('click', handlePostListClick);
        categoryList.addEventListener('click', handleMetaListClick);
        tagListEl.addEventListener('click', handleMetaListClick);
        document.getElementById('postCategoryCheckboxes').addEventListener('change', handleCategoryChoice);
        document.getElementById('postTagCheckboxes').addEventListener('change', handleTagChoice);
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

    function handlePostSubmit(e) {
        e.preventDefault();
        var checkedTagIds = [];
        document.querySelectorAll('#postTagCheckboxes input:checked').forEach(function (cb) { checkedTagIds.push(parseInt(cb.value, 10)); });
        var categoryRadio = document.querySelector('#postCategoryCheckboxes input:checked');
        var title = document.getElementById('postTitle').value.trim();
        if (!title) { showToast('请先填写文章标题。', 'warn'); return; }
        var body = {
            title: title,
            slug: generateSlug(title),
            summary: document.getElementById('postSummary').value.trim(),
            content_md: document.getElementById('postContent').value,
            content_html: document.getElementById('postContent').value,
            cover_image: document.getElementById('postCoverImage').value.trim(),
            category_id: categoryRadio ? parseInt(categoryRadio.value, 10) : 0,
            tag_ids: checkedTagIds,
            status: document.querySelector('input[name="postStatus"]:checked').value
        };
        var requestPromise = editingPostId ? request('/articles/' + editingPostId, { method: 'PUT', body: body }) : request('/articles', { method: 'POST', body: body });
        requestPromise.then(function () { closeEditor(); refreshAll(); }).catch(function (err) { reportError('保存失败', err); });
    }

    function handleMetaSubmit(e) {
        e.preventDefault();
        var type = document.getElementById('metaType').value;
        var newName = document.getElementById('metaName').value.trim();
        if (!newName) return;
        var body = { name: newName, slug: generateSlug(newName), description: '' };
        var requestPromise = editingMetaId ? request(metaEndpoint(type) + '/' + editingMetaId, { method: 'PUT', body: body }) : request(metaEndpoint(type), { method: 'POST', body: body });
        requestPromise.then(function () { closeMetaEditor(); refreshAll(); }).catch(function (err) { reportError('保存失败', err); });
    }

    function handlePostListClick(e) {
        var btn = e.target.closest('[data-action]');
        if (!btn) return;
        var id = parseInt(btn.dataset.id, 10);
        if (btn.dataset.action === 'edit') {
            openEditor(id);
        } else if (btn.dataset.action === 'delete') {
            var title = btn.dataset.title;
            confirmAction('删除文章', '确定要删除「' + (title || '这篇文章') + '」吗？', '此操作不可恢复', function () {
                request('/articles/' + id, { method: 'DELETE' }).then(refreshAll).catch(function (err) { reportError('删除失败', err); });
            });
        }
    }

    function handleMetaListClick(e) {
        var btn = e.target.closest('[data-action]');
        if (!btn) return;
        var type = btn.dataset.type;
        var id = parseInt(btn.dataset.id, 10);
        var name = btn.dataset.name;
        if (btn.dataset.action === 'edit-meta') {
            openMetaEditor(type, id, name);
        } else if (btn.dataset.action === 'delete-meta') {
            var count = parseInt(btn.dataset.count || '0', 10);
            var typeName = type === 'category' ? '分类' : '标签';
            var warn = count > 0 ? '该' + typeName + '下有 ' + count + ' 篇文章，删除后可能影响文章显示' : '';
            confirmAction('删除' + typeName, '确定要删除' + typeName + '「' + name + '」吗？', warn, function () {
                request(metaEndpoint(type) + '/' + id, { method: 'DELETE' }).then(refreshAll).catch(function (err) { reportError('删除失败', err); });
            });
        }
    }

    function handleCategoryChoice(e) {
        var radio = e.target;
        if (radio.type === 'radio') {
            document.querySelectorAll('#postCategoryCheckboxes .tag-checkbox').forEach(function (label) { label.classList.remove('checked'); });
            radio.closest('.tag-checkbox').classList.add('checked');
        }
    }

    function handleTagChoice(e) {
        var cb = e.target;
        if (cb.type === 'checkbox') cb.closest('.tag-checkbox').classList.toggle('checked', cb.checked);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
