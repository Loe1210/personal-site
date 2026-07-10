(function () {
  const app = window.AdminApp;
  if (!app) return;

  function taxonomyRow(item) {
    return `
      <article class="taxonomy-row">
        <strong>${app.escapeHTML(item.name)}</strong>
        <div class="taxonomy-row__slug">${app.escapeHTML(item.slug)}</div>
        <p class="taxonomy-row__desc">${app.escapeHTML(item.description || 'No description')}</p>
      </article>
    `;
  }

  async function loadLists() {
    const user = await app.requireAuth();
    if (!user) return;
    app.mountUser(user);
    app.bindLogout();

    const [categories, tags] = await Promise.all([
      app.request('/api/admin/categories', { headers: { Accept: 'application/json' } }),
      app.request('/api/admin/tags', { headers: { Accept: 'application/json' } })
    ]);

    document.getElementById('category-list').innerHTML = (categories.list || []).length
      ? categories.list.map(taxonomyRow).join('')
      : '<div class="admin-empty">No categories yet.</div>';
    document.getElementById('tag-list').innerHTML = (tags.list || []).length
      ? tags.list.map(taxonomyRow).join('')
      : '<div class="admin-empty">No tags yet.</div>';
  }

  document.getElementById('category-form').addEventListener('submit', async function (event) {
    event.preventDefault();
    try {
      await app.jsonRequest('/api/admin/categories', 'POST', {
        name: document.getElementById('category-name').value.trim(),
        slug: document.getElementById('category-slug').value.trim(),
        description: document.getElementById('category-description').value.trim()
      });
      app.setFeedback('category-feedback', 'Category created.', true);
      event.target.reset();
      loadLists();
    } catch (error) {
      app.setFeedback('category-feedback', error.message, false);
    }
  });

  document.getElementById('tag-form').addEventListener('submit', async function (event) {
    event.preventDefault();
    try {
      await app.jsonRequest('/api/admin/tags', 'POST', {
        name: document.getElementById('tag-name').value.trim(),
        slug: document.getElementById('tag-slug').value.trim(),
        description: document.getElementById('tag-description').value.trim()
      });
      app.setFeedback('tag-feedback', 'Tag created.', true);
      event.target.reset();
      loadLists();
    } catch (error) {
      app.setFeedback('tag-feedback', error.message, false);
    }
  });

  loadLists().catch(function (error) {
    app.showToast(error.message, 'error');
  });
})();
