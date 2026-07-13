namespace go media

struct FileRecord {
    1: i64 id
    2: string original_name
    3: string url
    4: string path
    5: string content_type
    6: i64 size
    7: string biz_type
    8: string created_at
}

struct GetFileRequest {
    1: i64 id
}

struct GetFileResponse {
    1: FileRecord file
}

service MediaService {
    GetFileResponse GetFile(1: GetFileRequest req)
}
