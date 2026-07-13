package service

import (
	"context"
	"strings"
	"time"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	"github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"

	"gorm.io/gorm"
)

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ListPublicArticles(ctx context.Context, req *articlemodel.ListArticlesRequest) (*articlemodel.ListArticlesResponse, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return List(ctx, page, pageSize, "published", req.Category, req.Tag, req.Keyword, false)
}

func ListAdminArticles(ctx context.Context, req *articlemodel.ListArticlesRequest) (*articlemodel.ListArticlesResponse, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return List(ctx, page, pageSize, req.Status, req.Category, req.Tag, req.Keyword, false)
}

func List(_ context.Context, page, pageSize int64, status string, categoryName string, tagName string, keyword string, is_all bool) (*articlemodel.ListArticlesResponse, error) {
	var total int64
	var records []db.Article

	query := db.DB.Model(&db.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if categoryName != "" {
		var category db.Category
		if err := db.DB.Where("name = ? OR slug = ?", categoryName, categoryName).First(&category).Error; err == nil {
			query = query.Where("category_id = ?", category.ID)
		} else {
			return &articlemodel.ListArticlesResponse{
				List:     []*articlemodel.Article{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
			}, nil
		}
	}
	if tagName != "" {
		var tag db.Tag
		if err := db.DB.Where("name = ? OR slug = ?", tagName, tagName).First(&tag).Error; err == nil {
			query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").Where("article_tags.tag_id = ?", tag.ID)
		} else {
			return &articlemodel.ListArticlesResponse{
				List:     []*articlemodel.Article{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
			}, nil
		}
	}
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ?", keyword, keyword)
	}
	query = query.Distinct()

	if err := query.Count(&total).Error; err != nil {
		return nil, errno.Internal
	}

	if is_all {
		if err := query.Order("is_top DESC, published_at DESC, created_at DESC").Find(&records).Error; err != nil {
			return nil, errno.Internal
		}
	} else {
		offset := (page - 1) * pageSize
		if err := query.Order("is_top DESC, published_at DESC, created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&records).Error; err != nil {
			return nil, errno.Internal
		}
	}

	tagsByArticleID, err := getTagsForArticles(records)
	if err != nil {
		return nil, err
	}

	items := make([]*articlemodel.Article, 0, len(records))
	for _, record := range records {
		item := toArticleModel(record, tagsByArticleID[record.ID])
		if item != nil {
			items = append(items, item)
		}
	}

	return &articlemodel.ListArticlesResponse{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func getTagsForArticles(articles []db.Article) (map[int64][]*tagmodel.Tag, error) {
	result := make(map[int64][]*tagmodel.Tag)
	if len(articles) == 0 {
		return result, nil
	}

	articleIDs := make([]int64, len(articles))
	for i, a := range articles {
		articleIDs[i] = a.ID
	}

	var articleTags []db.ArticleTag
	if err := db.DB.Where("article_id IN ?", articleIDs).Find(&articleTags).Error; err != nil {
		return nil, errno.Internal
	}

	tagIDs := make([]int64, 0, len(articleTags))
	tagMapByArticleID := make(map[int64][]int64)
	for _, at := range articleTags {
		tagIDs = append(tagIDs, at.TagID)
		tagMapByArticleID[at.ArticleID] = append(tagMapByArticleID[at.ArticleID], at.TagID)
	}

	tagResp, err := ListTags(context.Background(), &tagmodel.ListTagsRequest{})
	if err != nil {
		return nil, err
	}
	tagMap := make(map[int64]*tagmodel.Tag)
	for _, t := range tagResp.List {
		tagMap[t.ID] = t
	}

	for articleID, tids := range tagMapByArticleID {
		for _, tid := range tids {
			if t, ok := tagMap[tid]; ok {
				result[articleID] = append(result[articleID], t)
			}
		}
	}

	return result, nil
}

func GetPublicArticleByID(_ context.Context, id int64) (*articlemodel.GetArticleResponse, error) {
	var record db.Article
	if err := db.DB.Where("id = ? AND status = 'published'", id).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ArticleNotFound
		}
		return nil, errno.Internal
	}

	var tagRecords []db.Tag
	if err := db.DB.Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", record.ID).
		Find(&tagRecords).Error; err != nil {
		return nil, errno.Internal
	}
	tags := make([]*tagmodel.Tag, 0, len(tagRecords))
	for _, t := range tagRecords {
		tags = append(tags, toTagModel(&t))
	}

	return &articlemodel.GetArticleResponse{
		Article: toArticleDetailModel(record, tags),
	}, nil
}

func GetPublicArticleBySlug(_ context.Context, req *articlemodel.GetArticleBySlugRequest) (*articlemodel.GetArticleResponse, error) {
	var record db.Article
	if err := db.DB.Where("slug = ? AND status = 'published'", req.Slug).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ArticleNotFound
		}
		return nil, errno.Internal
	}

	var tagRecords []db.Tag
	if err := db.DB.Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", record.ID).
		Find(&tagRecords).Error; err != nil {
		return nil, errno.Internal
	}
	tags := make([]*tagmodel.Tag, 0, len(tagRecords))
	for _, t := range tagRecords {
		tags = append(tags, toTagModel(&t))
	}

	return &articlemodel.GetArticleResponse{
		Article: toArticleDetailModel(record, tags),
	}, nil
}

func Get(_ context.Context, id uint) (*articlemodel.GetArticleResponse, error) {
	var record db.Article
	if err := db.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ArticleNotFound
		}
		return nil, errno.Internal
	}

	var tagRecords []db.Tag
	if err := db.DB.Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", record.ID).
		Find(&tagRecords).Error; err != nil {
		return nil, errno.Internal
	}
	tags := make([]*tagmodel.Tag, 0, len(tagRecords))
	for _, t := range tagRecords {
		tags = append(tags, toTagModel(&t))
	}

	return &articlemodel.GetArticleResponse{
		Article: toArticleDetailModel(record, tags),
	}, nil
}

func GetAdjacentPublicArticles(_ context.Context, id int64) (*articlemodel.GetAdjacentArticlesResponse, error) {
	var current db.Article
	if err := db.DB.Where("id = ? AND status = 'published'", id).First(&current).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ArticleNotFound
		}
		return nil, errno.Internal
	}

	var records []db.Article
	if err := db.DB.Where("status = ?", "published").Order("is_top DESC, published_at DESC, created_at DESC").Find(&records).Error; err != nil {
		return nil, errno.Internal
	}

	prev, next := findAdjacentArticles(records, id)
	return &articlemodel.GetAdjacentArticlesResponse{
		Prev: prev,
		Next: next,
	}, nil
}
func CreateArticle(_ context.Context, req *articlemodel.CreateArticleRequest) (*articlemodel.CreateArticleResponse, error) {
	var count int64
	if err := db.DB.Model(&db.Article{}).Where("slug = ?", req.Slug).Count(&count).Error; err != nil {
		return nil, errno.Internal
	}
	if count > 0 {
		return nil, errno.SlugConflict
	}

	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	now := time.Now()
	record := &db.Article{
		Title:       req.Title,
		Slug:        slug,
		Summary:     req.Summary,
		ContentMd:   req.ContentMd,
		ContentHTML: "",
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
		IsTop:       0,
		AuthorID:    0,
	}

	if req.Status == "published" {
		record.PublishedAt = &now
	}

	if err := db.DB.Create(record).Error; err != nil {
		return nil, errno.Internal
	}

	if err := syncArticleTags(record.ID, req.TagIds); err != nil {
		return nil, err
	}

	resp, err := Get(context.Background(), uint(record.ID))
	if err != nil {
		return nil, err
	}
	return &articlemodel.CreateArticleResponse{Article: resp.Article, Message: "创建成功"}, nil
}

func UpdateArticle(_ context.Context, req *articlemodel.UpdateArticleRequest) (*articlemodel.UpdateArticleResponse, error) {
	id := uint(req.ID)
	var record db.Article
	if err := db.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ArticleNotFound
		}
		return nil, errno.Internal
	}

	if req.Slug != "" && req.Slug != record.Slug {
		var count int64
		if err := db.DB.Model(&db.Article{}).Where("slug = ? AND id <> ?", req.Slug, id).Count(&count).Error; err != nil {
			return nil, errno.Internal
		}
		if count > 0 {
			return nil, errno.SlugConflict
		}
		record.Slug = req.Slug
	}

	if req.Title != "" {
		record.Title = req.Title
	}
	if req.Summary != "" {
		record.Summary = req.Summary
	}
	if req.ContentMd != "" {
		record.ContentMd = req.ContentMd
		record.ContentHTML = ""
	}
	if req.CoverImage != "" {
		record.CoverImage = req.CoverImage
	}
	if req.CategoryID > 0 {
		record.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		if record.Status != "published" && req.Status == "published" {
			now := time.Now()
			record.PublishedAt = &now
		}
		record.Status = req.Status
	}
	record.IsTop = 0

	if err := db.DB.Save(&record).Error; err != nil {
		return nil, errno.Internal
	}

	if err := syncArticleTags(record.ID, req.TagIds); err != nil {
		return nil, err
	}

	resp, err := Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &articlemodel.UpdateArticleResponse{Article: resp.Article, Message: "更新成功"}, nil
}

func DeleteArticle(_ context.Context, req *articlemodel.DeleteArticleRequest) (*articlemodel.DeleteArticleResponse, error) {
	id := uint(req.ID)
	var record db.Article
	if err := db.DB.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ArticleNotFound
		}
		return nil, errno.Internal
	}

	if err := db.DB.Where("article_id = ?", id).Delete(&db.ArticleTag{}).Error; err != nil {
		return nil, errno.Internal
	}

	if err := db.DB.Delete(&record).Error; err != nil {
		return nil, errno.Internal
	}

	return &articlemodel.DeleteArticleResponse{Success: true, Message: "删除成功"}, nil
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func syncArticleTags(articleID int64, tagIDs []int64) error {
	if err := db.DB.Where("article_id = ?", articleID).Delete(&db.ArticleTag{}).Error; err != nil {
		return errno.Internal
	}

	if len(tagIDs) == 0 {
		return nil
	}

	var articleTags []db.ArticleTag
	for _, tagID := range tagIDs {
		if tagID > 0 {
			articleTags = append(articleTags, db.ArticleTag{
				ArticleID: articleID,
				TagID:     tagID,
			})
		}
	}

	if len(articleTags) > 0 {
		if err := db.DB.Create(&articleTags).Error; err != nil {
			return errno.Internal
		}
	}

	return nil
}

func findAdjacentArticles(records []db.Article, currentID int64) (*articlemodel.AdjacentArticle, *articlemodel.AdjacentArticle) {
	currentIndex := -1
	for i, record := range records {
		if record.ID == currentID {
			currentIndex = i
			break
		}
	}
	if currentIndex == -1 {
		return nil, nil
	}

	var prev *articlemodel.AdjacentArticle
	var next *articlemodel.AdjacentArticle
	if currentIndex > 0 {
		prev = toAdjacentArticle(records[currentIndex-1])
	}
	if currentIndex < len(records)-1 {
		next = toAdjacentArticle(records[currentIndex+1])
	}
	return prev, next
}

func toAdjacentArticle(record db.Article) *articlemodel.AdjacentArticle {
	publishedAt := ""
	if record.PublishedAt != nil {
		publishedAt = formatTime(*record.PublishedAt)
	}
	return &articlemodel.AdjacentArticle{
		ID:          record.ID,
		Title:       record.Title,
		Slug:        record.Slug,
		PublishedAt: publishedAt,
	}
}
func toArticleModel(record db.Article, tags []*tagmodel.Tag) *articlemodel.Article {
	publishedAt := ""
	if record.PublishedAt != nil {
		publishedAt = formatTime(*record.PublishedAt)
	}
	createdAt := formatTime(record.CreatedAt)

	article := &articlemodel.Article{
		ID:          record.ID,
		Title:       record.Title,
		Slug:        record.Slug,
		Summary:     record.Summary,
		CoverImage:  record.CoverImage,
		CategoryID:  record.CategoryID,
		Status:      record.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   formatTime(record.UpdatedAt),
		PublishedAt: publishedAt,
	}

	if tags != nil {
		article.TagIds = make([]int64, 0, len(tags))
		for _, t := range tags {
			article.TagIds = append(article.TagIds, t.ID)
		}
	}

	return article
}

func toArticleDetailModel(record db.Article, tags []*tagmodel.Tag) *articlemodel.Article {
	article := toArticleModel(record, tags)
	article.ContentMd = record.ContentMd
	article.ContentHTML = record.ContentHTML
	return article
}
