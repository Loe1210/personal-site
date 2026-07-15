.PHONY: micro-up micro-down micro-test micro-smoke proto-check proto-gen

micro-up:
	docker compose -f deploy/docker/compose.yaml up -d --build

micro-down:
	docker compose -f deploy/docker/compose.yaml down

micro-test:
	go test ./...

micro-smoke:
	powershell -ExecutionPolicy Bypass -File scripts/smoke/microservices_smoke.ps1

proto-check:
	go run ./scripts/proto/check

proto-gen: proto-check
	powershell -ExecutionPolicy Bypass -File scripts/proto/gen.ps1
