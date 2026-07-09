package article

import (
	"context"
	"strings"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
)

var mockArticles = []*articlemodel.Article{
	{
		ID:          1,
		Title:       "Hertz 路由、中间件与分层实践",
		Slug:        "hertz-router-middleware-layering",
		Summary:     "记录我学习 Hertz 路由、中间件和项目分层时的理解。",
		ContentMd:   "# Hertz 路由、中间件与分层实践\n\n这是第一篇文章内容。",
		ContentHTML: "<h1>Hertz 路由、中间件与分层实践</h1><p>这是第一篇文章内容。</p>",
		CoverImage:  "",
		CategoryID:  1,
		TagIds:      []int64{1, 2},
		Status:      "published",
		CreatedAt:   "2026-07-09 11:00:00",
		UpdatedAt:   "2026-07-09 11:00:00",
		PublishedAt: "2026-07-09 11:00:00",
	},
	{
		ID:          2,
		Title:       "实习里遇到的 Go 小坑",
		Slug:        "go-pitfalls-internship",
		Summary:     "总结实习过程中遇到的一些 Go 常见问题。",
		ContentMd:   "# 实习里遇到的 Go 小坑\n\n这是第二篇文章内容。",
		ContentHTML: "<h1>实习里遇到的 Go 小坑</h1><p>这是第二篇文章内容。</p>",
		CoverImage:  "",
		CategoryID:  2,
		TagIds:      []int64{2, 3},
		Status:      "published",
		CreatedAt:   "2026-07-09 11:10:00",
		UpdatedAt:   "2026-07-09 11:10:00",
		PublishedAt: "2026-07-09 11:10:00",
	},
}

func ListPublicArticles(_ context.Context, req *articlemodel.ListArticlesRequest) (*articlemodel.ListArticlesResponse, error) {
	list := make([]*articlemodel.Article, 0)

	keyword := strings.TrimSpace(req.Keyword)
	for _, item := range mockArticles {
		if item.Status != "published" {
			continue
		}
		if keyword != "" && !strings.Contains(strings.ToLower(item.Title), strings.ToLower(keyword)) {
			continue
		}
		list = append(list, item)
	}

	page := int64(1)
	pageSize := int64(len(list))
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	return &articlemodel.ListArticlesResponse{
		List:     list,
		Total:    int64(len(list)),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func GetPublicArticleBySlug(_ context.Context, req *articlemodel.GetArticleBySlugRequest) (*articlemodel.GetArticleResponse, error) {
	for _, item := range mockArticles {
		if item.Slug == req.Slug && item.Status == "published" {
			return &articlemodel.GetArticleResponse{
				Article: item,
			}, nil
		}
	}
	return nil, nil
}

func ListAdminArticles(_ context.Context, req *articlemodel.ListArticlesRequest) (*articlemodel.ListArticlesResponse, error) {
	list := make([]*articlemodel.Article, 0)

	keyword := strings.TrimSpace(req.Keyword)
	for _, item := range mockArticles {
		if keyword != "" && !strings.Contains(strings.ToLower(item.Title), strings.ToLower(keyword)) {
			continue
		}
		list = append(list, item)
	}

	page := int64(1)
	pageSize := int64(len(list))
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	return &articlemodel.ListArticlesResponse{
		List:     list,
		Total:    int64(len(list)),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func CreateArticle(_ context.Context, req *articlemodel.CreateArticleRequest) (*articlemodel.CreateArticleResponse, error) {
	nextID := int64(len(mockArticles) + 1)

	article := &articlemodel.Article{
		ID:          nextID,
		Title:       req.Title,
		Slug:        req.Slug,
		Summary:     req.Summary,
		ContentMd:   req.ContentMd,
		ContentHTML: req.ContentMd,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		TagIds:      req.TagIds,
		Status:      req.Status,
		CreatedAt:   "2026-07-09 12:00:00",
		UpdatedAt:   "2026-07-09 12:00:00",
		PublishedAt: "",
	}

	if article.Status == "" {
		article.Status = "draft"
	}

	mockArticles = append(mockArticles, article)

	return &articlemodel.CreateArticleResponse{
		Article: article,
		Message: "article created",
	}, nil
}

func UpdateArticle(_ context.Context, req *articlemodel.UpdateArticleRequest) (*articlemodel.UpdateArticleResponse, error) {
	for _, item := range mockArticles {
		if item.ID != req.ID {
			continue
		}

		item.Title = req.Title
		item.Slug = req.Slug
		item.Summary = req.Summary
		item.ContentMd = req.ContentMd
		item.ContentHTML = req.ContentMd
		item.CoverImage = req.CoverImage
		item.CategoryID = req.CategoryID
		item.TagIds = req.TagIds
		item.Status = req.Status
		item.UpdatedAt = "2026-07-09 12:30:00"

		if item.Status == "published" && item.PublishedAt == "" {
			item.PublishedAt = "2026-07-09 12:30:00"
		}
		if item.Status != "published" {
			item.PublishedAt = ""
		}

		return &articlemodel.UpdateArticleResponse{
			Article: item,
			Message: "article updated",
		}, nil
	}

	return nil, nil
}

func DeleteArticle(_ context.Context, req *articlemodel.DeleteArticleRequest) (*articlemodel.DeleteArticleResponse, error) {
	for i, item := range mockArticles {
		if item.ID != req.ID {
			continue
		}

		mockArticles = append(mockArticles[:i], mockArticles[i+1:]...)

		return &articlemodel.DeleteArticleResponse{
			Success: true,
			Message: "article deleted",
		}, nil
	}

	return &articlemodel.DeleteArticleResponse{
		Success: false,
		Message: "article not found",
	}, nil
}