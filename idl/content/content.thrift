namespace go content

struct Tag {
    1: i64 id
    2: string name
    3: string slug
}

struct Category {
    1: i64 id
    2: string name
    3: string slug
    4: string description
}

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
    11: list<Tag> tags
}

struct GetArticleByIDRequest {
    1: i64 id
}

struct ListPublicArticlesRequest {
    1: i64 page
    2: i64 page_size
    3: string keyword
}

struct ListArticlesResponse {
    1: list<Article> list
    2: i64 total
}

struct GetArticleResponse {
    1: Article article
}

struct CreateArticleRequest {
    1: Article article
}

struct UpdateArticleRequest {
    1: Article article
}

struct DeleteArticleRequest {
    1: i64 id
}

struct DeleteArticleResponse {
    1: bool success
}

service ContentService {
    GetArticleResponse GetArticleByID(1: GetArticleByIDRequest req)
    ListArticlesResponse ListPublicArticles(1: ListPublicArticlesRequest req)
    GetArticleResponse CreateArticle(1: CreateArticleRequest req)
    GetArticleResponse UpdateArticle(1: UpdateArticleRequest req)
    DeleteArticleResponse DeleteArticle(1: DeleteArticleRequest req)
}