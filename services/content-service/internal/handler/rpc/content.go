package rpc

import (
	"context"

	kitexcontent "github.com/Loe1210/personal-site/kitex_gen/content"
	"github.com/Loe1210/personal-site/services/content-service/internal/model"
	"github.com/Loe1210/personal-site/services/content-service/internal/service"
)

type Handler struct {
	articles *service.ArticleService
}

func NewHandler(articles *service.ArticleService) *Handler {
	return &Handler{articles: articles}
}

func (h *Handler) GetArticleByID(ctx context.Context, req *kitexcontent.GetArticleByIDRequest) (*kitexcontent.GetArticleResponse, error) {
	article, err := h.articles.GetArticleByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &kitexcontent.GetArticleResponse{Article: toPBArticle(article)}, nil
}

func (h *Handler) ListPublicArticles(ctx context.Context, req *kitexcontent.ListPublicArticlesRequest) (*kitexcontent.ListArticlesResponse, error) {
	result, err := h.articles.ListPublicArticles(ctx, model.ListFilter{
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
		Keyword:  req.GetKeyword(),
	})
	if err != nil {
		return nil, err
	}
	return toPBListResult(result), nil
}

func (h *Handler) CreateArticle(ctx context.Context, req *kitexcontent.CreateArticleRequest) (*kitexcontent.GetArticleResponse, error) {
	created, err := h.articles.CreateArticle(ctx, toInternalArticle(req.GetArticle()))
	if err != nil {
		return nil, err
	}
	return &kitexcontent.GetArticleResponse{Article: toPBArticle(created)}, nil
}

func (h *Handler) UpdateArticle(ctx context.Context, req *kitexcontent.UpdateArticleRequest) (*kitexcontent.GetArticleResponse, error) {
	updated, err := h.articles.UpdateArticle(ctx, toInternalArticle(req.GetArticle()))
	if err != nil {
		return nil, err
	}
	return &kitexcontent.GetArticleResponse{Article: toPBArticle(updated)}, nil
}

func (h *Handler) DeleteArticle(ctx context.Context, req *kitexcontent.DeleteArticleRequest) (*kitexcontent.DeleteArticleResponse, error) {
	if err := h.articles.DeleteArticle(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &kitexcontent.DeleteArticleResponse{Success: true}, nil
}

func toPBListResult(result *model.ListResult) *kitexcontent.ListArticlesResponse {
	if result == nil {
		return &kitexcontent.ListArticlesResponse{}
	}
	items := make([]*kitexcontent.Article, 0, len(result.List))
	for _, article := range result.List {
		items = append(items, toPBArticle(article))
	}
	return &kitexcontent.ListArticlesResponse{List: items, Total: result.Total}
}

func toPBArticle(article *model.ArticleDetail) *kitexcontent.Article {
	if article == nil {
		return nil
	}
	tags := make([]*kitexcontent.Tag, 0, len(article.Tags))
	for _, tag := range article.Tags {
		tags = append(tags, &kitexcontent.Tag{Id: tag.ID, Name: tag.Name, Slug: tag.Slug})
	}
	return &kitexcontent.Article{
		Id:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Summary:     article.Summary,
		ContentMd:   article.ContentMd,
		ContentHtml: article.ContentHTML,
		CoverImage:  article.CoverImage,
		CategoryId:  article.CategoryID,
		TagIds:      append([]int64(nil), article.TagIDs...),
		Status:      article.Status,
		Tags:        tags,
	}
}

func toInternalArticle(article *kitexcontent.Article) *model.ArticleDetail {
	if article == nil {
		return nil
	}
	tags := make([]model.TagDTO, 0, len(article.GetTags()))
	for _, tag := range article.GetTags() {
		tags = append(tags, model.TagDTO{ID: tag.GetId(), Name: tag.GetName(), Slug: tag.GetSlug()})
	}
	return &model.ArticleDetail{
		ID:          article.GetId(),
		Title:       article.GetTitle(),
		Slug:        article.GetSlug(),
		Summary:     article.GetSummary(),
		ContentMd:   article.GetContentMd(),
		ContentHTML: article.GetContentHtml(),
		CoverImage:  article.GetCoverImage(),
		CategoryID:  article.GetCategoryId(),
		TagIDs:      append([]int64(nil), article.GetTagIds()...),
		Status:      article.GetStatus(),
		Tags:        tags,
	}
}
