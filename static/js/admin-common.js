(function () {
  const root = document.getElementById('admin-toast-root');

  function showToast(message, type) {
    if (!root) return;
    const toast = document.createElement('div');
    toast.className = 'admin-toast admin-toast--' + (type || 'success');
    toast.textContent = message;
    root.appendChild(toast);
    window.setTimeout(function () {
      toast.remove();
    }, 2800);
  }

  async function request(url, options) {
    const config = Object.assign({ credentials: 'same-origin' }, options || {});
    const response = await fetch(url, config);
    const contentType = response.headers.get('content-type') || '';
    const payload = contentType.includes('application/json') ? await response.json() : null;

    if (!response.ok || (payload && payload.code !== 0)) {
      const message = payload && payload.message ? payload.message : 'Request failed';
      const error = new Error(message);
      error.status = response.status;
      error.payload = payload;
      throw error;
    }

    return payload ? payload.data : null;
  }

  function jsonRequest(url, method, body) {
    return request(url, {
      method: method || 'GET',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      },
      body: body ? JSON.stringify(body) : undefined
    });
  }

  async function requireAuth() {
    try {
      const me = await request('/api/admin/me', {
        headers: { Accept: 'application/json' }
      });
      return me.user;
    } catch (error) {
      if (error.status === 401) {
        window.location.href = '/admin/login';
        return null;
      }
      throw error;
    }
  }

  function mountUser(user) {
    const node = document.getElementById('admin-user-name');
    if (node && user) {
      node.textContent = user.nickname ? user.nickname + ' (' + user.username + ')' : user.username;
    }
  }

  function bindLogout() {
    const btn = document.getElementById('admin-logout-btn');
    if (!btn) return;
    btn.addEventListener('click', async function () {
      try {
        await jsonRequest('/api/admin/logout', 'POST');
        window.location.href = '/admin/login';
      } catch (error) {
        showToast(error.message, 'error');
      }
    });
  }

  function setFeedback(id, message, success) {
    const node = document.getElementById(id);
    if (!node) return;
    if (!message) {
      node.hidden = true;
      node.textContent = '';
      node.classList.remove('admin-success');
      return;
    }
    node.hidden = false;
    node.textContent = message;
    node.classList.toggle('admin-success', !!success);
  }

  function escapeHTML(input) {
    return String(input || '')
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  function collectCheckedValues(containerSelector) {
    return Array.from(document.querySelectorAll(containerSelector + ' input[type="checkbox"]:checked')).map(function (item) {
      return Number(item.value);
    });
  }

  window.AdminApp = {
    request: request,
    jsonRequest: jsonRequest,
    requireAuth: requireAuth,
    mountUser: mountUser,
    bindLogout: bindLogout,
    showToast: showToast,
    setFeedback: setFeedback,
    escapeHTML: escapeHTML,
    collectCheckedValues: collectCheckedValues
  };
})();
