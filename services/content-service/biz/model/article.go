package model

type ArticleRequest struct {
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Summary     string  `json:"summary"`
	ContentMd   string  `json:"content_md"`
	ContentHTML string  `json:"content_html"`
	CoverImage  string  `json:"cover_image"`
	CategoryID  int64   `json:"category_id"`
	TagIDs      []int64 `json:"tag_ids"`
	Status      string  `json:"status"`
}
