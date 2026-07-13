package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Loe1210/personal-site/services/content-service/internal/application"
)

type Article struct {
	ID          int64      `gorm:"primaryKey;autoIncrement"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Summary     string     `gorm:"type:text"`
	ContentMd   string     `gorm:"column:content_md;type:longtext"`
	ContentHTML string     `gorm:"column:content_html;type:longtext"`
	CoverImage  string     `gorm:"type:varchar(255)"`
	CategoryID  int64      `gorm:"default:0"`
	Status      string     `gorm:"type:varchar(32);index;default:'draft'"`
	IsTop       int        `gorm:"type:tinyint(1);default:0;index"`
	AuthorID    int64      `gorm:"default:0"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	PublishedAt *time.Time `gorm:"column:published_at;default:null"`
}

type ArticleTag struct {
	ArticleID int64 `gorm:"primaryKey"`
	TagID     int64 `gorm:"primaryKey"`
}

type Tag struct {
	ID   int64
	Name string
	Slug string
}

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) GetByID(ctx context.Context, id int64) (*application.ArticleDetail, error) {
	var article Article
	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	tags, err := r.tagsForArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	return toArticleDetail(article, tags), nil
}

func (r *ArticleRepository) List(ctx context.Context, filter application.ListFilter) (*application.ListResult, error) {
	var articles []Article
	query := r.db.WithContext(ctx).Model(&Article{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order("is_top DESC, published_at DESC, id DESC").Offset(int(offset)).Limit(int(filter.PageSize)).Find(&articles).Error; err != nil {
		return nil, err
	}
	result := make([]*application.ArticleDetail, 0, len(articles))
	for _, article := range articles {
		result = append(result, toArticleDetail(article, nil))
	}
	return &application.ListResult{List: result, Total: total}, nil
}

func (r *ArticleRepository) Create(ctx context.Context, detail *application.ArticleDetail) error {
	article := fromArticleDetail(detail)
	if err := r.db.WithContext(ctx).Create(&article).Error; err != nil {
		return err
	}
	detail.ID = article.ID
	return r.syncTags(ctx, article.ID, detail.TagIDs)
}

func (r *ArticleRepository) Update(ctx context.Context, detail *application.ArticleDetail) error {
	article := fromArticleDetail(detail)
	if err := r.db.WithContext(ctx).Model(&Article{}).Where("id = ?", detail.ID).Updates(article).Error; err != nil {
		return err
	}
	return r.syncTags(ctx, detail.ID, detail.TagIDs)
}

func (r *ArticleRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Where("article_id = ?", id).Delete(&ArticleTag{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&Article{}, id).Error
}

func (r *ArticleRepository) tagsForArticle(ctx context.Context, articleID int64) ([]application.TagDTO, error) {
	var tags []Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	result := make([]application.TagDTO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, application.TagDTO{ID: tag.ID, Name: tag.Name, Slug: tag.Slug})
	}
	return result, nil
}

func (r *ArticleRepository) syncTags(ctx context.Context, articleID int64, tagIDs []int64) error {
	if err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&ArticleTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	rows := make([]ArticleTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID > 0 {
			rows = append(rows, ArticleTag{ArticleID: articleID, TagID: tagID})
		}
	}
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&rows).Error
}

func toArticleDetail(article Article, tags []application.TagDTO) *application.ArticleDetail {
	tagIDs := make([]int64, 0, len(tags))
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
	}
	return &application.ArticleDetail{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Summary:     article.Summary,
		ContentMd:   article.ContentMd,
		ContentHTML: article.ContentHTML,
		CoverImage:  article.CoverImage,
		CategoryID:  article.CategoryID,
		TagIDs:      tagIDs,
		Status:      article.Status,
		Tags:        tags,
	}
}

func fromArticleDetail(detail *application.ArticleDetail) Article {
	return Article{
		ID:          detail.ID,
		Title:       detail.Title,
		Slug:        detail.Slug,
		Summary:     detail.Summary,
		ContentMd:   detail.ContentMd,
		ContentHTML: detail.ContentHTML,
		CoverImage:  detail.CoverImage,
		CategoryID:  detail.CategoryID,
		Status:      detail.Status,
	}
}
