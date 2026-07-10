(function () {
  const app = window.AdminApp;
  if (!app) return;

  const state = {
    categoryEditingId: null,
    tagEditingId: null,
  };

  function taxonomyRow(item, type) {
    return `
      <article class="taxonomy-row">
        <div class="taxonomy-row__main">
          <strong>${app.escapeHTML(item.name)}</strong>
          <div class="taxonomy-row__slug">${app.escapeHTML(item.slug)}</div>
          <p class="taxonomy-row__desc">${app.escapeHTML(item.description || 'No description')}</p>
        </div>
        <div class="taxonomy-row__actions">
          <button class="admin-button admin-button--ghost taxonomy-action" type="button" data-type="${type}" data-action="edit" data-id="${item.id}" data-name="${app.escapeHTML(item.name)}" data-slug="${app.escapeHTML(item.slug)}" data-description="${app.escapeHTML(item.description || '')}">Edit</button>
          <button class="admin-button admin-button--danger taxonomy-action" type="button" data-type="${type}" data-action="delete" data-id="${item.id}" data-name="${app.escapeHTML(item.name)}">Delete</button>
        </div>
      </article>
    `;
  }

  async function loadLists() {
    const [categories, tags] = await Promise.all([
      app.request('/api/admin/categories', { headers: { Accept: 'application/json' } }),
      app.request('/api/admin/tags', { headers: { Accept: 'application/json' } })
    ]);

    document.getElementById('category-list').innerHTML = (categories.list || []).length
      ? categories.list.map(function (item) { return taxonomyRow(item, 'category'); }).join('')
      : '<div class="admin-empty">No categories yet.</div>';
    document.getElementById('tag-list').innerHTML = (tags.list || []).length
      ? tags.list.map(function (item) { return taxonomyRow(item, 'tag'); }).join('')
      : '<div class="admin-empty">No tags yet.</div>';
  }

  function resetCategoryForm() {
    state.categoryEditingId = null;
    document.getElementById('category-form').reset();
    document.getElementById('category-submit-btn').textContent = 'Create Category';
    document.getElementById('category-cancel-btn').hidden = true;
  }

  function resetTagForm() {
    state.tagEditingId = null;
    document.getElementById('tag-form').reset();
    document.getElementById('tag-submit-btn').textContent = 'Create Tag';
    document.getElementById('tag-cancel-btn').hidden = true;
  }

  function fillCategoryForm(button) {
    state.categoryEditingId = Number(button.dataset.id);
    document.getElementById('category-name').value = button.dataset.name || '';
    document.getElementById('category-slug').value = button.dataset.slug || '';
    document.getElementById('category-description').value = button.dataset.description || '';
    document.getElementById('category-submit-btn').textContent = 'Update Category';
    document.getElementById('category-cancel-btn').hidden = false;
    app.setFeedback('category-feedback', 'Editing category #' + state.categoryEditingId, true);
  }

  function fillTagForm(button) {
    state.tagEditingId = Number(button.dataset.id);
    document.getElementById('tag-name').value = button.dataset.name || '';
    document.getElementById('tag-slug').value = button.dataset.slug || '';
    document.getElementById('tag-description').value = button.dataset.description || '';
    document.getElementById('tag-submit-btn').textContent = 'Update Tag';
    document.getElementById('tag-cancel-btn').hidden = false;
    app.setFeedback('tag-feedback', 'Editing tag #' + state.tagEditingId, true);
  }

  async function submitCategoryForm(event) {
    event.preventDefault();
    const payload = {
      name: document.getElementById('category-name').value.trim(),
      slug: document.getElementById('category-slug').value.trim(),
      description: document.getElementById('category-description').value.trim()
    };

    try {
      let successMessage = '';
      if (state.categoryEditingId) {
        await app.jsonRequest('/api/admin/categories/' + state.categoryEditingId, 'PUT', payload);
        successMessage = 'Category updated.';
      } else {
        await app.jsonRequest('/api/admin/categories', 'POST', payload);
        successMessage = 'Category created.';
      }
      resetCategoryForm();
      await loadLists();
      app.setFeedback('category-feedback', successMessage, true);
    } catch (error) {
      app.setFeedback('category-feedback', error.message, false);
    }
  }

  async function submitTagForm(event) {
    event.preventDefault();
    const payload = {
      name: document.getElementById('tag-name').value.trim(),
      slug: document.getElementById('tag-slug').value.trim(),
      description: document.getElementById('tag-description').value.trim()
    };

    try {
      let successMessage = '';
      if (state.tagEditingId) {
        await app.jsonRequest('/api/admin/tags/' + state.tagEditingId, 'PUT', payload);
        successMessage = 'Tag updated.';
      } else {
        await app.jsonRequest('/api/admin/tags', 'POST', payload);
        successMessage = 'Tag created.';
      }
      resetTagForm();
      await loadLists();
      app.setFeedback('tag-feedback', successMessage, true);
    } catch (error) {
      app.setFeedback('tag-feedback', error.message, false);
    }
  }

  async function handleDelete(type, id, name) {
    const confirmed = window.confirm('Delete ' + type + ' "' + name + '"?');
    if (!confirmed) return;

    try {
      if (type === 'category') {
        await app.jsonRequest('/api/admin/categories/' + id, 'DELETE');
        app.showToast('Category deleted.', 'success');
        if (state.categoryEditingId === id) {
          resetCategoryForm();
          app.setFeedback('category-feedback', '', false);
        }
      } else {
        await app.jsonRequest('/api/admin/tags/' + id, 'DELETE');
        app.showToast('Tag deleted.', 'success');
        if (state.tagEditingId === id) {
          resetTagForm();
          app.setFeedback('tag-feedback', '', false);
        }
      }
      await loadLists();
    } catch (error) {
      app.showToast(error.message, 'error');
    }
  }

  document.getElementById('category-form').addEventListener('submit', function (event) {
    submitCategoryForm(event);
  });

  document.getElementById('tag-form').addEventListener('submit', function (event) {
    submitTagForm(event);
  });

  document.getElementById('category-cancel-btn').addEventListener('click', function () {
    resetCategoryForm();
    app.setFeedback('category-feedback', '', false);
  });
  document.getElementById('tag-cancel-btn').addEventListener('click', function () {
    resetTagForm();
    app.setFeedback('tag-feedback', '', false);
  });

  document.addEventListener('click', function (event) {
    const button = event.target.closest('.taxonomy-action');
    if (!button) return;

    const type = button.dataset.type;
    const action = button.dataset.action;
    const id = Number(button.dataset.id);

    if (action === 'edit') {
      if (type === 'category') {
        fillCategoryForm(button);
      } else {
        fillTagForm(button);
      }
      return;
    }

    if (action === 'delete') {
      handleDelete(type, id, button.dataset.name || 'item');
    }
  });

  async function init() {
    const user = await app.requireAuth();
    if (!user) return;
    app.mountUser(user);
    app.bindLogout();
    resetCategoryForm();
    resetTagForm();
    app.setFeedback('category-feedback', '', false);
    app.setFeedback('tag-feedback', '', false);
    await loadLists();
  }

  init().catch(function (error) {
    app.showToast(error.message, 'error');
  });
})();
