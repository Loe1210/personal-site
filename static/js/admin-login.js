(function () {
  const app = window.AdminApp;
  if (!app) return;

  const form = document.getElementById('admin-login-form');
  const errorNode = document.getElementById('admin-login-error');
  const submitBtn = document.getElementById('admin-login-submit');

  async function checkSession() {
    try {
      await app.request('/api/admin/me', { headers: { Accept: 'application/json' } });
      window.location.href = '/admin';
    } catch (error) {
      if (error.status !== 401) {
        errorNode.hidden = false;
        errorNode.textContent = error.message;
      }
    }
  }

  form.addEventListener('submit', async function (event) {
    event.preventDefault();
    errorNode.hidden = true;
    submitBtn.disabled = true;
    try {
      await app.jsonRequest('/api/admin/login', 'POST', {
        username: document.getElementById('login-username').value.trim(),
        password: document.getElementById('login-password').value
      });
      window.location.href = '/admin';
    } catch (error) {
      errorNode.hidden = false;
      errorNode.textContent = error.message;
    } finally {
      submitBtn.disabled = false;
    }
  });

  checkSession();
})();
