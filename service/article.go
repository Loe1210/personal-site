package service

import (
	"context"
	"strings"
	"time"
	"gorm.io/gorm"
	mysqlDriver "github.com/go-sql-driver/mysql"

	dbmodel "github.com/Loe1210/personal-site/dal/db"
	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
	"github.com/Loe1210/personal-site/pkg/errno"
	
)

const timeLayout = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(timeLayout)
}

func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Local().Format(timeLayout)
}

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
		list = append(list, toArticleModel(dbmodel.DB, &records[i]))
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
		Article: toArticleModel(dbmodel.DB, &record),
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
		list = append(list, toArticleModel(dbmodel.DB, &records[i]))
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
	status := req.Status
	if status == "" {
		status = "draft"
	}

	if err := ensureCategoryExists(dbmodel.DB, req.CategoryID); err != nil {
		return nil, err
	}
	if err := ensureTagsExist(dbmodel.DB, req.TagIds); err != nil {
		return nil, err
	}

	record := &dbmodel.Article{
		Title:       req.Title,
		Slug:        req.Slug,
		Summary:     req.Summary,
		ContentMd:   req.ContentMd,
		ContentHTML: req.ContentMd,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		Status:      status,
	}

	if status == "published" {
		now := time.Now()
		record.PublishedAt = &now
	}

	if err := dbmodel.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
				return errno.SlugConflict
			}
			return errno.Internal
		}

		if err := replaceArticleTags(tx, record.ID, req.TagIds); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &articlemodel.CreateArticleResponse{
		Article: toArticleModel(dbmodel.DB, record),
		Message: "article created",
	}, nil
}

func UpdateArticle(_ context.Context, req *articlemodel.UpdateArticleRequest) (*articlemodel.UpdateArticleResponse, error) {
	var record dbmodel.Article

	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, nil
	}
	if err := ensureCategoryExists(dbmodel.DB, req.CategoryID); err != nil {
		return nil, err
	}
	if err := ensureTagsExist(dbmodel.DB, req.TagIds); err != nil {
		return nil, err
	}

	record.Title = req.Title
	record.Slug = req.Slug
	record.Summary = req.Summary
	record.ContentMd = req.ContentMd
	record.ContentHTML = req.ContentMd
	record.CoverImage = req.CoverImage
	record.CategoryID = req.CategoryID
	record.Status = req.Status

	if record.Status == "published" {
		if record.PublishedAt == nil {
			now := time.Now()
			record.PublishedAt = &now
		}
	} else {
		record.PublishedAt = nil
	}


	if err := dbmodel.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&record).Error; err != nil {
			if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
				return errno.SlugConflict
			}
			return errno.Internal
		}

		if err := replaceArticleTags(tx, record.ID, req.TagIds); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &articlemodel.UpdateArticleResponse{
		Article: toArticleModel(dbmodel.DB, &record),
		Message: "article updated",
	}, nil
}

func DeleteArticle(_ context.Context, req *articlemodel.DeleteArticleRequest) (*articlemodel.DeleteArticleResponse, error) {
	var rowsAffected int64

	if err := dbmodel.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", req.ID).Delete(&dbmodel.ArticleTag{}).Error; err != nil {
			return errno.Internal
		}

		result := tx.Delete(&dbmodel.Article{}, req.ID)
		if result.Error != nil {
			return errno.Internal
		}

		rowsAffected = result.RowsAffected
		return nil
	}); err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
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



func toArticleModel(tx *gorm.DB, item *dbmodel.Article) *articlemodel.Article {
	if item == nil {
		return nil
	}

	tagIDs, err := getArticleTagIDs(tx, item.ID)
	if err != nil {
		tagIDs = []int64{}
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
		TagIds:      tagIDs,
		Status:      item.Status,
		CreatedAt:   formatTime(item.CreatedAt),
		UpdatedAt:   formatTime(item.UpdatedAt),
		PublishedAt: formatTimePtr(item.PublishedAt),
	}
}



func ensureCategoryExists(tx *gorm.DB, categoryID int64) error {
	if categoryID == 0 {
		return nil
	}

	var count int64
	if err := tx.Model(&dbmodel.Category{}).
		Where("id = ?", categoryID).
		Count(&count).Error; err != nil {
		return errno.Internal
	}

	if count == 0 {
		return errno.CategoryNotFound
	}

	return nil
}

func ensureTagsExist(tx *gorm.DB, tagIDs []int64) error {
	tagIDs = uniqueInt64(tagIDs)
	if len(tagIDs) == 0 {
		return nil
	}

	var count int64
	if err := tx.Model(&dbmodel.Tag{}).
		Where("id IN ?", tagIDs).
		Count(&count).Error; err != nil {
		return errno.Internal
	}

	if count != int64(len(tagIDs)) {
		return errno.TagNotFound
	}

	return nil
}

func uniqueInt64(ids []int64) []int64 {
	if len(ids) == 0 {
		return ids
	}

	seen := make(map[int64]struct{}, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}

func replaceArticleTags(tx *gorm.DB, articleID int64, tagIDs []int64) error {
	tagIDs = uniqueInt64(tagIDs)

	if err := tx.Where("article_id = ?", articleID).
		Delete(&dbmodel.ArticleTag{}).Error; err != nil {
		return errno.Internal
	}

	if len(tagIDs) == 0 {
		return nil
	}

	relations := make([]dbmodel.ArticleTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		relations = append(relations, dbmodel.ArticleTag{
			ArticleID: articleID,
			TagID:     tagID,
		})
	}

	if err := tx.Create(&relations).Error; err != nil {
		return errno.Internal
	}

	return nil
}

func getArticleTagIDs(tx *gorm.DB, articleID int64) ([]int64, error) {
	var relations []dbmodel.ArticleTag
	if err := tx.Where("article_id = ?", articleID).
		Find(&relations).Error; err != nil {
		return nil, errno.Internal
	}

	tagIDs := make([]int64, 0, len(relations))
	for _, item := range relations {
		tagIDs = append(tagIDs, item.TagID)
	}

	return tagIDs, nil
}
