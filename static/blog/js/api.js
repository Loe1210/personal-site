/**
 * 博客 API 封装
 * 使用 IIFE 避免全局变量污染，仅导出 BlogAPI 和 formatDate
 */
(function () {
    'use strict';

    var API_BASE = '/api/content';

    var categoryCache = {};
    var tagCache = {};

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
                throw new Error(data.msg || data.message || 'API error');
            }
            return data.data;
        });
    }

    function buildCategoryCache(categories) {
        categoryCache = {};
        categories.forEach(function(c) {
            categoryCache[c.id] = c;
        });
    }

    function buildTagCache(tags) {
        tagCache = {};
        tags.forEach(function(t) {
            tagCache[t.id] = t;
        });
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

    function estimateReadingTime(article) {
        article = article || {};
        var text = [article.content_md, article.content_html, article.summary, article.title]
            .filter(Boolean)
            .join(' ')
            .replace(/<[^>]*>/g, ' ');
        var chineseChars = (text.match(/[\u4e00-\u9fff]/g) || []).length;
        var englishWords = (text.replace(/[\u4e00-\u9fff]/g, ' ').match(/[A-Za-z0-9_]+(?:[-'][A-Za-z0-9_]+)*/g) || []).length;
        var minutes = (chineseChars / 350) + (englishWords / 200);
        return Math.max(1, Math.ceil(minutes || 1));
    }

    function mapBackendArticle(article) {
        article = article || {};
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
            cover: article.cover_image || '',
            category: categoryName,
            category_id: article.category_id,
            tags: tagNames,
            tag_ids: article.tag_ids || [],
            created_at: article.created_at,
            updated_at: article.updated_at,
            published_at: article.published_at,
            reading_time: article.reading_time || estimateReadingTime(article),
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
            var qs = '?page=' + (opts.page || 1) + '&page_size=' + (opts.limit || 10);
            if (opts.category) {
                qs += '&category=' + encodeURIComponent(opts.category);
            }
            if (opts.tag) {
                qs += '&tag=' + encodeURIComponent(opts.tag);
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
        },

        getPost: function (id) {
            return request('/articles/' + encodeURIComponent(id)).then(function(data) {
                return mapBackendArticle(data.article || data);
            });
        },

        getAdjacentPosts: function (id) {
            return request('/articles/' + encodeURIComponent(id) + '/adjacent').then(function(data) {
                data = data || {};
                return {
                    prev: data.prev || null,
                    next: data.next || null
                };
            });
        },

        getCategories: function () {
            return request('/categories').then(function(data) {
                var list = data.list || data || [];
                var mapped = list.map(mapBackendCategory);
                buildCategoryCache(mapped);
                return mapped;
            });
        },

        getTags: function () {
            return request('/tags').then(function(data) {
                var list = data.list || data || [];
                var mapped = list.map(mapBackendTag);
                buildTagCache(mapped);
                return mapped;
            });
        }
    };

    window.BlogAPI = BlogAPI;
    window.formatDate = formatDate;

})();
