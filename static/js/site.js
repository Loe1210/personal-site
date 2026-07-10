(function () {
  async function fetchJSON(url) {
    const response = await fetch(url, {
      headers: {
        Accept: "application/json",
      },
      credentials: "same-origin",
    });

    const data = await response.json();
    if (!response.ok || data.code !== 0) {
      throw new Error(data.message || "request failed");
    }
    return data.data;
  }

  function escapeHTML(input) {
    return String(input || "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  function markdownToHTML(markdown) {
    const safe = escapeHTML(markdown || "");

    const codeHandled = safe.replace(/```([\s\S]*?)```/g, function (_, code) {
      return `<pre><code>${code.trim()}</code></pre>`;
    });

    const lines = codeHandled.split(/\n{2,}/).map(function (block) {
      const trimmed = block.trim();
      if (!trimmed) {
        return "";
      }
      if (trimmed.startsWith("<pre><code>")) {
        return trimmed;
      }
      if (/^###\s+/.test(trimmed)) {
        return `<h3>${trimmed.replace(/^###\s+/, "")}</h3>`;
      }
      if (/^##\s+/.test(trimmed)) {
        return `<h2>${trimmed.replace(/^##\s+/, "")}</h2>`;
      }
      if (/^#\s+/.test(trimmed)) {
        return `<h1>${trimmed.replace(/^#\s+/, "")}</h1>`;
      }
      if (/^>\s+/.test(trimmed)) {
        return `<blockquote>${trimmed.replace(/^>\s+/, "")}</blockquote>`;
      }
      if (/^- /.test(trimmed)) {
        const items = trimmed
          .split("\n")
          .map(function (line) {
            return line.replace(/^- /, "").trim();
          })
          .filter(Boolean)
          .map(function (item) {
            return `<li>${item}</li>`;
          })
          .join("");
        return `<ul>${items}</ul>`;
      }
      return `<p>${trimmed.replace(/\n/g, "<br>")}</p>`;
    });

    return lines.join("");
  }

  window.PersonalSite = {
    fetchJSON: fetchJSON,
    escapeHTML: escapeHTML,
    markdownToHTML: markdownToHTML,
  };
})();
