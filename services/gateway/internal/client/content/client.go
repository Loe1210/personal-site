package content

import "context"

type ListPublicArticlesRequest struct {
	Page     int64
	PageSize int64
	Keyword  string
}

type Tag struct {
	ID   int64
	Name string
	Slug string
}

type Article struct {
	ID          int64
	Title       string
	Slug        string
	Summary     string
	ContentMD   string
	ContentHTML string
	CoverImage  string
	CategoryID  int64
	TagIDs      []int64
	Status      string
	Tags        []Tag
}

type ListPublicArticlesResponse struct {
	List  []Article
	Total int64
}

type ArticleClient interface {
	ListPublicArticles(ctx context.Context, req ListPublicArticlesRequest) (*ListPublicArticlesResponse, error)
	GetArticleByID(ctx context.Context, id int64) (*Article, error)
}
