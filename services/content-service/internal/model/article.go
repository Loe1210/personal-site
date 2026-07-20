package model

import "time"

type ArticleDetail struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	ContentMd   string     `json:"content_md"`
	ContentHTML string     `json:"content_html"`
	CoverImage  string     `json:"cover_image"`
	CategoryID  int64      `json:"category_id"`
	TagIDs      []int64    `json:"tag_ids"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Tags        []TagDTO   `json:"tags"`
}

type TagDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ListFilter struct {
	Page     int64
	PageSize int64
	Status   string
	Keyword  string
}

type ListResult struct {
	List  []*ArticleDetail `json:"list"`
	Total int64            `json:"total"`
}

type AdjacentArticles struct {
	Prev *ArticleDetail `json:"prev"`
	Next *ArticleDetail `json:"next"`
}
