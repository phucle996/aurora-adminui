---
name: adminui-design-system
description: Use this guide when building or updating the adminui so it keeps the current Vite + React + TanStack stack while matching the visual language, information density, and operational dashboard feel of controlplane/ui.
---

# Admin UI Design Skill

## Purpose

Use this guide whenever work touches `adminui` UI, layout, theming, page composition, component styling, or interaction design.

This guide exists to keep `adminui` visually aligned with `controlplane/ui` while **preserving the current `adminui` framework**:

- `Vite`
- `React 19`
- `TanStack Router`
- `TanStack Query`
- `Tailwind CSS v4`
- `Radix UI`
- `shadcn`-style component primitives

Do **not** migrate `adminui` to `Next.js`.
Do **not** rebuild it as a copy-paste clone of `controlplane/ui`.
Do **not** keep the current generic demo-dashboard styling if it conflicts with the controlplane look.

The goal is:

- same product family feel as `controlplane/ui`
- same operational B2B dashboard tone
- same spacing discipline, card structure, tokens, and hierarchy
- same seriousness and clarity
- but implemented with `adminui`'s existing stack and routing model

## Canonical References

When implementing design decisions, prefer these files as the source of truth:

- Controlplane theme tokens: [globals.css](../controlplane/ui/src/app/globals.css)
- Controlplane card structure: [ComponentCard.tsx](../controlplane/ui/src/components/common/ComponentCard.tsx)
- Controlplane breadcrumb/title rhythm: [PageBreadCrumb.tsx](../controlplane/ui/src/components/common/PageBreadCrumb.tsx)
- Controlplane layout shell/sidebar: [AppSidebar.tsx](../controlplane/ui/src/layout/AppSidebar.tsx)
- Admin UI current theme entry: [theme.css](./src/styles/theme.css)
- Admin UI current style entry: [index.css](./src/styles/index.css)
- Admin UI root shell: [__root.tsx](./src/routes/__root.tsx)
- Admin UI current sidebar data: [sidebar-data.ts](./src/components/layout/data/sidebar-data.ts)

If a design decision is ambiguous, choose the option that makes `adminui` feel closer to `controlplane/ui`.

## Non-Negotiables

### Keep the framework

`adminui` must stay on:

- `Vite`
- `TanStack Router`
- `TanStack Query`
- `Radix`
- `shadcn`-style primitives
- `Tailwind v4`

Do not introduce:

- `Next.js`
- `MUI`
- `Ant Design`
- a second component system
- inline one-off styling that bypasses the token system

### Match the controlplane visual language

`adminui` should look like it belongs to the same product suite as `controlplane/ui`, especially in:

- color behavior
- density
- card rhythm
- typography hierarchy
- sidebar feel
- form treatment
- table treatment
- status badge treatment
- chart tone

### Favor operational dashboard aesthetics

This is not a marketing site and not a generic SaaS starter dashboard.

The desired tone is:

- precise
- information-dense
- modern
- calm
- slightly enterprise
- trustworthy
- not playful
- not over-animated

## Design Direction

### Overall visual identity

The visual identity should feel like:

- a cloud/infrastructure operations console
- crisp white surfaces on soft gray backgrounds
- clear blue brand actions
- restrained use of accent colors
- strong sectioning
- readable dense data layouts

Avoid:

- floating “glassmorphism” UI
- big gradients as the main page treatment
- excessive blur
- oversized rounded pills everywhere
- purple-first palettes
- soft consumer app aesthetics
- dribbble-like decoration that hurts scannability

## Typography

### Target hierarchy

Match the controlplane hierarchy instead of the current generic `Inter/Manrope` admin template feel.

Preferred direction:

- page and section headings should visually feel like `controlplane/ui`
- body copy should be clean, neutral, and compact
- metadata text should stay quiet and secondary

### Font guidance

`controlplane/ui` uses `Outfit` as the primary family.

For `adminui`, the preferred implementation is:

1. adopt `Outfit` as the primary UI family as well
2. if a second font is kept, use it sparingly and never let it dominate the app shell

If `Inter` and `Manrope` remain temporarily:

- use one as body fallback only
- do not mix them arbitrarily across cards, headers, and tables
- move the overall feel toward the proportions and weight balance of `controlplane/ui`

### Typography rules

- page title: strong, compact, dark, not oversized
- section title: medium weight, calm, high contrast
- helper text: smaller and muted
- stat values: large but not flashy
- table headings: compact, uppercase only if already established, otherwise sentence case
- breadcrumb text: small and muted

Avoid:

- giant hero typography
- thin light weights for operational text
- large font jumps between adjacent sections

## Color System

### Primary palette

Port the **semantics** of the controlplane palette into `adminui`.

Core palette direction:

- brand blue for primary actions and active states
- neutral grays for surfaces, borders, and layout scaffolding
- green for healthy/success
- red for destructive/error
- amber/orange for warning and pending

The controlplane token families worth mirroring are:

- `brand-*`
- `gray-*`
- `success-*`
- `error-*`
- `warning-*`
- `blue-light-*`

### Token strategy for adminui

Refactor `adminui/src/styles/theme.css` so semantic tokens map to the controlplane tone:

- `--background` should match the calm off-white/gray page feel
- `--card` should be white or near-white
- `--border` should be a clear light neutral, not too faint
- `--primary` should be a blue close in behavior to controlplane `brand-500`
- `--muted` should support dense dashboard sections without looking washed out
- `--ring` should match the blue focus ring language
- chart tokens should come from the same family as controlplane status/brand colors

### Surface layering

Default page treatment should be:

- app shell background: soft gray
- cards: white
- elevated overlays: white with stronger shadow
- sidebar: white or near-white, not tinted novelty colors

Dark mode, if kept, should also feel like the same family as `controlplane/ui`, not a separate theme identity.

## Radius, Border, and Shadow

### Radius

Use moderate radius similar to controlplane cards:

- cards: around `rounded-2xl`
- form controls: consistent medium radius
- badges: rounded but not oversized
- dialogs/popovers: same family as cards

Avoid:

- ultra-square harsh edges
- exaggerated pill-heavy shapes everywhere

### Borders

Borders are a major part of the controlplane look.

Rules:

- cards should usually have a visible, light border
- tables should rely on subtle separators
- inputs should have clear boundaries
- tabs and section dividers should be understated but explicit

### Shadows

Use shadows sparingly and consistently.

Target feel:

- subtle card elevation
- stronger shadow only for overlays, popovers, dropdowns, dialogs
- do not make every card look like it is floating dramatically

Mirror the controlplane pattern:

- tiny shadow for standard cards
- medium shadow for overlays
- strong shadow only for special surfaces

## Layout System

### Shell

`adminui` shell should be brought closer to `controlplane/ui`:

- left sidebar as the main navigation spine
- content area with clear horizontal padding
- page content stacked with deliberate vertical rhythm
- top-level pages should feel anchored, not like floating demo widgets

### Content width

Use a wide but controlled app layout:

- comfortable for tables and dashboards
- not cramped
- not edge-to-edge raw

Prefer:

- consistent horizontal padding
- 24px to 32px style rhythm on desktop
- tighter but still breathable spacing on smaller screens

### Page composition

Default page structure should follow this rhythm:

1. breadcrumb / section path
2. page title + short contextual description
3. top action bar if needed
4. summary cards or filters
5. primary content cards
6. secondary details, tables, charts, or forms

This should feel extremely close to `controlplane/ui` pages.

## Sidebar and Navigation

### Navigation tone

The current `adminui` sample sidebar data feels too template-like.

The desired navigation should feel:

- product-specific
- operational
- domain-grouped
- sober

### Sidebar styling rules

Match the controlplane sidebar treatment:

- compact label sizing
- strong active state
- muted inactive items
- clean grouping headers
- icons that support hierarchy without overpowering text

Use active states similar in spirit to controlplane:

- light brand-tinted active background
- blue text/icon emphasis
- restrained hover states

### Menu content rules

Do not keep placeholder/demo information architecture once real product sections exist.

Navigation labels should be:

- short
- operational
- domain-based
- unambiguous

## Cards

### Card is the primary content unit

`adminui` should adopt the same card philosophy as controlplane:

- cards are the default section container
- each card should have a clear header/body structure
- cards should be visually predictable across the app

### Card structure

Preferred pattern:

- outer card: rounded, bordered, white background
- optional header: title + description + header actions
- body padding: consistent and not cramped

Recommended rhythm based on controlplane:

- header padding larger than body micro-content spacing
- body supports forms, tables, stats, or charts without changing the shell look

Avoid:

- ad hoc panel styles on every page
- mixing 4-5 different card radii and border treatments

## Forms

### Form design goals

Admin forms should feel like the `New Virtual Machine` page in spirit:

- guided
- segmented
- easy to scan
- not noisy
- each step or section visually grouped

### Form rules

- group related fields inside cards or clear subsections
- use concise labels
- helper text should exist only where it reduces user risk
- validation states should be visible but not aggressive
- destructive actions must look distinct from primary submit actions

### Inputs

Inputs should feel production-grade:

- medium height
- clear border
- quiet placeholder
- strong focus ring
- no soft ghosty outlines

### Selectors and radio groups

For package/image/resource choices:

- make options feel like selectable infrastructure SKUs
- use grid cards where appropriate
- emphasize active selection with border + subtle brand background
- include icons only when they help recognition

The current `New Virtual Machine` treatment in controlplane is a good direction for high-value selectors.

## Tables

### Table philosophy

Tables should feel like operational data surfaces, not generic data grid demos.

Rules:

- compact but readable rows
- quiet header background if needed
- strong sorting/filter affordances only when implemented
- statuses should be easy to scan
- row actions should not visually dominate the row

### Table styling

- neutral borders
- small/medium row height
- clear hover state
- muted secondary metadata
- stronger emphasis on primary object name

Avoid:

- heavy zebra striping
- oversized padding
- colorful row backgrounds unless status demands it

## Badges and Status

Statuses should mirror controlplane semantics:

- healthy/running/active -> green
- pending/provisioning -> amber or orange
- failed/error -> red
- inactive/stopped -> gray
- informational/assigned -> blue if needed

Badge rules:

- small
- readable
- consistent casing
- consistent radius
- no rainbow badge explosion on one page

## Charts and Metrics

Charts should support operators, not decorate dashboards.

Rules:

- use restrained palettes
- prioritize clarity over flourish
- grid lines should be subtle
- legends should be compact
- titles should explain the metric, not marketing language

When showing multiple metric cards/charts:

- align them in predictable grids
- reuse the same card structure as controlplane
- keep metric formatting consistent

Avoid:

- neon palettes
- too many series in one chart
- gratuitous animation

## Empty States, Loading, and Errors

### Empty states

Keep empty states practical:

- explain what is missing
- explain what the user can do next
- provide the primary CTA if relevant

Avoid whimsical illustrations unless they already exist in the product family.

### Loading states

Prefer:

- skeletons for cards/tables
- stable layout during background refresh
- no full-page flicker for small data refreshes

### Error states

Errors should be:

- visible
- brief
- actionable where possible

Keep the same operational tone as controlplane.

## Motion

Motion should be subtle and purposeful.

Allowed:

- small fade/slide transitions for overlays
- top-loading indicator for route changes
- gentle skeleton transitions

Avoid:

- large springy animations
- motion that changes page structure dramatically
- decorative animation on dashboards

## Responsive Behavior

`adminui` must preserve the same seriousness on smaller screens.

Rules:

- cards stack cleanly
- tables degrade gracefully with overflow or responsive summaries
- form sections collapse without losing hierarchy
- top action bars wrap predictably
- side navigation should remain understandable when collapsed

Do not solve mobile by shrinking everything until it becomes dense noise.

## Accessibility

Every redesign should preserve or improve:

- visible focus states
- contrast
- keyboard access
- semantic labels
- tab order

The controlplane look should not be reproduced at the expense of usability.

## Implementation Strategy

When modernizing `adminui` toward the controlplane style, work in this order:

1. **Theme tokens**
   - update `adminui/src/styles/theme.css`
   - align semantic colors, radii, ring, surface, sidebar tokens

2. **Global shell**
   - update root shell and sidebar styling
   - remove generic starter-dashboard feel

3. **Core primitives**
   - card
   - button
   - badge
   - input
   - select
   - table wrappers
   - dialog/popover/dropdown shells

4. **Page scaffolds**
   - breadcrumb/header blocks
   - filters/action bars
   - stats strips

5. **Domain pages**
   - users
   - tenants
   - plans
   - infrastructure/admin operations

Always fix tokens and primitives first before styling a single page in isolation.

## Mapping Rules From Controlplane UI

### Translate, do not clone

Because the frameworks differ, treat `controlplane/ui` as a visual and structural reference, not a code donor.

Examples:

- `ComponentCard.tsx` in controlplane should inspire an `adminui` card wrapper, not be ported literally
- controlplane theme tokens should be translated into `adminui` semantic variables
- controlplane page title rhythm should be reproduced with `adminui` layout components

### Semantic translation examples

- controlplane `brand-500` behavior -> `adminui` `--primary`
- controlplane `gray-50` page background -> `adminui` `--background`
- controlplane white cards with light border -> `adminui` `--card` + `--border`
- controlplane focus ring -> `adminui` `--ring`

## Anti-Patterns

Do not do any of the following unless the product intentionally changes direction:

- keep the stock shadcn-admin demo look
- use purple as the dominant accent
- use heavy gradients as page backgrounds
- mix multiple unrelated radii systems
- introduce one-off color values in components
- use oversized pills for all tabs, badges, and buttons
- make charts brighter than the data table they support
- use casual consumer-product copy in operational admin pages
- copy `controlplane/ui` component code directly into `adminui` without adapting it
- switch frameworks just to chase design parity

## Page-Level Heuristics

When building a new `adminui` page, ask:

1. Would this page visually belong next to controlplane pages?
2. Does the page read like an operations console instead of a template?
3. Are the tokens semantic and reusable?
4. Are the cards and headers using the same rhythm as the rest of the app?
5. Would a future engineer know where to extend the pattern?

If the answer is “no” to any of these, refine before shipping.

## Expected End State

When this guide is followed, `adminui` should:

- feel like the admin/control sibling of `controlplane/ui`
- keep its existing engineering stack
- lose the generic starter-dashboard identity
- support complex data-heavy pages without style drift
- have a stable token and primitive system that future pages can reuse

That is the target standard for all future `adminui` work.
