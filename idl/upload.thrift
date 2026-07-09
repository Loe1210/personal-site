namespace go upload

struct UploadFile {
    1: i64 file_id
    2: string file_name
    3: string file_url
    4: string file_path
    5: string mime_type
    6: i64 size
    7: string biz_type
    8: string created_at
}

struct UploadImageRequest {
    1: string biz_type (api.form="biz_type")
}

struct UploadImageResponse {
    1: UploadFile upload
}

struct GetUploadInfoRequest {
    1: i64 id (api.path="id")
}

struct GetUploadInfoResponse {
    1: UploadFile upload
}

service UploadService {
    UploadImageResponse UploadImage(1: UploadImageRequest req) (api.post="/api/admin/upload")
    GetUploadInfoResponse GetUploadInfo(1: GetUploadInfoRequest req) (api.get="/api/admin/uploads/:id")
}