# Personal Site Admin Frontend Plan

## 1. Goal

This phase adds the first usable admin console for the personal site.

The target is not a heavy CMS. The target is a compact admin frontend that makes the current backend capabilities actually usable in the browser.

The first release must support:

- session login
- current user detection
- article list
- article create and edit
- article delete
- category create and list
- tag create and list
- cover upload in the editor
- clear feedback for unauthenticated and forbidden states

## 2. Scope

The first admin frontend release includes these pages:

- `/admin/login`
- `/admin`
- `/admin/articles`
- `/admin/articles/new`
- `/admin/articles/:id/edit`
- `/admin/taxonomy`

Upload is integrated into the article editor instead of a standalone page in the first round.

## 3. Design Direction

The admin console should stay visually aligned with the public site, but move toward efficiency instead of atmosphere.

Keep:

- dark base
- grid / glow language
- soft glass cards
- green accent
- clean typography

Strengthen:

- form readability
- content density
- operation clarity
- status feedback
- permission feedback

Reduce:

- oversized animation
- decorative noise
- over-immersive first-screen effects

## 4. Information Architecture

Primary admin navigation:

- Dashboard
- Articles
- Taxonomy
- View Site
- Logout

Page responsibilities:

### Login

- username/password login
- redirect logged-in users to `/admin`
- show login failure message

### Dashboard

- show current user
- quick actions
- show article/category/tag counts

### Articles

- list articles
- filter by status
- keyword search
- create new article
- edit article
- delete article

### Editor

- title
- slug
- summary
- markdown body
- category
- tags
- status
- cover upload
- save draft / publish

### Taxonomy

- category list + create form
- tag list + create form

## 5. API Integration

Reuse existing backend APIs:

- `POST /api/admin/login`
- `POST /api/admin/logout`
- `GET /api/admin/me`
- `GET /api/admin/articles`
- `POST /api/admin/articles`
- `PUT /api/admin/articles/:id`
- `DELETE /api/admin/articles/:id`
- `GET /api/admin/categories`
- `POST /api/admin/categories`
- `GET /api/admin/tags`
- `POST /api/admin/tags`
- `POST /api/admin/upload`

Implementation note:

- the current editor reads existing article data by loading the admin article list and selecting by id, because the current frontend plan avoids adding a new backend article-detail route in this round.

## 6. Session Strategy

The admin frontend follows the existing session-based backend.

Rules:

- do not store tokens in local storage
- rely on browser cookies
- use `GET /api/admin/me` as the login-state source of truth
- redirect unauthenticated users to `/admin/login`
- show a clear forbidden state for `403`

## 7. Directory Layout

```text
templates/
  components/
    admin-layout.html
    admin-sidebar.html
    admin-topbar.html
  pages/
    admin/
      login.html
      dashboard.html
      articles.html
      article-edit.html
      taxonomy.html

static/
  css/
    admin.css
    admin-login.css
    admin-dashboard.css
    admin-articles.css
    admin-editor.css
    admin-taxonomy.css
  js/
    admin-common.js
    admin-login.js
    admin-dashboard.js
    admin-articles.js
    admin-editor.js
    admin-taxonomy.js
```

## 8. Execution Order

1. admin document
2. admin layout, sidebar, topbar
3. login page and login guard
4. dashboard
5. article list
6. article editor with upload
7. taxonomy page
8. validation and devlog

## 9. Acceptance Criteria

This phase is complete when:

- login works from the browser
- logout works from the browser
- current user is shown in the admin shell
- article list renders
- article create works
- article edit works
- article delete works
- category create works
- tag create works
- cover upload works in the editor
- unauthenticated users are redirected to login
- forbidden operations show a clear message
- the admin frontend matches the public site visual language
