(function () {
  const app = window.AdminApp;
  if (!app) return;

  const page = document.querySelector('.article-editor-page');
  if (!page) return;

  const articleID = page.getAttribute('data-article-id');
  const titleInput = document.getElementById('article-title');
  const slugInput = document.getElementById('article-slug');
  const summaryInput = document.getElementById('article-summary');
  const contentInput = document.getElementById('article-content');
  const statusSelect = document.getElementById('article-status');
  const categorySelect = document.getElementById('article-category');
  const tagsNode = document.getElementById('article-tags');
  const coverInput = document.getElementById('article-cover');
  const coverFile = document.getElementById('article-cover-file');
  const coverFileName = document.getElementById('article-cover-file-name');
  const coverUploadBtn = document.getElementById('article-cover-upload');
  const coverPreviewWrap = document.getElementById('article-cover-preview');
  const coverPreviewImage = document.getElementById('article-cover-preview-image');
  const saveDraftBtn = document.getElementById('article-save-draft');
  const publishBtn = document.getElementById('article-publish');

  const MAX_COVER_SIZE = 5 * 1024 * 1024;
  const ALLOWED_COVER_TYPES = new Set(['image/jpeg', 'image/png', 'image/webp', 'image/gif']);

  let categories = [];
  let tags = [];

  function renderCategories() {
    const options = ['<option value="0">No category</option>'];
    categories.forEach(function (item) {
      options.push(`<option value="${item.id}">${app.escapeHTML(item.name)}</option>`);
    });
    categorySelect.innerHTML = options.join('');
  }

  function renderTags(selectedIDs) {
    const selected = new Set((selectedIDs || []).map(Number));
    tagsNode.innerHTML = tags.map(function (item) {
      const checked = selected.has(Number(item.id)) ? 'checked' : '';
      return `<label class="checkbox-chip"><input type="checkbox" value="${item.id}" ${checked}> <span>${app.escapeHTML(item.name)}</span></label>`;
    }).join('');
  }

  function updateCoverPreview(url) {
    if (!url) {
      coverPreviewWrap.hidden = true;
      return;
    }
    coverPreviewImage.src = url;
    coverPreviewWrap.hidden = false;
  }

  function setCoverUploadState(uploading) {
    if (!coverUploadBtn) return;
    coverUploadBtn.disabled = uploading;
    coverUploadBtn.textContent = uploading ? 'Uploading...' : 'Upload Cover';
  }

  function setCoverFileName(file) {
    if (!coverFileName) return;
    coverFileName.textContent = file ? file.name : '未选择文件';
  }

  function validateCoverFile(file) {
    if (!file) {
      return 'Please choose a file first.';
    }
    if (file.size <= 0) {
      return 'The selected file is empty.';
    }
    if (file.size > MAX_COVER_SIZE) {
      return 'Image size must be less than or equal to 5MB.';
    }
    if (!ALLOWED_COVER_TYPES.has(file.type)) {
      return 'Only JPG, JPEG, PNG, WEBP and GIF images are allowed.';
    }
    return '';
  }

  async function loadInitialData() {
    const user = await app.requireAuth();
    if (!user) return null;
    app.mountUser(user);
    app.bindLogout();

    const [categoryResp, tagResp, articleResp] = await Promise.all([
      app.request('/api/admin/categories', { headers: { Accept: 'application/json' } }),
      app.request('/api/admin/tags', { headers: { Accept: 'application/json' } }),
      articleID ? app.request('/api/admin/articles', { headers: { Accept: 'application/json' } }) : Promise.resolve(null)
    ]);

    categories = categoryResp.list || [];
    tags = tagResp.list || [];
    renderCategories();
    renderTags([]);

    if (!articleID || !articleResp) return null;
    const match = (articleResp.list || []).find(function (item) { return String(item.id) === String(articleID); });
    if (!match) throw new Error('Article not found');
    return match;
  }

  function fillArticle(article) {
    titleInput.value = article.title || '';
    slugInput.value = article.slug || '';
    summaryInput.value = article.summary || '';
    contentInput.value = article.content_md || '';
    statusSelect.value = article.status || 'draft';
    categorySelect.value = String(article.category_id || 0);
    renderTags(article.tag_ids || []);
    coverInput.value = article.cover_image || '';
    updateCoverPreview(article.cover_image || '');
  }

  async function uploadCover() {
    const file = coverFile.files && coverFile.files[0];
    const validationMessage = validateCoverFile(file);
    if (validationMessage) {
      app.setFeedback('editor-feedback', validationMessage, false);
      app.showToast(validationMessage, 'error');
      return;
    }

    const formData = new FormData();
    formData.append('biz_type', 'article_cover');
    formData.append('file', file);

    try {
      await app.requireAuth();
      setCoverUploadState(true);
      const data = await app.request('/api/admin/upload', {
        method: 'POST',
        body: formData
      });
      coverInput.value = data.upload.file_url;
      updateCoverPreview(data.upload.file_url);
      coverFile.value = '';
      setCoverFileName(null);
      app.setFeedback('editor-feedback', 'Cover uploaded successfully.', true);
      app.showToast('Cover uploaded', 'success');
    } catch (error) {
      app.setFeedback('editor-feedback', error.message, false);
      app.showToast(error.message, 'error');
    } finally {
      setCoverUploadState(false);
    }
  }

  function payloadWithStatus(statusValue) {
    return {
      title: titleInput.value.trim(),
      slug: slugInput.value.trim(),
      summary: summaryInput.value.trim(),
      content_md: contentInput.value,
      cover_image: coverInput.value.trim(),
      category_id: Number(categorySelect.value || 0),
      tag_ids: app.collectCheckedValues('#article-tags'),
      status: statusValue
    };
  }

  async function saveArticle(statusValue) {
    const payload = payloadWithStatus(statusValue);
    if (!payload.title || !payload.slug) {
      app.setFeedback('editor-feedback', 'Title and slug are required.', false);
      return;
    }
    try {
      await app.requireAuth();
      if (articleID) {
        await app.jsonRequest('/api/admin/articles/' + articleID, 'PUT', payload);
      } else {
        await app.jsonRequest('/api/admin/articles', 'POST', payload);
      }
      app.setFeedback('editor-feedback', 'Article saved successfully.', true);
      app.showToast('Article saved', 'success');
      if (!articleID) {
        window.setTimeout(function () {
          window.location.href = '/admin/articles';
        }, 800);
      }
    } catch (error) {
      app.setFeedback('editor-feedback', error.message, false);
      app.showToast(error.message, 'error');
    }
  }

  coverUploadBtn.addEventListener('click', uploadCover);
  coverFile.addEventListener('change', function () {
    const file = coverFile.files && coverFile.files[0];
    setCoverFileName(file || null);
  });
  coverInput.addEventListener('input', function () {
    updateCoverPreview(coverInput.value.trim());
  });
  saveDraftBtn.addEventListener('click', function () { saveArticle('draft'); });
  publishBtn.addEventListener('click', function () { saveArticle('published'); });

  loadInitialData().then(function (article) {
    setCoverFileName(null);
    if (article) fillArticle(article);
  }).catch(function (error) {
    app.setFeedback('editor-feedback', error.message, false);
  });
})();