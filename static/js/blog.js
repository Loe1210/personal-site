(function () {
  const root = window.PersonalSite;
  if (!root) return;

  const listEl = document.getElementById("blog-list");
  const categoryFilter = document.getElementById("category-filter");
  const tagFilter = document.getElementById("tag-filter");
  const summaryEl = document.getElementById("blog-filter-summary");

  const state = {
    category: "",
    tag: "",
    categories: [],
    tags: [],
  };

  function readStateFromURL() {
    const params = new URLSearchParams(window.location.search);
    state.category = params.get("category") || "";
    state.tag = params.get("tag") || "";
  }

  function writeStateToURL() {
    const params = new URLSearchParams();
    if (state.category) params.set("category", state.category);
    if (state.tag) params.set("tag", state.tag);
    const nextURL = params.toString() ? "/blog?" + params.toString() : "/blog";
    window.history.pushState({ category: state.category, tag: state.tag }, "", nextURL);
  }

  function renderChips(container, items, type, activeValue) {
    if (!container) return;
    const allLabel = type === "category" ? "全部分类" : "全部标签";
    const chips = [
      `<button class="filter-chip${activeValue === "" ? " is-active" : ""}" data-type="${type}" data-value="">${allLabel}</button>`
    ];

    items.forEach(function (item) {
      chips.push(
        `<button class="filter-chip${activeValue === item.slug ? " is-active" : ""}" data-type="${type}" data-value="${root.escapeHTML(item.slug)}">${root.escapeHTML(item.name)}</button>`
      );
    });

    container.innerHTML = chips.join("");
  }

  function buildCategoryMap() {
    const map = {};
    state.categories.forEach(function (item) {
      map[item.id] = item.name;
    });
    return map;
  }

  function buildTagMap() {
    const map = {};
    state.tags.forEach(function (item) {
      map[item.id] = item.name;
    });
    return map;
  }

  function getCategoryName(slug) {
    const match = state.categories.find(function (item) { return item.slug === slug; });
    return match ? match.name : "";
  }

  function getTagName(slug) {
    const match = state.tags.find(function (item) { return item.slug === slug; });
    return match ? match.name : "";
  }

  function renderSummary(total) {
    if (!summaryEl) return;

    const parts = [];
    const categoryName = state.category ? getCategoryName(state.category) : "";
    const tagName = state.tag ? getTagName(state.tag) : "";

    if (categoryName) {
      parts.push("分类: " + categoryName);
    }
    if (tagName) {
      parts.push("标签: #" + tagName);
    }

    if (!parts.length) {
      summaryEl.textContent = "当前展示全部已发布文章，可通过下方分类和标签快速切换视图。";
      return;
    }

    const totalText = typeof total === "number" ? " · 共 " + total + " 篇" : "";
    summaryEl.textContent = "当前筛选：" + parts.join(" / ") + totalText;
  }

  function cardHTML(article, categoryMap, tagMap) {
    const category = categoryMap[article.category_id] || "Uncategorized";
    const tags = (article.tag_ids || []).map(function (id) { return tagMap[id]; }).filter(Boolean);
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
          <h2 class="article-card__title">
            <a href="/blog/${root.escapeHTML(article.slug)}">${root.escapeHTML(article.title)}</a>
          </h2>
          <p class="article-card__summary">${root.escapeHTML(article.summary || "No summary yet")}</p>
          <div class="article-card__tags">
            ${tags.map(function (name) {
              return `<span class="article-card__tag">#${root.escapeHTML(name)}</span>`;
            }).join("")}
          </div>
        </div>
      </article>
    `;
  }

  async function loadTaxonomy() {
    const [categories, tags] = await Promise.all([
      root.fetchJSON("/api/categories"),
      root.fetchJSON("/api/tags")
    ]);

    state.categories = categories.list || [];
    state.tags = tags.list || [];
    renderChips(categoryFilter, state.categories, "category", state.category);
    renderChips(tagFilter, state.tags, "tag", state.tag);
    renderSummary();
  }

  async function loadArticles() {
    if (!listEl) return;
    listEl.innerHTML = '<div class="empty-state">正在加载文章...</div>';

    const params = new URLSearchParams();
    if (state.category) params.set("category", state.category);
    if (state.tag) params.set("tag", state.tag);

    try {
      const articleURL = params.toString() ? "/api/articles?" + params.toString() : "/api/articles";
      const articles = await root.fetchJSON(articleURL);
      const list = articles.list || [];
      const categoryMap = buildCategoryMap();
      const tagMap = buildTagMap();

      renderSummary(list.length);
      renderChips(categoryFilter, state.categories, "category", state.category);
      renderChips(tagFilter, state.tags, "tag", state.tag);

      if (!list.length) {
        listEl.innerHTML = '<div class="empty-state">当前筛选条件下没有文章。</div>';
        return;
      }

      listEl.innerHTML = list.map(function (article) {
        return cardHTML(article, categoryMap, tagMap);
      }).join("");
    } catch (error) {
      renderSummary();
      listEl.innerHTML = `<div class="empty-state">文章列表加载失败：${root.escapeHTML(error.message)}</div>`;
    }
  }

  async function refreshArticles() {
    await loadArticles();
  }

  document.addEventListener("click", function (event) {
    const chip = event.target.closest(".filter-chip");
    if (!chip) return;

    const type = chip.getAttribute("data-type");
    const value = chip.getAttribute("data-value") || "";
    if (type === "category") {
      state.category = value;
    }
    if (type === "tag") {
      state.tag = value;
    }

    writeStateToURL();
    refreshArticles();
  });

  window.addEventListener("popstate", function () {
    readStateFromURL();
    renderChips(categoryFilter, state.categories, "category", state.category);
    renderChips(tagFilter, state.tags, "tag", state.tag);
    refreshArticles();
  });

  async function init() {
    readStateFromURL();
    await loadTaxonomy();
    await loadArticles();
  }

  init();
})();
