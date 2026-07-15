# Task 3 Report - Chunk reception, tmp area, and resume

Status: DONE

## Completed
- Added `TmpStorage` for streamed chunk writes into `static/uploads/tmp` by default.
- Added `ChunkService` to stream chunk bodies to disk, record chunk metadata after successful writes, and update upload task progress.
- Added a chunk upload HTTP endpoint that reads from the request body stream instead of loading the whole payload into memory.
- Preserved the existing small-file upload flow.
- Wired the chunk service into `media-service` startup and route registration.

## Verification
- `go test ./services/media-service/internal/dal/storage -run TestTmpStorageWritesChunkToTmpPath -count=1`
- `go test ./services/media-service/internal/service -run TestChunkServiceWritesChunkToTmpPath -count=1`
- `go test ./services/media-service/...`

## Notes
- Chunk paths are written under the local tmp directory using deterministic names, which keeps resume bookkeeping stable.
- The task still relies on the request-supplied user id for now; later auth integration should replace that with shared session context.

## Commit
- `9510418`
