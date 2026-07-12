/**
 * 博客 API 封装
 * 使用 IIFE 避免全局变量污染，仅导出 BlogAPI 和 formatDate
 */
(function () {
    'use strict';

    var API_BASE = '/api';
    var USE_MOCK = false;

    var categoryCache = {};
    var tagCache = {};
    var categoriesLoaded = false;
    var tagsLoaded = false;

var MOCK_POSTS = [
    {
        id: 'hello-internship',
        slug: 'hello-internship',
        title: '实习第一周：从校园到工位的思维切换',
        summary: '记录入职第一周的环境配置、代码仓库上手、第一份 PR 全过程，以及从学生到工程师的心态转变。',
        cover: '',
        category: '实习笔记',
        tags: ['实习', '入职', 'Git'],
        created_at: '2026-07-01T09:00:00+08:00',
        updated_at: '2026-07-01T09:00:00+08:00',
        reading_time: 6,
        content: [
            '# 实习第一周：从校园到工位的思维切换',
            '',
            '入职第一周，最大的感受是：**工程世界和学生世界的距离，比想象中大得多**。',
            '',
            '## 一、环境配置：第一天就踩坑',
            '',
            '配环境永远是新人的第一道门槛。公司用的技术栈和学校学的有重叠，但版本、规范、CI 流程完全不同。',
            '',
            '```bash',
            '# 克隆代码仓库',
            'git clone git@github.com:company/repo.git',
            'cd repo',
            '',
            '# 安装依赖',
            'npm install',
            '',
            '# 启动开发服务',
            'npm run dev',
            '```',
            '',
            '> 提示：遇到依赖冲突先别急着 `--force`，先看 README 和 CONTRIBUTING.md。',
            '',
            '## 二、第一份 PR：小而完整',
            '',
            '导师没让我一上来就改核心逻辑，而是分配了一个文档错别字 + 一个单元测试补充的小任务。这其实是培养**完整提交工作流**的好机会：',
            '',
            '1. 新建分支 `fix/typo-in-readme`',
            '2. 改动 → 本地测试 → commit',
            '3. push → 在 GitHub 上开 PR',
            '4. 关联 issue → 等待 review',
            '5. 根据 review 意见 amend → 重新 push',
            '6. merge',
            '',
            '## 三、心态转变',
            '',
            '- 学校里写代码：**能跑就行**',
            '- 工程里写代码：**可读、可测、可维护、可回滚**',
            '',
            '这三条会贯穿整个实习，先记下来，慢慢体会。',
            ''
        ].join('\n')
    },
    {
        id: 'learn-react-hooks',
        slug: 'learn-react-hooks',
        title: 'React Hooks 学习笔记：useEffect 的依赖陷阱',
        summary: '从一次内存泄漏 bug 出发，深入理解 useEffect 的依赖数组、清理函数和闭包陷阱。',
        cover: '',
        category: '学习笔记',
        tags: ['React', 'Hooks', '前端'],
        created_at: '2026-07-05T20:30:00+08:00',
        updated_at: '2026-07-05T20:30:00+08:00',
        reading_time: 8,
        content: [
            '# React Hooks 学习笔记：useEffect 的依赖陷阱',
            '',
            '实习项目里遇到一个组件卸载后还在 setState 的警告，追根溯源是 `useEffect` 的依赖没写对。',
            '',
            '## 问题复现',
            '',
            '```jsx',
            'function Timer() {',
            '  const [count, setCount] = useState(0);',
            '  useEffect(() => {',
            '    setInterval(() => setCount(c => c + 1), 1000);',
            '  }, []);',
            '  return <div>{count}</div>;',
            '}',
            '```',
            '',
            '组件卸载后，定时器仍在运行，setState 报警告。',
            '',
            '## 正确写法',
            '',
            '```jsx',
            'useEffect(() => {',
            '  const id = setInterval(() => setCount(c => c + 1), 1000);',
            '  return () => clearInterval(id);',
            '}, []);',
            '```',
            '',
            '## 三条经验',
            '',
            '1. **凡是有副作用，必有清理**',
            '2. **依赖数组不要撒谎**，该写的都得写',
            '3. 想跳过依赖用 `useRef`，别用空数组糊弄',
            ''
        ].join('\n')
    }
];

var MOCK_CATEGORIES = [
    { name: '实习笔记', slug: 'internship', count: 1 },
    { name: '学习笔记', slug: 'learning', count: 1 }
];

var MOCK_TAGS = [
    { name: '实习', slug: 'internship', count: 1 },
    { name: 'React', slug: 'react', count: 1 },
    { name: 'Hooks', slug: 'hooks', count: 1 },
    { name: '前端', slug: 'frontend', count: 1 }
];

function formatDate(iso) {
    if (!iso) return '';
    var d = new Date(iso);
    if (isNaN(d.getTime())) return '';
    var pad = function (n) { return n < 10 ? '0' + n : '' + n; };
    return d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate());
}

function request(path, options) {
    options = options || {};
    return fetch(API_BASE + path, {
        method: options.method || 'GET',
        headers: Object.assign({
            'Content-Type': 'application/json'
        }, options.headers || {}),
        body: options.body ? JSON.stringify(options.body) : undefined,
        credentials: 'same-origin'
    }).then(function (res) {
        if (!res.ok) throw new Error('API ' + path + ' returned ' + res.status);
        return res.json();
    }).then(function (data) {
        if (data.code !== 0) {
            throw new Error(data.message || 'API error');
        }
        return data.data;
    });
}

var STORAGE_KEY = 'hins_blog_posts';
var CATEGORIES_KEY = 'hins_blog_categories';
var TAGS_KEY = 'hins_blog_tags';

function loadManagedPosts() {
    try {
        var stored = localStorage.getItem(STORAGE_KEY);
        if (stored) return JSON.parse(stored);
    } catch (e) {}
    return null;
}

function loadManagedCategories() {
    try {
        var stored = localStorage.getItem(CATEGORIES_KEY);
        if (stored) return JSON.parse(stored);
    } catch (e) {}
    return null;
}

function loadManagedTags() {
    try {
        var stored = localStorage.getItem(TAGS_KEY);
        if (stored) return JSON.parse(stored);
    } catch (e) {}
    return null;
}

function buildCategoryCache(categories) {
    categoryCache = {};
    categories.forEach(function(c) {
        categoryCache[c.id] = c;
    });
    categoriesLoaded = true;
}

function buildTagCache(tags) {
    tagCache = {};
    tags.forEach(function(t) {
        tagCache[t.id] = t;
    });
    tagsLoaded = true;
}

function getCategoryNameById(categoryId) {
    if (categoryId && categoryCache[categoryId]) {
        return categoryCache[categoryId].name;
    }
    return '';
}

function getTagNamesByIds(tagIds) {
    if (!tagIds || !Array.isArray(tagIds)) return [];
    return tagIds.map(function(id) {
        if (tagCache[id]) {
            return tagCache[id].name;
        }
        return null;
    }).filter(function(name) { return name !== null; });
}

function mapBackendArticle(article) {
    var categoryName = '';
    if (article.category_name) {
        categoryName = article.category_name;
    } else if (article.category && article.category.name) {
        categoryName = article.category.name;
    } else if (article.category_id) {
        categoryName = getCategoryNameById(article.category_id);
    }

    var tagNames = [];
    if (article.tags && Array.isArray(article.tags)) {
        tagNames = article.tags.map(function(t) { return t.name || t; });
    } else if (article.tag_ids && Array.isArray(article.tag_ids)) {
        tagNames = getTagNamesByIds(article.tag_ids);
    }

    return {
        id: article.id,
        slug: article.slug,
        title: article.title,
        summary: article.summary,
        cover: article.cover_image,
        category: categoryName,
        category_id: article.category_id,
        tags: tagNames,
        tag_ids: article.tag_ids || [],
        created_at: article.created_at,
        updated_at: article.updated_at,
        published_at: article.published_at,
        reading_time: article.reading_time || 5,
        content: article.content_md || article.content_html || article.content,
        status: article.status
    };
}

function mapBackendCategory(cat) {
    return {
        id: cat.id,
        name: cat.name,
        slug: cat.slug,
        count: cat.article_count || 0
    };
}

function mapBackendTag(tag) {
    return {
        id: tag.id,
        name: tag.name,
        slug: tag.slug,
        count: tag.article_count || 0
    };
}

function getCategorySlugByName(name) {
    for (var id in categoryCache) {
        if (categoryCache[id].name === name) {
            return categoryCache[id].slug;
        }
    }
    return name;
}

function getTagSlugByName(name) {
    for (var id in tagCache) {
        if (tagCache[id].name === name) {
            return tagCache[id].slug;
        }
    }
    return name;
}

var BlogAPI = {
    getPosts: function (opts) {
        opts = opts || {};
        var self = this;
        return self.ensureMetadataLoaded().then(function() {
            if (USE_MOCK) {
                var managed = loadManagedPosts();
                var allPosts = managed || MOCK_POSTS;
                var page = opts.page || 1;
                var limit = opts.limit || 10;
                var list = allPosts.slice().sort(function (a, b) {
                    return new Date(b.created_at) - new Date(a.created_at);
                });
                list = list.filter(function (p) { return p.status !== 'draft'; });
                if (opts.category) {
                    list = list.filter(function (p) { return p.category === opts.category; });
                }
                if (opts.tag) {
                    list = list.filter(function (p) { return (p.tags || []).indexOf(opts.tag) !== -1; });
                }
                if (opts.search) {
                    var q = opts.search.toLowerCase();
                    list = list.filter(function (p) {
                        return (p.title || '').toLowerCase().indexOf(q) !== -1 ||
                               (p.summary || '').toLowerCase().indexOf(q) !== -1;
                    });
                }
                var total = list.length;
                var start = (page - 1) * limit;
                var posts = list.slice(start, start + limit).map(function (p) {
                    var copy = {};
                    for (var k in p) {
                        if (k !== 'content') copy[k] = p[k];
                    }
                    return copy;
                });
                return { posts: posts, total: total, page: page, limit: limit };
            }
            var qs = '?page=' + (opts.page || 1) + '&page_size=' + (opts.limit || 10);
            if (opts.category) {
                var categorySlug = getCategorySlugByName(opts.category);
                qs += '&category=' + encodeURIComponent(categorySlug);
            }
            if (opts.tag) {
                var tagSlug = getTagSlugByName(opts.tag);
                qs += '&tag=' + encodeURIComponent(tagSlug);
            }
            if (opts.search) qs += '&keyword=' + encodeURIComponent(opts.search);
            return request('/articles' + qs).then(function(data) {
                var articles = (data.list || []).map(mapBackendArticle);
                return {
                    posts: articles,
                    total: data.total || articles.length,
                    page: data.page || opts.page || 1,
                    limit: data.page_size || opts.limit || 10
                };
            });
        });
    },

    getPost: function (id) {
        var self = this;
        return self.ensureMetadataLoaded().then(function() {
            if (USE_MOCK) {
                var managed = loadManagedPosts();
                var allPosts = managed || MOCK_POSTS;
                var found = allPosts.filter(function (p) { return String(p.id) === String(id); })[0];
                return found
                    ? found
                    : Promise.reject(new Error('post not found'));
            }
            return request('/articles/id/' + encodeURIComponent(id)).then(function(data) {
                return mapBackendArticle(data.article);
            });
        });
    },

    getAdjacentPosts: function (slugOrId) {
        return Promise.resolve({ prev: null, next: null });
    },

    getCategories: function () {
        if (USE_MOCK) {
            var managedPosts = loadManagedPosts();
            var allPosts = managedPosts || MOCK_POSTS;
            var managedCats = loadManagedCategories();
            var cats = {};
            if (managedCats) {
                managedCats.forEach(function (c) { cats[c] = 0; });
            }
            allPosts.forEach(function (p) {
                if (p.category) cats[p.category] = (cats[p.category] || 0) + 1;
            });
            var catList = Object.keys(cats).map(function (name) {
                return { id: name, name: name, slug: name, count: cats[name] };
            });
            buildCategoryCache(catList);
            return Promise.resolve(catList);
        }
        var self = this;
        return request('/categories').then(function(data) {
            var list = data.list || data || [];
            var mapped = list.map(mapBackendCategory);
            
            return self._getPublishedArticlesForCount().then(function(articles) {
                var catCount = {};
                articles.forEach(function(article) {
                    var catId = article.category_id;
                    if (catId) {
                        catCount[catId] = (catCount[catId] || 0) + 1;
                    }
                });
                mapped.forEach(function(cat) {
                    cat.count = catCount[cat.id] || 0;
                });
                buildCategoryCache(mapped);
                return mapped;
            });
        });
    },

    getTags: function () {
        if (USE_MOCK) {
            var managedPosts = loadManagedPosts();
            var allPosts = managedPosts || MOCK_POSTS;
            var managedTags = loadManagedTags();
            var tags = {};
            if (managedTags) {
                managedTags.forEach(function (t) { tags[t] = 0; });
            }
            allPosts.forEach(function (p) {
                (p.tags || []).forEach(function (t) {
                    tags[t] = (tags[t] || 0) + 1;
                });
            });
            var tagList = Object.keys(tags).map(function (name) {
                return { id: name, name: name, slug: name, count: tags[name] };
            });
            buildTagCache(tagList);
            return Promise.resolve(tagList);
        }
        var self = this;
        return request('/tags').then(function(data) {
            var list = data.list || data || [];
            var mapped = list.map(mapBackendTag);
            
            return self._getPublishedArticlesForCount().then(function(articles) {
                var tagCount = {};
                articles.forEach(function(article) {
                    (article.tag_ids || []).forEach(function(tagId) {
                        tagCount[tagId] = (tagCount[tagId] || 0) + 1;
                    });
                });
                mapped.forEach(function(tag) {
                    tag.count = tagCount[tag.id] || 0;
                });
                buildTagCache(mapped);
                return mapped;
            });
        });
    },

    _allPublishedArticles: null,

    _getPublishedArticlesForCount: function() {
        if (this._allPublishedArticles) {
            return Promise.resolve(this._allPublishedArticles);
        }
        var self = this;
        var allArticles = [];
        var currentPage = 1;
        var pageSize = 100;

        function fetchPage() {
            return request('/articles?page=' + currentPage + '&page_size=' + pageSize).then(function(data) {
                var list = data.list || [];
                allArticles = allArticles.concat(list);
                if (data.total > currentPage * pageSize) {
                    currentPage++;
                    return fetchPage();
                }
                self._allPublishedArticles = allArticles;
                return allArticles;
            });
        }
        return fetchPage();
    },

    ensureMetadataLoaded: function() {
        var promises = [];
        if (!categoriesLoaded) {
            promises.push(this.getCategories());
        }
        if (!tagsLoaded) {
            promises.push(this.getTags());
        }
        if (promises.length === 0) {
            return Promise.resolve();
        }
        return Promise.all(promises);
    }
};

window.BlogAPI = BlogAPI;
window.formatDate = formatDate;

})();
