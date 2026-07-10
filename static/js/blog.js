(function () {
  const root = window.PersonalSite;
  if (!root) return;

  const listEl = document.getElementById("blog-list");
  const categoryFilter = document.getElementById("category-filter");
  const tagFilter = document.getElementById("tag-filter");

  const state = {
    category: "",
    tag: "",
    categories: [],
    tags: [],
  };

  function renderChips(container, items, type, activeValue) {
    if (!container) return;
    const allLabel = type === "category" ? "All Categories" : "All Tags";
    const chips = [`<button class="filter-chip${activeValue === "" ? " is-active" : ""}" data-type="${type}" data-value="">${allLabel}</button>`];

    items.forEach(function (item) {
      chips.push(`<button class="filter-chip${activeValue === item.slug ? " is-active" : ""}" data-type="${type}" data-value="${root.escapeHTML(item.slug)}">${root.escapeHTML(item.name)}</button>`);
    });

    container.innerHTML = chips.join("");
  }

  function cardHTML(article, categoryMap, tagMap) {
    const category = categoryMap[article.category_id] || "Uncategorized";
    const tags = (article.tag_ids || []).map(function (id) { return tagMap[id]; }).filter(Boolean);

    return `
      <article class="article-card">
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
      </article>
    `;
  }

  async function loadList() {
    if (!listEl) return;
    listEl.innerHTML = '<div class="empty-state">Loading posts...</div>';

    try {
      const [articles, categories, tags] = await Promise.all([
        root.fetchJSON("/api/articles"),
        root.fetchJSON("/api/categories"),
        root.fetchJSON("/api/tags"),
      ]);

      state.categories = categories.list || [];
      state.tags = tags.list || [];

      renderChips(categoryFilter, state.categories, "category", state.category);
      renderChips(tagFilter, state.tags, "tag", state.tag);

      const categoryMap = {};
      state.categories.forEach(function (item) { categoryMap[item.id] = item.name; });

      const tagMap = {};
      const tagSlugMap = {};
      state.tags.forEach(function (item) {
        tagMap[item.id] = item.name;
        tagSlugMap[item.id] = item.slug;
      });

      let filtered = articles.list || [];
      if (state.category) {
        const categoryIDs = state.categories.filter(function (item) { return item.slug === state.category; }).map(function (item) { return item.id; });
        filtered = filtered.filter(function (item) { return categoryIDs.includes(item.category_id); });
      }
      if (state.tag) {
        filtered = filtered.filter(function (item) {
          return (item.tag_ids || []).some(function (id) { return tagSlugMap[id] === state.tag; });
        });
      }

      if (!filtered.length) {
        listEl.innerHTML = '<div class="empty-state">No posts match the current filters.</div>';
        return;
      }

      listEl.innerHTML = filtered.map(function (article) {
        return cardHTML(article, categoryMap, tagMap);
      }).join("");
    } catch (error) {
      listEl.innerHTML = `<div class="empty-state">Post list failed to load: ${root.escapeHTML(error.message)}</div>`;
    }
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
    loadList();
  });

  loadList();
})();
