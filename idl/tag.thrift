namespace go tag

struct Tag {
    1: i64 id
    2: string name
    3: string slug
    4: string description
    5: string created_at
    6: string updated_at
}

struct CreateTagRequest {
    1: string name (api.body="name")
    2: string slug (api.body="slug")
    3: string description (api.body="description")
}

struct CreateTagResponse {
    1: Tag tag
    2: string message
}

struct UpdateTagRequest {
    1: i64 id (api.path="id")
    2: string name (api.body="name")
    3: string slug (api.body="slug")
    4: string description (api.body="description")
}

struct UpdateTagResponse {
    1: Tag tag
    2: string message
}

struct DeleteTagRequest {
    1: i64 id (api.path="id")
}

struct DeleteTagResponse {
    1: string message
}

struct ListTagsRequest {
}

struct ListTagsResponse {
    1: list<Tag> list
}

service TagService {
    CreateTagResponse CreateTag(1: CreateTagRequest req) (api.post="/api/admin/tags")
    UpdateTagResponse UpdateTag(1: UpdateTagRequest req) (api.put="/api/admin/tags/:id")
    DeleteTagResponse DeleteTag(1: DeleteTagRequest req) (api.delete="/api/admin/tags/:id")
    ListTagsResponse ListTags(1: ListTagsRequest req) (api.get="/api/tags")
    ListTagsResponse ListAdminTags(1: ListTagsRequest req) (api.get="/api/admin/tags")
}
