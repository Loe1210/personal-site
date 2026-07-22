package content

import (
	"testing"

	kitexcontent "github.com/Loe1210/personal-site/kitex_gen/content"
)

func TestArticleFromPBMapsCoreFields(t *testing.T) {
	got := articleFromPB(&kitexcontent.Article{
		Id:          12,
		Title:       "Deep Refactor",
		Slug:        "deep-refactor",
		Summary:     "summary",
		ContentMd:   "markdown body",
		ContentHtml: "<p>html body</p>",
		CoverImage:  "/covers/a.png",
		Status:      "published",
		CategoryId:  3,
		TagIds:      []int64{5, 8},
		Tags: []*kitexcontent.Tag{
			{Id: 5, Name: "Go", Slug: "go"},
			{Id: 8, Name: "Microservices", Slug: "microservices"},
		},
	})

	if got.ID != 12 {
		t.Fatalf("expected id 12, got %d", got.ID)
	}
	if got.Title != "Deep Refactor" || got.Slug != "deep-refactor" || got.Summary != "summary" {
		t.Fatalf("unexpected text fields: %#v", got)
	}
	if got.ContentMD != "markdown body" || got.ContentHTML != "<p>html body</p>" || got.CoverImage != "/covers/a.png" {
		t.Fatalf("unexpected content fields: %#v", got)
	}
	if got.Status != "published" || got.CategoryID != 3 {
		t.Fatalf("unexpected metadata fields: %#v", got)
	}
	if len(got.TagIDs) != 2 || got.TagIDs[0] != 5 || got.TagIDs[1] != 8 {
		t.Fatalf("unexpected tag ids: %#v", got.TagIDs)
	}
	if len(got.Tags) != 2 || got.Tags[0].Name != "Go" || got.Tags[1].Slug != "microservices" {
		t.Fatalf("unexpected tags: %#v", got.Tags)
	}
}

func TestArticleFromPBNilIsZeroArticle(t *testing.T) {
	got := articleFromPB(nil)
	if got.ID != 0 || got.Title != "" || len(got.Tags) != 0 {
		t.Fatalf("expected zero article, got %#v", got)
	}
}
