(function () {
  const root = window.PersonalSite;
  if (!root) return;

  const page = document.querySelector(".article-page");
  if (!page) return;

  const slug = page.getAttribute("data-slug");
  const titleEl = document.getElementById("article-title");
  const summaryEl = document.getElementById("article-summary");
  const metaEl = document.getElementById("article-meta");
  const bodyEl = document.getElementById("article-body");
  const coverWrap = document.getElementById("article-cover-wrap");
  const coverEl = document.getElementById("article-cover");

  function chip(text) {
    return `<span class="article-meta-chip">${root.escapeHTML(text)}</span>`;
  }

  async function loadArticle() {
    try {
      const [articleResp, categoriesResp, tagsResp] = await Promise.all([
        root.fetchJSON(`/api/articles/${encodeURIComponent(slug)}`),
        root.fetchJSON("/api/categories"),
        root.fetchJSON("/api/tags"),
      ]);

      const article = articleResp.article;
      const categoryMap = {};
      (categoriesResp.list || []).forEach(function (item) { categoryMap[item.id] = item.name; });
      const tagMap = {};
      (tagsResp.list || []).forEach(function (item) { tagMap[item.id] = item.name; });

      titleEl.textContent = article.title || "未命名文章";
      summaryEl.textContent = article.summary || "这篇文章还没有摘要。";

      const metaParts = [];
      if (article.published_at || article.created_at) {
        metaParts.push(chip(article.published_at || article.created_at));
      }
      if (article.category_id && categoryMap[article.category_id]) {
        metaParts.push(chip(categoryMap[article.category_id]));
      }
      (article.tag_ids || []).forEach(function (id) {
        if (tagMap[id]) {
          metaParts.push(chip(`#${tagMap[id]}`));
        }
      });
      metaEl.innerHTML = metaParts.join("");

      if (article.cover_image) {
        coverEl.src = article.cover_image;
        coverEl.alt = article.title || "article cover";
        coverWrap.hidden = false;
      }

      bodyEl.innerHTML = root.markdownToHTML(article.content_md || article.content_html || "");
    } catch (error) {
      titleEl.textContent = "文章加载失败";
      summaryEl.textContent = error.message;
      bodyEl.innerHTML = `<div class="empty-state">正文加载失败：${root.escapeHTML(error.message)}</div>`;
    }
  }

  loadArticle();
})();
