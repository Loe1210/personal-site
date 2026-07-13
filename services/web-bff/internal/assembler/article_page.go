package assembler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ArticleDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Summary     string `json:"summary"`
	ContentHTML string `json:"content_html"`
	CoverImage  string `json:"cover_image"`
}

type ArticlePageDTO struct {
	Article *ArticleDTO `json:"article"`
}

type ContentClient interface {
	GetArticleByID(ctx context.Context, id int64) (*ArticleDTO, error)
}

type ArticlePageAssembler struct {
	content ContentClient
}

func NewArticlePageAssembler(content ContentClient) *ArticlePageAssembler {
	return &ArticlePageAssembler{content: content}
}

func (a *ArticlePageAssembler) BuildArticlePage(ctx context.Context, id int64) (*ArticlePageDTO, error) {
	article, err := a.content.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ArticlePageDTO{Article: article}, nil
}

type HTTPContentClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPContentClient(baseURL string) *HTTPContentClient {
	return &HTTPContentClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPContentClient) GetArticleByID(ctx context.Context, id int64) (*ArticleDTO, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/articles/%d", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("content-service returned %d", resp.StatusCode)
	}
	var envelope struct {
		Data ArticleDTO `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	return &envelope.Data, nil
}
