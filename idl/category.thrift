namespace go category

struct Category {
    1: i64 id
    2: string name
    3: string slug
    4: string description
    5: string created_at
    6: string updated_at
}

struct CreateCategoryRequest {
    1: string name (api.body="name")
    2: string slug (api.body="slug")
    3: string description (api.body="description")
}

struct CreateCategoryResponse {
    1: Category category
    2: string message
}

struct ListCategoriesRequest {
}

struct ListCategoriesResponse {
    1: list<Category> list
}

service CategoryService {
    CreateCategoryResponse CreateCategory(1: CreateCategoryRequest req) (api.post="/api/admin/categories")
    ListCategoriesResponse ListCategories(1: ListCategoriesRequest req) (api.get="/api/categories")
    ListCategoriesResponse ListAdminCategories(1: ListCategoriesRequest req) (api.get="/api/admin/categories")
}