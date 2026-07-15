# Upload Hardening Task 1 Review Package

- Base: `ebcf46e`
- Head: `b003256`
- Task: upload task & metadata model

## What changed
- Added `UploadTask` and `UploadChunk` domain models under `services/media-service/internal/model/`.
- Added repositories for upload task state and chunk records.
- Wired the media-service migration path to include the new upload tables.
- Extended the file record schema with upload metadata needed by resumable uploads.
- Added a repository test covering task/chunk persistence and progress updates.

## Diff summary
```text
 .superpowers/sdd/upload-hardening-task-1-brief.md  |  70 ++++++++++++
 .superpowers/sdd/upload-hardening-task-1-report.md |  19 ++++
 go.mod                                             |  12 ++
 go.sum                                             |  22 ++++
 .../internal/dal/db/file_repository.go             |  11 +-
 services/media-service/internal/dal/db/migrate.go  |   2 +-
 .../internal/dal/db/upload_chunk_repository.go     |  83 ++++++++++++++
 .../internal/dal/db/upload_task_repository.go      | 125 +++++++++++++++++++++
 .../internal/dal/db/upload_task_repository_test.go |  79 +++++++++++++
 services/media-service/internal/model/file.go      |   5 +-
 .../media-service/internal/model/upload_chunk.go   |  13 +++
 .../media-service/internal/model/upload_task.go    |  29 +++++
 .../media-service/migrations/002_upload_tasks.sql  |  39 +++++++
```

## Verification
- `go test ./services/media-service/...`

## Reviewer focus
- Check schema/index choices for upload task lookups and expiration cleanup.
- Check repository method semantics around versioning and progress updates.
- Check whether the additional sqlite test dependency is acceptable for the repo style.
