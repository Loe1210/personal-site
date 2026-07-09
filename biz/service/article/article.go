package article

import (
	"context"
	"strconv"
	"strings"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
)



func ListPublicArticles(_ context.Context, req *articlemodel.ListArticlesRequest) (*articlemodel.ListArticlesResponse, error) {
	var records []dbmodel.Article

	query := dbmodel.DB.Where("status = ?", "published")
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	list := make([]*articlemodel.Article, 0, len(records))
	for i := range records {
		list = append(list, toArticleModel(&records[i]))
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
	var record dbmodel.Article

	err := dbmodel.DB.Where("slug = ? AND status = ?", req.Slug, "published").First(&record).Error
	if err != nil {
		return nil, nil
	}

	return &articlemodel.GetArticleResponse{
		Article: toArticleModel(&record),
	}, nil
}

func ListAdminArticles(_ context.Context, req *articlemodel.ListArticlesRequest) (*articlemodel.ListArticlesResponse, error) {
	var records []dbmodel.Article

	query := dbmodel.DB.Model(&dbmodel.Article{})

	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if strings.TrimSpace(req.Status) != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	list := make([]*articlemodel.Article, 0, len(records))
	for i := range records {
		list = append(list, toArticleModel(&records[i]))
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
	now := "2026-07-09 14:40:00"
	status := req.Status
	if status == "" {
		status = "draft"
	}

	record := &dbmodel.Article{
		Title:       req.Title,
		Slug:        req.Slug,
		Summary:     req.Summary,
		ContentMd:   req.ContentMd,
		ContentHTML: req.ContentMd,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		TagIds:      joinTagIDs(req.TagIds),
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: "",
	}

	if status == "published" {
		record.PublishedAt = now
	}

	if err := dbmodel.DB.Create(record).Error; err != nil {
		return nil, err
	}

	return &articlemodel.CreateArticleResponse{
		Article: toArticleModel(record),
		Message: "article created",
	}, nil
}

func UpdateArticle(_ context.Context, req *articlemodel.UpdateArticleRequest) (*articlemodel.UpdateArticleResponse, error) {
	var record dbmodel.Article

	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, nil
	}

	record.Title = req.Title
	record.Slug = req.Slug
	record.Summary = req.Summary
	record.ContentMd = req.ContentMd
	record.ContentHTML = req.ContentMd
	record.CoverImage = req.CoverImage
	record.CategoryID = req.CategoryID
	record.TagIds = joinTagIDs(req.TagIds)
	record.Status = req.Status
	record.UpdatedAt = "2026-07-09 15:00:00"

	if record.Status == "published" {
		if record.PublishedAt == "" {
			record.PublishedAt = "2026-07-09 15:00:00"
		}
	} else {
		record.PublishedAt = ""
	}

	if err := dbmodel.DB.Save(&record).Error; err != nil {
		return nil, err
	}

	return &articlemodel.UpdateArticleResponse{
		Article: toArticleModel(&record),
		Message: "article updated",
	}, nil
}

func DeleteArticle(_ context.Context, req *articlemodel.DeleteArticleRequest) (*articlemodel.DeleteArticleResponse, error) {
	result := dbmodel.DB.Delete(&dbmodel.Article{}, req.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return &articlemodel.DeleteArticleResponse{
			Success: false,
			Message: "article not found",
		}, nil
	}

	return &articlemodel.DeleteArticleResponse{
		Success: true,
		Message: "article deleted",
	}, nil
}

func parseTagIDs(tagIDs string) []int64 {
	if strings.TrimSpace(tagIDs) == "" {
		return []int64{}
	}

	parts := strings.Split(tagIDs, ",")
	result := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			continue
		}
		result = append(result, v)
	}
	return result
}

func toArticleModel(item *dbmodel.Article) *articlemodel.Article {
	if item == nil {
		return nil
	}

	return &articlemodel.Article{
		ID:          item.ID,
		Title:       item.Title,
		Slug:        item.Slug,
		Summary:     item.Summary,
		ContentMd:   item.ContentMd,
		ContentHTML: item.ContentHTML,
		CoverImage:  item.CoverImage,
		CategoryID:  item.CategoryID,
		TagIds:      parseTagIDs(item.TagIds),
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		PublishedAt: item.PublishedAt,
	}
}

func joinTagIDs(tagIDs []int64) string {
	if len(tagIDs) == 0 {
		return ""
	}

	parts := make([]string, 0, len(tagIDs))
	for _, id := range tagIDs {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
}