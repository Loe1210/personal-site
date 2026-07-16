# Task 3 Fix Brief - Upload chunk race guard

You are fixing the remaining Task 3 review finding in the local file upload hardening plan.

## Finding to fix
- High: `services/media-service/internal/service/chunk_service.go:62-99` and `services/media-service/internal/dal/db/upload_task_repository.go:64-71` allow a chunk upload to overwrite task state after a cancel/complete race. The service reads the task, writes the chunk, then blindly calls `UpdateProgress` with the last-read status; there is no status/version guard in the update. A chunk that started before cancellation can still persist and flip progress back onto a task that should no longer accept chunks.

## Required direction
- Add a status/version guard to the task progress update so stale chunk uploads cannot overwrite a task that was canceled or completed after the task was read.
- Keep local storage only and preserve the existing small-file upload flow.
- Update any tests that cover the task state transition behavior.

## Likely files
- `services/media-service/internal/dal/db/upload_task_repository.go`
- `services/media-service/internal/service/chunk_service.go`
- `services/media-service/internal/service/upload_task_service.go`
- `services/media-service/internal/service/chunk_service_test.go`
- `services/media-service/internal/service/upload_task_service_test.go`

## Report file
Append your fix note to:
`.superpowers/sdd/upload-hardening-task-3-report.md`

## Important
- Do not revert or overwrite edits made by others.
- Run the relevant tests and commit your changes.
