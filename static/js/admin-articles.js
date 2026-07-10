(function () {
  const app = window.AdminApp;
  if (!app) return;

  const listNode = document.getElementById('articles-list');
  const searchInput = document.getElementById('article-search-input');
  const statusFilter = document.getElementById('article-status-filter');
  const filterBtn = document.getElementById('article-filter-btn');

  function articleCard(item) {
    return `
      <article class="admin-item-card">
        <div class="admin-item-card__meta">
          <span class="admin-badge">${app.escapeHTML(item.status || '')}</span>
          <span>${app.escapeHTML(item.updated_at || item.created_at || '')}</span>
          <span>${app.escapeHTML(item.slug || '')}</span>
        </div>
        <h3 class="admin-item-card__title">${app.escapeHTML(item.title || 'Untitled')}</h3>
        <p class="admin-item-card__desc">${app.escapeHTML(item.summary || 'No summary')}</p>
        <div class="admin-item-card__actions">
          <a class="admin-button admin-button--ghost" href="/admin/articles/${item.id}/edit">Edit</a>
          <button class="admin-button admin-button--ghost" type="button" data-delete-id="${item.id}">Delete</button>
        </div>
      </article>
    `;
  }

  async function loadArticles() {
    listNode.innerHTML = '<div class="admin-empty">Loading articles...</div>';
    const params = new URLSearchParams();
    if (searchInput.value.trim()) params.set('keyword', searchInput.value.trim());
    if (statusFilter.value) params.set('status', statusFilter.value);

    try {
      const user = await app.requireAuth();
      if (!user) return;
      app.mountUser(user);
      app.bindLogout();
      const data = await app.request('/api/admin/articles?' + params.toString(), { headers: { Accept: 'application/json' } });
      if (!data.list || !data.list.length) {
        listNode.innerHTML = '<div class="admin-empty">No articles found.</div>';
        return;
      }
      listNode.innerHTML = data.list.map(articleCard).join('');
    } catch (error) {
      app.setFeedback('articles-feedback', error.message, false);
      listNode.innerHTML = '<div class="admin-empty">Failed to load articles.</div>';
    }
  }

  filterBtn.addEventListener('click', loadArticles);
  document.addEventListener('click', async function (event) {
    const btn = event.target.closest('[data-delete-id]');
    if (!btn) return;
    const id = btn.getAttribute('data-delete-id');
    if (!window.confirm('Delete this article?')) return;
    try {
      await app.requireAuth();
      await app.request('/api/admin/articles/' + id, { method: 'DELETE', headers: { Accept: 'application/json' } });
      app.showToast('Article deleted', 'success');
      loadArticles();
    } catch (error) {
      app.showToast(error.message, 'error');
    }
  });

  loadArticles();
})();
