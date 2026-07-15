package model

type Tag struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	ArticleCount int64  `json:"article_count"`
}
