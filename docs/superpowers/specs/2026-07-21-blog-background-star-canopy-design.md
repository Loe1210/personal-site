# Blog Background And Star Canopy Design

## Goal

Unify every `/blog` page on the current `docs/background.png` art direction, add a lightweight animated star canopy, and refine sidebar and directory interactions.

## Experience

- Copy `docs/background.png` to the deployed blog background asset and use it on blog list, categories, tags, and post detail pages.
- Add `blog-page` classes to the post detail body so it shares the same fixed background.
- Add two CSS-only star layers to the top of each page. Stars use muted cyan and ice-blue tones, twinkle slowly, and never block clicks.
- Avatar hover scales the circular image and brightens its blue halo. Keyboard focus receives the same visible treatment.
- Remove the `directory-hero` bottom border on category and tag pages.

## Constraints

- Do not change article APIs, routes, or database data.
- Preserve the existing sidebar layout and responsive behavior.
- Disable star and avatar movement for `prefers-reduced-motion`.
- Cover art remains available on article detail pages only.

## Validation

- Source-level tests assert shared background asset references, star canopy CSS, avatar hover behavior, and border removal.
- Run JavaScript syntax checks, `go test ./...`, and rebuild the frontend container.