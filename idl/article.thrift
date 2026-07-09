namespace go article

struct Article {
    1: i64 id
    2: string title
    3: string slug
    4: string summary
    5: string content_md
    6: string content_html
    7: string cover_image
    8: i64 category_id
    9: list<i64> tag_ids
    10: string status
    11: string created_at
    12: string updated_at
    13: string published_at
}

struct CreateArticleRequest {
    1: string title (api.body="title")
    2: string slug (api.body="slug")
    3: string summary (api.body="summary")
    4: string content_md (api.body="content_md")
    5: string cover_image (api.body="cover_image")
    6: i64 category_id (api.body="category_id")
    7: list<i64> tag_ids (api.body="tag_ids")
    8: string status (api.body="status")
}

struct CreateArticleResponse {
    1: Article article
    2: string message
}

struct UpdateArticleRequest {
    1: i64 id (api.path="id")
    2: string title (api.body="title")
    3: string slug (api.body="slug")
    4: string summary (api.body="summary")
    5: string content_md (api.body="content_md")
    6: string cover_image (api.body="cover_image")
    7: i64 category_id (api.body="category_id")
    8: list<i64> tag_ids (api.body="tag_ids")
    9: string status (api.body="status")
}

struct UpdateArticleResponse {
    1: Article article
    2: string message
}

struct DeleteArticleRequest {
    1: i64 id (api.path="id")
}

struct DeleteArticleResponse {
    1: bool success
    2: string message
}

struct GetArticleByIDRequest {
    1: i64 id (api.path="id")
}

struct GetArticleBySlugRequest {
    1: string slug (api.path="slug")
}

struct GetArticleResponse {
    1: Article article
}

struct ListArticlesRequest {
    1: i64 page (api.query="page")
    2: i64 page_size (api.query="page_size")
    3: string tag (api.query="tag")
    4: string category (api.query="category")
    5: string keyword (api.query="keyword")
    6: string status (api.query="status")
}

struct ListArticlesResponse {
    1: list<Article> list
    2: i64 total
    3: i64 page
    4: i64 page_size
}

struct PublishArticleRequest {
    1: i64 id (api.path="id")
    2: string status (api.body="status")
}

struct PublishArticleResponse {
    1: Article article
    2: string message
}

service ArticleService {
    CreateArticleResponse CreateArticle(1: CreateArticleRequest req) (api.post="/api/admin/articles")
    UpdateArticleResponse UpdateArticle(1: UpdateArticleRequest req) (api.put="/api/admin/articles/:id")
    DeleteArticleResponse DeleteArticle(1: DeleteArticleRequest req) (api.delete="/api/admin/articles/:id")
    GetArticleResponse GetArticleByID(1: GetArticleByIDRequest req) (api.get="/api/admin/articles/:id")
    GetArticleResponse GetArticleBySlug(1: GetArticleBySlugRequest req) (api.get="/api/articles/:slug")
    ListArticlesResponse ListArticles(1: ListArticlesRequest req) (api.get="/api/articles")
    ListArticlesResponse ListAdminArticles(1: ListArticlesRequest req) (api.get="/api/admin/articles")
    PublishArticleResponse PublishArticle(1: PublishArticleRequest req) (api.patch="/api/admin/articles/:id/publish")
}