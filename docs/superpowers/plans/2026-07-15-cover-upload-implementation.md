# 封面上传实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 后台支持本地封面图片上传，图片落地到 media-service 并通过 nginx 暴露为可访问 URL，文章保存时自动写入 `cover_image`。

**Architecture:** 后台编辑器只负责选择本地文件与回填 URL；媒体文件由 media-service 统一接收、校验和落盘；frontend/nginx 通过共享卷直接提供静态访问。文章数据仍只保存 URL，不保存二进制。

**Tech Stack:** Hertz, media-service, Nginx, Docker Compose, Go, 原生前端 JavaScript

## Global Constraints

- 封面文件只允许常见图片类型上传。
- 文章表单最终提交的是 `cover_image` URL，不提交文件本体。
- 上传后文件必须可通过 nginx 静态地址访问。
- 保持现有 `/api/media/upload` 路由与微服务结构。

---

### Task 1: media-service 图片上传约束

**Files:**
- Modify: `services/media-service/internal/service/media_service.go`
- Modify: `services/media-service/internal/service/media_service_test.go`

**Interfaces:**
- Consumes: `model.UploadInput{FileName, Content, ContentType, BizType}`
- Produces: `model.FileRecord.URL`, `model.FileRecord.Path`

- [x] **Step 1: Write the failing test**
- [x] **Step 2: Run test to verify it fails**
- [x] **Step 3: Write minimal implementation**
- [x] **Step 4: Run test to verify it passes**

### Task 2: shared static upload volume

**Files:**
- Modify: `deploy/docker/compose.yaml`
- Modify: `frontend/nginx/default.conf`

**Interfaces:**
- Consumes: media-service 落盘目录与 public_base_path
- Produces: nginx 可直接访问的 `/static/uploads/images/...`

- [x] **Step 1: Mount shared upload volume to media-service and frontend**
- [x] **Step 2: Verify compose config parses**
- [x] **Step 3: Rebuild containers and confirm static URL returns 200**

### Task 3: admin 封面上传控件

**Files:**
- Modify: `static/admin/index.html`
- Modify: `static/admin/js/admin.js`
- Modify: `static/admin/css/admin.css`

**Interfaces:**
- Consumes: `/api/media/upload`
- Produces: hidden `cover_image` value + preview UI

- [x] **Step 1: Add local file chooser and preview area**
- [x] **Step 2: Upload selected file and store returned URL**
- [x] **Step 3: Clear button resets preview and hidden field**
- [x] **Step 4: Keep article save payload using `cover_image`**

### Task 4: end-to-end verification

**Files:**
- None

**Interfaces:**
- Consumes: `/api/media/upload`, `/static/uploads/images/...`
- Produces: confirmed upload response and public file URL

- [x] **Step 1: Upload a real PNG through the API**
- [x] **Step 2: Confirm returned URL points to nginx-served file**
- [x] **Step 3: Rebuild frontend and confirm new admin assets load**
