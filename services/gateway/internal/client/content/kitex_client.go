package content

import (
	"context"

	kitexcontent "github.com/Loe1210/personal-site/kitex_gen/content"
	"github.com/Loe1210/personal-site/kitex_gen/content/contentservice"
)

type KitexArticleClient struct {
	cli contentservice.Client
}

func NewKitexArticleClient(cli contentservice.Client) *KitexArticleClient {
	return &KitexArticleClient{cli: cli}
}

func (c *KitexArticleClient) ListPublicArticles(ctx context.Context, req ListPublicArticlesRequest) (*ListPublicArticlesResponse, error) {
	resp, err := c.cli.ListPublicArticles(ctx, &kitexcontent.ListPublicArticlesRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}

	items := make([]Article, 0, len(resp.GetList()))
	for _, item := range resp.GetList() {
		items = append(items, articleFromPB(item))
	}

	return &ListPublicArticlesResponse{
		List:  items,
		Total: resp.GetTotal(),
	}, nil
}

func (c *KitexArticleClient) GetArticleByID(ctx context.Context, id int64) (*Article, error) {
	resp, err := c.cli.GetArticleByID(ctx, &kitexcontent.GetArticleByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}

	article := articleFromPB(resp.GetArticle())
	return &article, nil
}

func articleFromPB(article *kitexcontent.Article) Article {
	if article == nil {
		return Article{}
	}

	tags := make([]Tag, 0, len(article.GetTags()))
	for _, tag := range article.GetTags() {
		tags = append(tags, Tag{
			ID:   tag.GetId(),
			Name: tag.GetName(),
			Slug: tag.GetSlug(),
		})
	}

	return Article{
		ID:          article.GetId(),
		Title:       article.GetTitle(),
		Slug:        article.GetSlug(),
		Summary:     article.GetSummary(),
		ContentMD:   article.GetContentMd(),
		ContentHTML: article.GetContentHtml(),
		CoverImage:  article.GetCoverImage(),
		CategoryID:  article.GetCategoryId(),
		TagIDs:      append([]int64(nil), article.GetTagIds()...),
		Status:      article.GetStatus(),
		Tags:        tags,
	}
}
