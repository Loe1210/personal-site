.PHONY: micro-up micro-down micro-test micro-smoke proto-check proto-gen micro-redeploy deploy-static deploy-frontend deploy-code verify-prod

SERVER ?= 117.72.95.156
REMOTE_DIR ?= /opt/personal-web
DEPLOY_ARCHIVE ?= personal-web-static-update.tar.gz
LOCAL_COMPOSE_FILES = -f deploy/docker/compose.yaml
PROD_COMPOSE_FILES = -f deploy/docker/compose.yaml

micro-up:
	docker compose -f deploy/docker/compose.yaml up -d --build

micro-down:
	docker compose -f deploy/docker/compose.yaml down

micro-redeploy:
	docker compose $(LOCAL_COMPOSE_FILES) up -d --build frontend

micro-test:
	go test ./...

micro-smoke:
	powershell -ExecutionPolicy Bypass -File scripts/smoke/microservices_smoke.ps1

proto-check:
	go run ./scripts/proto/check

proto-gen: proto-check
	powershell -ExecutionPolicy Bypass -File scripts/proto/gen.ps1

deploy-static:
	tar -czf $(DEPLOY_ARCHIVE) static pet frontend Dockerfile blog_index_ui_test.go
	scp $(DEPLOY_ARCHIVE) root@$(SERVER):/tmp/$(DEPLOY_ARCHIVE)
	ssh root@$(SERVER) "cd $(REMOTE_DIR) && tar -xzf /tmp/$(DEPLOY_ARCHIVE) && docker compose $(PROD_COMPOSE_FILES) up -d --build frontend"

deploy-frontend: deploy-static verify-prod
deploy-code:
	ssh root@$(SERVER) "cd $(REMOTE_DIR) && git fetch origin $(BRANCH) && git checkout $(BRANCH) && git pull --ff-only origin $(BRANCH) && docker compose $(PROD_COMPOSE_FILES) up -d --build frontend media-service content-service web-bff gateway"
verify-prod:
	curl -I --max-time 10 http://$(SERVER):8080/blog/
	curl -s --max-time 10 http://$(SERVER):8080/blog/ | grep -E 'blog.css|pet.js'
