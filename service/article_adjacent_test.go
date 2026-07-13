package service

import (
	"testing"
	"time"

	articlemodel "github.com/Loe1210/personal-site/biz/model/article"
	"github.com/Loe1210/personal-site/dal/db"
)

func TestFindAdjacentArticlesReturnsNeighborsInPublicOrder(t *testing.T) {
	publishedAt := func(day int) *time.Time {
		value := time.Date(2026, 7, day, 10, 0, 0, 0, time.UTC)
		return &value
	}

	records := []db.Article{
		{ID: 11, Title: "Pinned", Slug: "pinned", IsTop: 1, PublishedAt: publishedAt(12), CreatedAt: time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)},
		{ID: 10, Title: "Newest", Slug: "newest", IsTop: 0, PublishedAt: publishedAt(11), CreatedAt: time.Date(2026, 7, 11, 10, 0, 0, 0, time.UTC)},
		{ID: 9, Title: "Current", Slug: "current", IsTop: 0, PublishedAt: publishedAt(10), CreatedAt: time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC)},
		{ID: 8, Title: "Older", Slug: "older", IsTop: 0, PublishedAt: publishedAt(9), CreatedAt: time.Date(2026, 7, 9, 10, 0, 0, 0, time.UTC)},
	}

	prev, next := findAdjacentArticles(records, 9)
	if prev == nil || prev.ID != 10 || prev.Title != "Newest" {
		t.Fatalf("expected previous article 10/Newest, got %#v", prev)
	}
	if next == nil || next.ID != 8 || next.Title != "Older" {
		t.Fatalf("expected next article 8/Older, got %#v", next)
	}
}

func TestFindAdjacentArticlesReturnsNilWhenCurrentMissing(t *testing.T) {
	records := []db.Article{{ID: 1, Title: "Only", Slug: "only"}}

	prev, next := findAdjacentArticles(records, 999)
	if prev != nil || next != nil {
		t.Fatalf("expected nil neighbors, got prev=%#v next=%#v", prev, next)
	}
}

var _ = articlemodel.Article{}
