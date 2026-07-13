package article

type AdjacentArticle struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	PublishedAt string `json:"published_at"`
}

type GetAdjacentArticlesResponse struct {
	Prev *AdjacentArticle `json:"prev"`
	Next *AdjacentArticle `json:"next"`
}
