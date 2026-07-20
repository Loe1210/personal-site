# Blog Index Breathing Refresh Design

Date: 2026-07-17

## Goal

Refine the blog index experience so it feels closer to the user's visual references while preserving the existing dark mountain atmosphere and current split-page information architecture.

This round focuses on four concrete outcomes:

1. Make the sidebar social buttons visually match the reference more closely.
2. Replace the inline WeChat reveal with a centered modal that can later hold real text and QR imagery.
3. Shrink and soften the sidebar navigation so it feels lighter and less bulky.
4. Add the provided animated pet at the bottom of the sidebar and further reduce article card visual weight to improve breathing room.

## Confirmed Direction

### 1. Social Buttons

The sidebar social area will use a two-row, three-column circular button grid aligned with the reference style.

Buttons:
- WeChat: highlighted accent button
- Email: white circular link button
- GitHub: white circular link button
- Remaining slots: reserved for future links if needed, but current implementation will prioritize the three active controls and keep layout visually consistent

Visual rules:
- White circular buttons for ordinary links
- One accent-colored WeChat button
- Tighter icon sizing and more even spacing than the current implementation
- Keep the existing dark sidebar palette around the buttons

### 2. WeChat Modal

WeChat interaction will change from an inline expander to a centered modal overlay.

First implementation will include:
- Page dimming overlay
- Centered white modal card
- Close button in the upper-right corner
- Placeholder title text
- Placeholder descriptive paragraph
- Placeholder image area using a temporary image block
- Space reserved for later replacement with a real QR code or social image asset

Content strategy for now:
- Use neutral placeholder copy indicating this area is for WeChat contact / communication
- Use a placeholder image block or existing neutral stand-in so the modal structure is complete even before final content arrives

### 3. Sidebar Navigation

The navigation list will stay in the sidebar but become lighter and more reference-like.

Changes:
- Reduce font size slightly
- Reduce row height and padding
- Reduce icon footprint slightly
- Keep icon + label + active chevron pattern
- Preserve current links, including the admin entry in the navigation list
- Remove any remaining oversized or heavy card feel from the navigation area

Desired result:
- The nav reads as a clean vertical list instead of a stack of large blocks
- More empty space appears between social area, stats, nav, and pet area

### 4. Sidebar Pet

The provided assets in `/pet` will be integrated as the sidebar's bottom visual.

Available assets:
- `pet/pet.json`
- `pet/spritesheet.webp`

Placement:
- Desktop-first placement at the very bottom edge of the sidebar
- The pet should feel like it is standing on the sidebar floor, similar to the reference
- The pet becomes the primary visual element in the sidebar footer area

Behavior:
- Use the provided animation metadata rather than replacing it with a fake CSS-only mascot
- Keep the pet decorative and non-blocking
- Avoid overlap with nav items, footer copy, or viewport clipping

Responsive rule:
- If space becomes too tight on smaller widths, the pet may scale down or hide before it breaks layout integrity

### 5. Article Card Breathing Room

The homepage article list will be reduced in visual weight again.

Changes:
- Narrow the readable card width slightly
- Reduce cover height further
- Increase vertical spacing between title, meta, summary, and CTA
- Add more outer spacing between cards
- Keep the single-column stream structure already chosen
- Preserve the current dark glass styling language

Desired result:
- Cards feel more elegant and less oversized
- The page gains more breathing room without becoming sparse or weak

## Technical Approach

### HTML / Template Updates

Files expected to change:
- `static/blog/index.html`
- `static/blog/categories.html`
- `static/blog/tags.html`

Template work will:
- Update social button markup toward the reference grid
- Add WeChat modal container near the page shell or body end
- Keep admin inside nav
- Reserve a dedicated sidebar bottom mount point for the pet
- Bump cache-busting versions for changed assets

### CSS Updates

Primary file:
- `static/blog/css/blog.css`

CSS work will:
- Restyle social buttons
- Style the modal overlay and centered modal card
- Reduce nav scale and spacing
- Create a fixed sidebar-bottom pet area
- Further tighten article card dimensions and spacing
- Maintain the current mountain/dark-glass visual identity

### JavaScript Updates

Files expected to change:
- `static/blog/js/sidebar.js`
- `static/blog/js/list.js`
- possibly a small dedicated pet helper script if the current sidebar script should stay focused

JS work will:
- Open and close the WeChat modal
- Support overlay click and close button dismissal
- Optionally keep copy support for displayed contact text
- Load / initialize the pet resource from `/pet`
- Keep existing blog list behavior intact

### Pet Integration

Implementation will first inspect `pet.json` and determine the runtime contract for the provided pet asset.

Preferred approach:
- Use the provided metadata directly if it already defines sprite frame behavior
- If it only defines raw frame layout, create a very small renderer that reads the sheet and animates it in place
- Keep the implementation isolated so it does not interfere with the blog list scripts

## Risks And Safeguards

### Risk: Sidebar Becomes Too Busy

Safeguard:
- Reduce nav scale before adding the pet
- Keep the pet visually low and separate from the main nav list
- Remove or soften unnecessary footer weight if needed

### Risk: Modal Feels Out Of Style

Safeguard:
- Use a bright modal card inside a soft dim overlay, following the user's reference while still sitting on the darker site background

### Risk: Pet Asset Contract Is Non-Obvious

Safeguard:
- Inspect `pet.json` before implementation and adapt the renderer to the asset format rather than assuming a generic sprite layout

### Risk: Homepage Cards Become Too Small

Safeguard:
- Reduce size incrementally, preserving readability and touch target quality

## Verification Plan

After implementation:
- Run `go test ./...`
- Rebuild Docker frontend deployment
- Confirm `/blog/` serves updated asset versions
- Confirm WeChat modal opens and closes correctly
- Confirm article CTA links still navigate correctly
- Confirm the pet displays at the sidebar bottom without overlap

## Out Of Scope

This round will not:
- Redesign category or tag page information architecture again
- Replace the whole visual language with the reference site's pastel theme
- Introduce heavy external pet frameworks if the local assets can be rendered simply
- Finalize permanent WeChat text or QR content beyond a polished placeholder modal