# Tag Sort And Sidebar Indicator Design

## Goal

Turn tag selection into a focused article archive while preserving the existing constellation landing view and shared blog sidebar.

## User Experience

- The tag page heading is `Tag.sort()`.
- The category page heading is `Category.sort()`.
- `/blog/tags` with no `focus` query renders the constellation.
- Selecting a tag changes the same route to `/blog/tags?focus=<tag>` and swaps the main content to a tag archive.
- The archive heading shows the selected tag and its article count.
- Each archive item links to its article and shows title, publication date, and any other associated tags.
- A return control clears `focus` and restores the constellation without a full page reload.
- Empty and failed states explain why no results are visible.

## Data Flow

- Fetch tag and category metadata when the page loads.
- Do not fetch every article until a tag is selected.
- On selection, request the public article list and hydrate only records missing tag data through `BlogAPI.getPost(id)`.
- Filter hydrated records by the selected tag. Detail requests may fail independently without breaking the archive view.

## Sidebar Navigation

- Insert one decorative navigation indicator into each desktop sidebar.
- On pointer enter or keyboard focus, the indicator moves to the hovered or focused item using its measured top position and height.
- On pointer leave or focus exit, it returns to the current page item.
- Navigation links remain normal links; the indicator never intercepts clicks.
- On compact layouts, the indicator is disabled so the existing responsive navigation remains intact.

## Styling

- Preserve the current blue background, sidebar dimensions, typography, and motion language.
- The archive uses vertically stacked, dark translucent cards with a slim bright edge, mirroring the information hierarchy in the reference while matching the existing site palette.
- The moving sidebar indicator uses a soft pale-blue translucent rectangle with a short transform transition.
- Respect `prefers-reduced-motion` by disabling nonessential movement.

## Validation

- Source-level tests cover the tag archive view, detail hydration fallback, return action, renamed page headings, and navigation indicator hooks.
- Run `node --check` for edited scripts, `go test ./...`, and rebuild the frontend container.