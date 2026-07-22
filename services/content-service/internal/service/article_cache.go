package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type ArticleCache interface {
	GetArticle(ctx context.Context, id int64) (*ArticleDetail, bool, error)
	SetArticle(ctx context.Context, article *ArticleDetail) error
	DeleteArticle(ctx context.Context, id int64) error
}

type LocalArticleCache struct {
	mu    sync.RWMutex
	items map[int64]*ArticleDetail
}

func NewLocalArticleCache() *LocalArticleCache {
	return &LocalArticleCache{items: map[int64]*ArticleDetail{}}
}

func (c *LocalArticleCache) GetArticle(_ context.Context, id int64) (*ArticleDetail, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	article, ok := c.items[id]
	return cloneArticle(article), ok, nil
}

func (c *LocalArticleCache) SetArticle(_ context.Context, article *ArticleDetail) error {
	if article == nil || article.ID <= 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[article.ID] = cloneArticle(article)
	return nil
}

func (c *LocalArticleCache) DeleteArticle(_ context.Context, id int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
	return nil
}

type RedisArticleCache struct {
	pool   *redis.Pool
	prefix string
	ttl    time.Duration
}

func NewRedisArticleCache(pool *redis.Pool, prefix string, ttl time.Duration) *RedisArticleCache {
	if prefix == "" {
		prefix = "content:article:"
	}
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &RedisArticleCache{pool: pool, prefix: prefix, ttl: ttl}
}

func (c *RedisArticleCache) GetArticle(_ context.Context, id int64) (*ArticleDetail, bool, error) {
	if c == nil || c.pool == nil {
		return nil, false, nil
	}
	conn := c.pool.Get()
	defer conn.Close()
	payload, err := redis.Bytes(conn.Do("GET", c.key(id)))
	if err == redis.ErrNil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var article ArticleDetail
	if err := json.Unmarshal(payload, &article); err != nil {
		return nil, false, err
	}
	return &article, true, nil
}

func (c *RedisArticleCache) SetArticle(_ context.Context, article *ArticleDetail) error {
	if c == nil || c.pool == nil || article == nil || article.ID <= 0 {
		return nil
	}
	payload, err := json.Marshal(article)
	if err != nil {
		return err
	}
	conn := c.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", c.key(article.ID), int(c.ttl.Seconds()), payload)
	return err
}

func (c *RedisArticleCache) DeleteArticle(_ context.Context, id int64) error {
	if c == nil || c.pool == nil {
		return nil
	}
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", c.key(id))
	return err
}

func (c *RedisArticleCache) key(id int64) string {
	return fmt.Sprintf("%s%d", c.prefix, id)
}

func cloneArticle(article *ArticleDetail) *ArticleDetail {
	if article == nil {
		return nil
	}
	clone := *article
	clone.TagIDs = append([]int64(nil), article.TagIDs...)
	clone.Tags = append([]TagDTO(nil), article.Tags...)
	return &clone
}
