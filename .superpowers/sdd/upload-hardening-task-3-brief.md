# Task 3 Brief - Chunk reception, tmp area, and resume

You are implementing Task 3 of the local file upload hardening plan in `docs/superpowers/plans/2026-07-15-local-file-upload-hardening.md`.

## Where this task fits
Task 1 added upload task/chunk persistence. Task 2 added the upload task state machine and route wiring. This task must add streamed chunk reception into a local temporary area and support resume bookkeeping.

## Requirements
- Use the repositories and models from Tasks 1 and 2.
- Keep local disk storage only.
- Do not read the whole large file into memory at once.
- Store chunk bytes directly into the temporary upload area.
- Record chunk metadata only after the write succeeds.
- Keep the existing small-file upload flow intact.

## Files in scope
- Create: `services/media-service/internal/dal/storage/tmp_storage.go`
- Create: `services/media-service/internal/dal/storage/tmp_storage_test.go`
- Modify: `services/media-service/internal/service/media_service.go`
- Create: `services/media-service/internal/service/chunk_service.go`
- Create: `services/media-service/internal/service/chunk_service_test.go`
- Modify: `services/media-service/biz/upload/handler.go`
- Modify: `services/media-service/biz/upload/route.go`

## Context you should assume
- Task 1 and Task 2 are already committed and should not be reverted.
- Upload task ownership comes from repository lookups using `upload_id` plus `user_id`.
- The temp path must live under the local upload tmp directory, not in a shared object store.
- You are not alone in the codebase. Do not revert or overwrite edits made by others.

## Test-first target
Start with the failing test from the plan:

```go
func TestChunkServiceWritesChunkToTmpPath(t *testing.T) {
    // upload one chunk, assert temp file exists and metadata is recorded
}
```

Then make it pass with the smallest change set that preserves the plan.

## Report file
Write your full report to:
`.superpowers/sdd/upload-hardening-task-3-report.md`

Report status should include:
- what you implemented
- tests run and results
- any concerns or follow-up notes

## Important
- Commit your changes when done and include the commit hash in the report.
- Keep the change scoped to Task 3 only.
