# Task 2 Report - Upload task service and state machine

Status: DONE

## Completed
- Added `UploadTaskService` with `InitUpload`, `GetUpload`, `CancelUpload`, and `CompleteUpload`.
- Added upload task sizing/default chunk calculation and 24-hour task expiration default.
- Added HTTP task endpoints under `/upload/tasks/*`.
- Wired upload task repositories into `media-service` startup and route registration.
- Preserved the existing `/upload` small-file flow.
- Tightened the size-limit test so it validates the intended failure path.

## Verification
- `go test ./services/media-service/internal/service -run TestInitUploadRejectsTooLargeFile -count=1`
- `go test ./services/media-service/...`

## Notes
- `CompleteUpload` currently only transitions status to `completed`; actual chunk merge and file record creation remain for later tasks.
- Task ownership is enforced through repository lookups using `upload_id` plus `user_id`.
