(function () {
  const app = window.AdminApp;
  if (!app) return;

  async function init() {
    try {
      const user = await app.requireAuth();
      if (!user) return;
      app.mountUser(user);
      app.bindLogout();

      const [articles, categories, tags] = await Promise.all([
        app.request('/api/admin/articles', { headers: { Accept: 'application/json' } }),
        app.request('/api/admin/categories', { headers: { Accept: 'application/json' } }),
        app.request('/api/admin/tags', { headers: { Accept: 'application/json' } })
      ]);

      document.getElementById('dashboard-greeting').textContent = 'Welcome back, ' + (user.nickname || user.username);
      document.getElementById('dashboard-subtitle').textContent = 'Your session is active. Use the shortcuts below to manage content.';
      document.getElementById('stat-articles').textContent = articles.total || 0;
      document.getElementById('stat-published').textContent = (articles.list || []).filter(function (item) { return item.status === 'published'; }).length;
      document.getElementById('stat-categories').textContent = (categories.list || []).length;
      document.getElementById('stat-tags').textContent = (tags.list || []).length;
    } catch (error) {
      app.showToast(error.message, 'error');
    }
  }

  init();
})();
