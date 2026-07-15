# Task 2 Brief - Upload task service and state machine

You are implementing Task 2 of the local file upload hardening plan in `docs/superpowers/plans/2026-07-15-local-file-upload-hardening.md`.

## Where this task fits
Task 1 already added upload task/chunk persistence, migrations, and repository tests. This task must build the service/state-machine layer on top of those repositories so later chunk upload and merge work has a real API contract.

## Requirements
- Use the existing upload task/chunk repositories from Task 1.
- Keep local disk storage only. Do not introduce OSS, MinIO, or CDN distribution.
- Preserve the existing small-file upload flow and `/api/media/upload`.
- The task state machine must support create, query, cancel, and complete.
- A task is only accessible by its creator.
- Follow the plan’s intended interfaces:
  - `InitUpload`
  - `GetUpload`
  - `CancelUpload`
  - `CompleteUpload`

## Files in scope
- Create: `services/media-service/internal/service/upload_task_service.go`
- Create: `services/media-service/internal/service/upload_task_service_test.go`
- Modify: `services/media-service/cmd/main.go`
- Modify: `services/media-service/cmd/router.go`
- Modify: `services/media-service/biz/router.go`
- Create: `services/media-service/biz/upload/task_handler.go`
- Create: `services/media-service/biz/upload/task_route.go`

## Context you should assume
- `services/media-service/internal/model/upload_task.go` and `upload_chunk.go` already exist.
- `services/media-service/internal/dal/db/upload_task_repository.go` and `upload_chunk_repository.go` already exist.
- The database migration and test sqlite dependency are already in place.
- Do not remove or rewrite unrelated changes made by other tasks.

## Test-first target
Start with the failing test from the plan:

```go
func TestInitUploadRejectsTooLargeFile(t *testing.T) {
    // init with file_size larger than configured limit, expect error
}
```

Then make it pass with the smallest change set that preserves the plan.

## Report file
Write your full report to:
`.superpowers/sdd/upload-hardening-task-2-report.md`

Report status should include:
- what you implemented
- tests run and results
- any concerns or follow-up notes

## Important
- You are not alone in the codebase. Do not revert or overwrite edits made by others.
- Keep your work tightly scoped to Task 2.
- Commit your changes when done and include the commit hash in the report.
