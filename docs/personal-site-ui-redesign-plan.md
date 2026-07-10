# Personal Site UI Redesign Plan

## 1. Goal

The current UI is functional but not yet distinctive enough.

The redesign goal is:

- keep the site technical and personal
- make the visual identity stronger
- improve first impression
- keep the blog readable
- keep the admin efficient
- avoid generic template feeling

This redesign will follow a strict sequence:

1. reference study
2. style direction proposals
3. page-level UI mockups
4. visual direction confirmation
5. implementation

No frontend rework should begin until the page-level mockups are approved.

---

## 2. Reference Study

This redesign does not rely on a single inspiration source. Instead, it extracts strengths from several personal sites and portfolio patterns.

### 2.1 SimonAKing/HomePage

Already confirmed as a major reference for:

- dark immersive landing feeling
- centered first-screen composition
- sparse but deliberate navigation
- technical mood

Reference:

- [SimonAKing/HomePage](https://github.com/SimonAKing/HomePage/tree/master)

### 2.2 Lee Robinson

Useful reference for:

- clean writing-first structure
- strong personal positioning
- simple but confident information hierarchy
- content clarity without noise

Reference:

- [leerob.com](https://leerob.com/)

### 2.3 Brittany Chiang

Useful reference for:

- polished developer-personal-brand balance
- structured section rhythm
- strong portfolio-grade spacing and typography
- excellent “professional but personal” tone

Reference:

- [brittanychiang.com](https://brittanychiang.com/)

### 2.4 Rauno Freiberg

Useful reference for:

- design confidence
- reduced but memorable copy
- strong identity through spacing and tone
- “crafted” feeling instead of blog-template feeling

Reference:

- [rauno.me](https://rauno.me/)

### 2.5 Creative Bloq Portfolio Pattern Summary

Useful as a pattern summary rather than a source to copy:

- strong portfolio sites reveal personality, not only content
- grid systems can feel distinctive if typography and spacing are intentional
- minimalist layouts still work when the hierarchy is sharp

Reference:

- [Creative Bloq portfolio examples](https://www.creativebloq.com/portfolios/examples-712368)

---

## 3. Design Problems To Solve

The current UI needs improvement in these areas:

### 3.1 Identity is not sharp enough

The site runs, but it does not yet feel like a memorable personal brand.

### 3.2 Page rhythm is still too plain

The current page structure is usable, but the visual pacing is not yet strong enough to create a premium feel.

### 3.3 Frontend and admin feel disconnected

The public site and admin console should feel like two modes of one system, not two unrelated interfaces.

### 3.4 The homepage needs stronger narrative control

The first screen should immediately answer:

- who you are
- what this site is about
- why someone should continue reading

---

## 4. Redesign Principles

The redesign must follow these principles:

### 4.1 No template look

Avoid:

- generic hero + card + footer structure
- overused startup gradients
- interchangeable SaaS layouts

### 4.2 Technical but warm

The site should feel like a backend engineer’s personal station, not a corporate product page.

### 4.3 Strong typography first

The redesign should rely more on:

- hierarchy
- spacing
- contrast
- alignment

and less on decorative effects.

### 4.4 Public site and admin share one system

They can differ in density and purpose, but should share:

- color language
- motion language
- border language
- card language
- typography logic

### 4.5 Reading experience matters

This is still a blog-driven site, so article detail pages must prioritize:

- readable width
- comfortable line-height
- strong heading hierarchy
- code block clarity

---

## 5. Proposed Style Directions

The redesign phase will produce 3 candidate visual directions.

### Direction A: Terminal Gallery

Keywords:

- black graphite base
- precise grid
- terminal-inspired details
- glowing green accents
- centered hero

Strengths:

- closest to the SimonAKing mood
- very aligned with Go backend learning identity
- memorable technical atmosphere

Risks:

- can become too cold if not balanced with softer copy and spacing

### Direction B: Editorial Engineer

Keywords:

- dark editorial layout
- high-contrast serif + sans pairing
- writing-first composition
- more mature and reflective tone

Strengths:

- better for long-form blog reading
- feels more like a serious technical writer’s home

Risks:

- less “wow” on first screen if not art-directed carefully

### Direction C: Crafted Minimal Lab

Keywords:

- reduced layout
- strong whitespace
- subtle motion
- premium minimalism
- product-designer level refinement

Strengths:

- strongest “high taste” direction
- could make the site feel timeless and personal

Risks:

- easiest to make too empty if execution is weak

---

## 6. Mockup Scope

Before implementation, the following page mockups must be produced:

### Public Pages

- Home
- Blog list
- Article detail
- About

### Admin Pages

- Admin login
- Admin dashboard
- Admin articles list
- Admin article editor
- Admin taxonomy

The mockups do not need to be pixel-perfect production assets yet, but they must be strong enough for visual approval.

---

## 7. Mockup Output Strategy

The redesign review will happen in two rounds.

### Round 1: Direction Boards

Output:

- 2 to 3 homepage-first visual directions
- each direction includes:
  - hero composition
  - type hierarchy
  - color treatment
  - nav pattern
  - card treatment

Goal:

- choose one main direction

### Round 2: Full Page Mockups

After one direction is selected, produce:

- Home
- Blog list
- Article detail
- About
- Admin login
- Admin dashboard
- Admin articles
- Admin editor
- Admin taxonomy

Goal:

- approve the full UI system before development

---

## 8. Visual System To Define Before Coding

Before implementation, the chosen direction must define:

- primary background
- accent color
- typography pair
- card radius system
- border opacity system
- glow usage rules
- spacing scale
- button styles
- badge styles
- form styles
- empty state style

This avoids visual drift during coding.

---

## 9. Development Sequence After Approval

Once the mockups are approved, development should follow this order:

1. shared design tokens
2. homepage
3. blog list
4. article detail
5. about
6. shared admin shell
7. admin login
8. admin articles
9. admin editor
10. admin taxonomy
11. responsive tuning
12. visual polish

---

## 10. Acceptance Criteria

The redesign is successful only if:

- the homepage feels memorable within 3 seconds
- the site no longer looks like a generic starter template
- article reading experience clearly improves
- admin and public UI feel like one product family
- the chosen direction feels personal, technical, and premium
- visual approval happens before code implementation

---

## 11. Current Default Execution Plan

The next design step should be:

1. produce 3 homepage-centered UI direction boards
2. let the user choose one
3. produce full-page mockups for that direction
4. only then start implementation

This is the safest way to improve the UI without repeating the previous “code first, redesign later” cycle.

## 12. Selection Result

The approved visual direction is:

- `Direction A: Terminal Gallery`

Implementation status:

- public frontend implementation has started and been applied to:
  - Home
  - Blog list
  - Article detail
  - About

Not included in this redesign coding round:

- backend API changes
- admin console visual redesign
- data model changes

Current execution rule:

- continue public visual polishing on top of `Terminal Gallery`
- keep backend stable
- decide separately whether the admin console should be visually upgraded to match this public direction
