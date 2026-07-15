# Task 1 Report - Upload task & metadata model

Status: DONE

## Completed
- Added upload task and upload chunk domain models in `services/media-service/internal/model/`.
- Added upload task and upload chunk GORM repositories with create, query, save, and progress update support.
- Added DB migration wiring for `upload_tasks` and `upload_chunks`.
- Added a repository test covering task persistence, chunk persistence, and progress reload.
- Extended the `files` table schema to keep upload and digest metadata needed by later chunk merge and resume steps.

## Verification
- `go test ./services/media-service/internal/dal/db -run TestUploadTaskRepositoryStoresStateAndChunks -count=1`
- `go test ./services/media-service/internal/dal/db ./services/media-service/internal/model`
- `go test ./services/media-service/...`

## Notes
- This task only establishes the metadata and persistence base for resumable upload flow.
- No code-level blockers remain for Task 1.
