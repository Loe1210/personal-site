(function () {
  const root = window.PersonalSite;
  if (!root) return;

  const latestList = document.getElementById("latest-posts-list");
  const tagCloud = document.getElementById("home-tag-cloud");

  function createPostCard(article, categories, tags) {
    const category = categories[article.category_id] || "Uncategorized";
    const tagNames = (article.tag_ids || []).map(function (id) {
      return tags[id];
    }).filter(Boolean);
    const cover = article.cover_image
      ? `<img class="article-card__cover-img" src="${root.escapeHTML(article.cover_image)}" alt="${root.escapeHTML(article.title)}">`
      : '<div class="article-card__cover-glow"></div><div class="article-card__cover-grid"></div>';

    return `
      <article class="article-card">
        <div class="article-card__cover">
          ${cover}
        </div>
        <div class="article-card__body">
          <div class="article-card__meta">
            <span class="article-card__category">${root.escapeHTML(category)}</span>
            <span class="article-card__date">${root.escapeHTML(article.published_at || article.created_at || "")}</span>
          </div>
          <h3 class="article-card__title">
            <a href="/blog/${root.escapeHTML(article.slug)}">${root.escapeHTML(article.title)}</a>
          </h3>
          <p class="article-card__summary">${root.escapeHTML(article.summary || "No summary yet")}</p>
          <div class="article-card__tags">
            ${tagNames.map(function (name) {
              return `<span class="article-card__tag">#${root.escapeHTML(name)}</span>`;
            }).join("")}
          </div>
        </div>
      </article>
    `;
  }

  async function loadHome() {
    if (!latestList) return;

    try {
      const [articles, categories, tags] = await Promise.all([
        root.fetchJSON("/api/articles?page=1&page_size=3"),
        root.fetchJSON("/api/categories"),
        root.fetchJSON("/api/tags"),
      ]);

      const categoryMap = {};
      (categories.list || []).forEach(function (item) {
        categoryMap[item.id] = item.name;
      });

      const tagMap = {};
      (tags.list || []).forEach(function (item) {
        tagMap[item.id] = item.name;
      });

      const items = (articles.list || []).slice(0, 3);
      if (!items.length) {
        latestList.innerHTML = '<div class="empty-state">还没有已发布的文章。</div>';
      } else {
        latestList.innerHTML = items.map(function (article) {
          return createPostCard(article, categoryMap, tagMap);
        }).join("");
      }

      const tagNames = (tags.list || []).slice(0, 8).map(function (item) {
        return `<span>${root.escapeHTML(item.name)}</span>`;
      });
      if (tagCloud && tagNames.length) {
        tagCloud.innerHTML = tagNames.join("");
      }
    } catch (error) {
      latestList.innerHTML = `<div class="empty-state">首页数据加载失败：${root.escapeHTML(error.message)}</div>`;
    }
  }

  loadHome();
})();
